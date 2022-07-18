//go:build linux || android
// +build linux android

package sers

// Copyright 2012 Michael Meier. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type baseport struct {
	fd           int
	f            *os.File
	platformData termiosPlatformData
}

func takeOverFD(fd int, fn string) (SerialPort, error) {
	bp := &baseport{
		fd: fd,
	}

	tio, err := bp.getattr()
	if err != nil {
		return nil, &Error{"putting fd in non-canonical mode", err}
	}

	cfmakeraw(tio)

	err = bp.setattr(tio)
	if err != nil {
		return nil, &Error{"putting fd in non-canonical mode", err}
	}

	bp.f = os.NewFile(uintptr(fd), fn)

	return bp, nil
}

// TakeOver accepts an open *os.File and returns a SerialPort representing the
// open file.
//
// Attention: This calls the .Fd() method of the *os.File and thus renders the
// deadline functionality unusable. Furthermore blocked readers may remain
// stuck after a Close() if no data arrives.
func TakeOver(f *os.File) (SerialPort, error) {
	if f == nil {
		return nil, &ParameterError{"f", "needs to be non-nil"}
	}
	bp := &baseport{fd: int(f.Fd()), f: f}

	return bp, nil
}

func (bp *baseport) Read(b []byte) (int, error) {
	n, err := bp.f.Read(b)

	// timeout gets reported as EOF
	if err == io.EOF {
		err = termiosSersTimeout{}
	}
	return n, err
}

func (b *baseport) Close() error {
	return b.f.Close()
}

func (bp *baseport) Write(b []byte) (int, error) {
	return bp.f.Write(b)
}

func (bp *baseport) SetMode(baudrate, databits, parity, stopbits, handshake int) error {
	if baudrate <= 0 {
		return &ParameterError{"baudrate", "has to be > 0"}
	}

	var datamask uint
	switch databits {
	case 5:
		datamask = unix.CS5
	case 6:
		datamask = unix.CS6
	case 7:
		datamask = unix.CS7
	case 8:
		datamask = unix.CS8
	default:
		return &ParameterError{"databits", "has to be 5, 6, 7 or 8"}
	}

	if stopbits != 1 && stopbits != 2 {
		return &ParameterError{"stopbits", "has to be 1 or 2"}
	}
	var stopmask uint
	if stopbits == 2 {
		stopmask = unix.CSTOPB
	}

	var parmask uint
	switch parity {
	case N:
		parmask = 0
	case E:
		parmask = unix.PARENB
	case O:
		parmask = unix.PARENB | unix.PARODD
	default:
		return &ParameterError{"parity", "has to be N, E or O"}
	}

	var flowmask uint
	switch handshake {
	case NO_HANDSHAKE:
		flowmask = 0
	case RTSCTS_HANDSHAKE:
		flowmask = unix.CRTSCTS
	default:
		return &ParameterError{"handshake", "has to be NO_HANDSHAKE or RTSCTS_HANDSHAKE"}
	}

	tio, err := bp.getattr()
	if err != nil {
		return &Error{"getattr", err}
	}

	tio.Cflag &^= unix.CSIZE
	tio.Cflag |= cflagtype(datamask)

	tio.Cflag &^= unix.PARENB | unix.PARODD
	tio.Cflag |= cflagtype(parmask)

	tio.Cflag &^= unix.CSTOPB
	tio.Cflag |= cflagtype(stopmask)

	tio.Cflag &^= unix.CRTSCTS
	tio.Cflag |= cflagtype(flowmask)

	err = bp.fillBaudrate(tio, baudrate)
	if err != nil {
		return err
	}

	if err := bp.setattr(tio); err != nil {
		return &Error{"setattr", err}
	}

	/*if err := bp.SetBaudRate(baudrate); err != nil {
		return err
	}*/

	return nil
}

func (bp *baseport) GetMode() (mode Mode, err error) {
	var tio *unix.Termios
	tio, err = bp.getattr()
	if err != nil {
		return
	}

	tioCharSize := tio.Cflag & unix.CSIZE
	switch tioCharSize {
	case unix.CS5:
		mode.DataBits = 5
	case unix.CS6:
		mode.DataBits = 6
	case unix.CS7:
		mode.DataBits = 7
	case unix.CS8:
		mode.DataBits = 8
	default:
		err = fmt.Errorf("unknown character size field (%#08x) in termios", tioCharSize)
	}

	mode.Stopbits = 1
	if tio.Cflag&unix.CSTOPB != 0 {
		mode.Stopbits = 2
	}

	mode.Parity = N
	switch tio.Cflag & (unix.PARENB | unix.PARODD) {
	case unix.PARENB | unix.PARODD:
		mode.Parity = O
	case unix.PARENB:
		mode.Parity = E
	}

	mode.Handshake = NO_HANDSHAKE
	if tio.Cflag&unix.CRTSCTS != 0 {
		mode.Handshake = RTSCTS_HANDSHAKE
	}

	/*mode.Baudrate, err = bp.getBaudrate()
	if err != nil {
		return
	}*/
	//panic("baud rate getting not yet supported")
	mode.Baudrate, err = bp.extractBaudrate(tio)
	if err != nil {
		return
	}

	return
}

func (bp *baseport) SetReadParams(minread int, timeout float64) error {
	inttimeout := int(timeout * 10)
	if inttimeout < 0 {
		return &ParameterError{"timeout", "needs to be 0 or higher"}
	}
	// if a timeout is desired but too small for the termios timeout
	// granularity, set the minimum timeout
	if timeout > 0 && inttimeout == 0 {
		inttimeout = 1
	}

	tio, err := bp.getattr()
	if err != nil {
		return &Error{"getattr", err}
	}

	tio.Cc[unix.VMIN] = cctype(minread)
	tio.Cc[unix.VTIME] = cctype(inttimeout)

	//fmt.Printf("baud rates from termios: %d, %d\n", tio.c_ispeed, tio.c_ospeed)

	err = bp.setattr(tio)
	if err != nil {
		return &Error{"setattr", err}
	}

	return nil
}

func (bp *baseport) SetBreak(on bool) error {
	var (
		op       uint   = unix.TIOCCBRK
		opstring string = "setting break"
	)
	if on {
		op, opstring = unix.TIOCSBRK, "clearing break"
	}

	var onint int = 0
	if on {
		onint = 1
	}

	err := unix.IoctlSetInt(bp.fd, op, onint)

	if err != nil {
		return &Error{fmt.Sprintf("ioctl: %s", opstring), err}
	}

	return nil
}

func Open(fn string) (SerialPort, error) {
	// the order of system calls is taken from Apple's SerialPortSample
	// open the TTY device read/write, nonblocking, i.e. not waiting
	// for the CARRIER signal and without the TTY controlling the process
	fd, err := syscall.Open(fn, syscall.O_RDWR|
		syscall.O_NOCTTY|
		syscall.O_NONBLOCK,
		0666)
	if err != nil {
		return nil, fmt.Errorf("open %s: %v", fn, err)
	}

	s, err := takeOverFD(fd, fn)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (bp *baseport) SetDeadline(t time.Time) error      { return bp.f.SetDeadline(t) }
func (bp *baseport) SetReadDeadline(t time.Time) error  { return bp.f.SetReadDeadline(t) }
func (bp *baseport) SetWriteDeadline(t time.Time) error { return bp.f.SetWriteDeadline(t) }

type termiosSersTimeout struct{}

func (tst termiosSersTimeout) Error() string {
	return "timeout"
}

func (tst termiosSersTimeout) Timeout() bool {
	return true
}

func cfmakeraw(tio *unix.Termios) {
	tio.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP |
		unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	tio.Oflag &^= unix.OPOST
	tio.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	tio.Cflag &^= (unix.CSIZE | unix.PARENB)
	tio.Cflag |= unix.CS8
}

// the baud rate enum functions can be shared across termios platforms.
// they might be unused if a decent method for setting nontraditional baudrates
// is available.

func lookupbaudrate(br int) (cflagtype, error) {
	switch br {
	case 50:
		return unix.B50, nil
	case 75:
		return unix.B75, nil
	case 110:
		return unix.B110, nil
	case 134:
		return unix.B134, nil
	case 150:
		return unix.B150, nil
	case 200:
		return unix.B200, nil
	case 300:
		return unix.B300, nil
	case 600:
		return unix.B600, nil
	case 1200:
		return unix.B1200, nil
	case 1800:
		return unix.B1800, nil
	case 2400:
		return unix.B2400, nil
	case 4800:
		return unix.B4800, nil
	case 9600:
		return unix.B9600, nil
	case 19200:
		return unix.B19200, nil
	case 38400:
		return unix.B38400, nil
	case 57600:
		return unix.B57600, nil
	case 115200:
		return unix.B115200, nil
	case 230400:
		return unix.B230400, nil
	case 460800:
		return unix.B460800, nil
	case 500000:
		return unix.B500000, nil
	case 576000:
		return unix.B576000, nil
	case 921600:
		return unix.B921600, nil
	case 1000000:
		return unix.B1000000, nil
	case 1152000:
		return unix.B1152000, nil
	case 1500000:
		return unix.B1500000, nil
	case 2000000:
		return unix.B2000000, nil
	case 2500000:
		return unix.B2500000, nil
	case 3000000:
		return unix.B3000000, nil
	case 3500000:
		return unix.B3500000, nil
	case 4000000:
		return unix.B4000000, nil
	}

	return 0, fmt.Errorf("unsupported baud rate %d (only traditional baud rates allowed)", br)
}

func speedtobaudrate(cflag cflagtype) (int, error) {
	switch cflag & unix.CBAUD {
	case unix.B50:
		return 50, nil
	case unix.B75:
		return 75, nil
	case unix.B110:
		return 110, nil
	case unix.B134:
		return 134, nil
	case unix.B150:
		return 150, nil
	case unix.B200:
		return 200, nil
	case unix.B300:
		return 300, nil
	case unix.B600:
		return 600, nil
	case unix.B1200:
		return 1200, nil
	case unix.B1800:
		return 1800, nil
	case unix.B2400:
		return 2400, nil
	case unix.B4800:
		return 4800, nil
	case unix.B9600:
		return 9600, nil
	case unix.B19200:
		return 19200, nil
	case unix.B38400:
		return 38400, nil
	case unix.B57600:
		return 57600, nil
	case unix.B115200:
		return 115200, nil
	case unix.B230400:
		return 230400, nil
	case unix.B460800:
		return 460800, nil
	case unix.B500000:
		return 500000, nil
	case unix.B576000:
		return 576000, nil
	case unix.B921600:
		return 921600, nil
	case unix.B1000000:
		return 1000000, nil
	case unix.B1152000:
		return 1152000, nil
	case unix.B1500000:
		return 1500000, nil
	case unix.B2000000:
		return 2000000, nil
	case unix.B2500000:
		return 2500000, nil
	case unix.B3000000:
		return 3000000, nil
	case unix.B3500000:
		return 3500000, nil
	case unix.B4000000:
		return 4000000, nil
	}

	// if none of the above matched but we have CBAUDEX set
	if cflag&unix.CBAUDEX != 0 {
		return 0, fmt.Errorf("cannot read extended baud rate setting (sers without termios2 support)")
	}

	return 0, fmt.Errorf("unknown baud rate encoding (no enum for 0x%08x)", cflag)
}
