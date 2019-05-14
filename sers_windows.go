// +build windows

package sers

// taken from https://github.com/tarm/goserial
// and slightly modified

// (C) 2011, 2012 Tarmigan Casebolt, Benjamin Siegert, Michael Meier
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type serialPort struct {
	f  *os.File
	fd syscall.Handle
	rl sync.Mutex
	wl sync.Mutex
}

type structDCB struct {
	DCBlength, BaudRate                            uint32
	flags                                          [4]byte
	wReserved, XonLim, XoffLim                     uint16
	ByteSize, Parity, StopBits                     byte
	XonChar, XoffChar, ErrorChar, EofChar, EvtChar byte
	wReserved1                                     uint16
}

type structTimeouts struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

type opError struct {
	op       string
	filename string
	err      error
}

func (o opError) Error() string {
	return o.op + " " + o.filename + ": " + o.err.Error()
}

//func openPort(name string) (rwc io.ReadWriteCloser, err error) { // TODO
func Open(name string) (rwc SerialPort, err error) {
	if len(name) > 0 && name[0] != '\\' {
		name = "\\\\.\\" + name
	}

	h, err := syscall.CreateFile(syscall.StringToUTF16Ptr(name),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL|syscall.FILE_FLAG_OVERLAPPED,
		0)
	if err != nil {
		return nil, opError{"open", name, err}
	}
	f := os.NewFile(uintptr(h), name)
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	/*if err = setCommState(h, baud); err != nil {
		return
	}*/
	if err = setupComm(h, 64, 64); err != nil {
		return
	}
	if err = setCommTimeouts(h, 0.0); err != nil {
		return
	}
	if err = setCommMask(h); err != nil {
		return
	}

	port := new(serialPort)
	port.f = f
	port.fd = h

	return port, nil
}

func (p *serialPort) Close() error {
	return p.f.Close()
}

func (p *serialPort) Write(buf []byte) (int, error) {
	p.wl.Lock()
	defer p.wl.Unlock()

	return p.f.Write(buf)
}

func (p *serialPort) Read(buf []byte) (int, error) {
	if p == nil || p.f == nil {
		return 0, fmt.Errorf("Invalid port on read %v %v", p, p.f)
	}

	p.rl.Lock()
	defer p.rl.Unlock()

	return p.f.Read(buf)
}

func (p *serialPort) SetBreak(on bool) error {
	var opstring string = "ClearCommBreak"
	if on {
		opstring = "SetCommBreak"
	}

	var (
		r1  uintptr
		err error
	)
	if on {
		r1, _, err = syscall.Syscall(nSetCommBreak, 1, uintptr(p.fd), 0, 0)
	} else {
		r1, _, err = syscall.Syscall(nClearCommBreak, 1, uintptr(p.fd), 0, 0)
	}
	if r1 == 0 {
		return &Error{opstring, err}
	}

	return nil
}

func (p *serialPort) SetDeadline(t time.Time) error      { return p.f.SetDeadline(t) }
func (p *serialPort) SetReadDeadline(t time.Time) error  { return p.f.SetReadDeadline(t) }
func (p *serialPort) SetWriteDeadline(t time.Time) error { return p.f.SetWriteDeadline(t) }

var (
	nSetCommState,
	nSetCommTimeouts,
	nSetCommMask,
	nSetupComm,
	nCreateEvent,
	nSetCommBreak,
	nClearCommBreak uintptr
)

func init() {
	k32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		panic("LoadLibrary " + err.Error())
	}
	defer syscall.FreeLibrary(k32)

	nSetCommState = getProcAddr(k32, "SetCommState")
	nSetCommTimeouts = getProcAddr(k32, "SetCommTimeouts")
	nSetCommMask = getProcAddr(k32, "SetCommMask")
	nSetupComm = getProcAddr(k32, "SetupComm")
	nCreateEvent = getProcAddr(k32, "CreateEventW")
	nSetCommBreak = getProcAddr(k32, "SetCommBreak")
	nClearCommBreak = getProcAddr(k32, "ClearCommBreak")
}

func getProcAddr(lib syscall.Handle, name string) uintptr {
	addr, err := syscall.GetProcAddress(lib, name)
	if err != nil {
		panic(name + " " + err.Error())
	}
	return addr
}

func setCommState(h syscall.Handle, baud, databits, parity, handshake int) error {
	var params structDCB
	params.DCBlength = uint32(unsafe.Sizeof(params))

	params.flags[0] = 0x01  // fBinary
	params.flags[0] |= 0x10 // Assert DSR

	params.ByteSize = byte(databits)

	params.BaudRate = uint32(baud)
	//params.ByteSize = 8

	switch parity {
	case N:
		params.flags[0] &^= 0x02
		params.Parity = 0 // NOPARITY
	case E:
		params.flags[0] |= 0x02
		params.Parity = 2 // EVENPARITY
	case O:
		params.flags[0] |= 0x02
		params.Parity = 1 // ODDPARITY
	default:
		return StringError("invalid parity setting")
	}

	switch handshake {
	case NO_HANDSHAKE:
		// TODO: reset handshake
	default:
		return StringError("only NO_HANDSHAKE is supported on windows")
	}

	r, _, err := syscall.Syscall(nSetCommState, 2, uintptr(h), uintptr(unsafe.Pointer(&params)), 0)
	if r == 0 {
		return err
	}
	return nil
}

func setCommTimeouts(h syscall.Handle, constTimeout float64) error {
	var timeouts structTimeouts
	const MAXDWORD = 1<<32 - 1
	timeouts.ReadIntervalTimeout = MAXDWORD
	timeouts.ReadTotalTimeoutMultiplier = MAXDWORD
	//timeouts.ReadTotalTimeoutConstant = MAXDWORD - 1
	if constTimeout == 0 {
		timeouts.ReadTotalTimeoutConstant = MAXDWORD - 1
	} else {
		timeouts.ReadTotalTimeoutConstant = uint32(constTimeout * 1000.0)
	}

	/* From http://msdn.microsoft.com/en-us/library/aa363190(v=VS.85).aspx

		 For blocking I/O see below:

		 Remarks:

		 If an application sets ReadIntervalTimeout and
		 ReadTotalTimeoutMultiplier to MAXDWORD and sets
		 ReadTotalTimeoutConstant to a value greater than zero and
		 less than MAXDWORD, one of the following occurs when the
		 ReadFile function is called:

		 If there are any bytes in the input buffer, ReadFile returns
		       immediately with the bytes in the buffer.

		 If there are no bytes in the input buffer, ReadFile waits
	               until a byte arrives and then returns immediately.

		 If no bytes arrive within the time specified by
		       ReadTotalTimeoutConstant, ReadFile times out.
	*/

	r, _, err := syscall.Syscall(nSetCommTimeouts, 2, uintptr(h), uintptr(unsafe.Pointer(&timeouts)), 0)
	if r == 0 {
		return err
	}
	return nil
}

func setupComm(h syscall.Handle, in, out int) error {
	r, _, err := syscall.Syscall(nSetupComm, 3, uintptr(h), uintptr(in), uintptr(out))
	if r == 0 {
		return err
	}
	return nil
}

func setCommMask(h syscall.Handle) error {
	const EV_RXCHAR = 0x0001
	r, _, err := syscall.Syscall(nSetCommMask, 2, uintptr(h), EV_RXCHAR, 0)
	if r == 0 {
		return err
	}
	return nil
}

func (sp *serialPort) SetMode(baudrate, databits, parity, stopbits, handshake int) error {
	if err := setCommState(sp.fd, baudrate, databits, parity, handshake); err != nil {
		return err
	}
	//return StringError("SetMode not implemented yet on Windows")
	return nil
}

func (sp *serialPort) SetReadParams(minread int, timeout float64) error {
	// TODO: minread is ignored!
	return setCommTimeouts(sp.fd, timeout)
}

type winSersTimeout struct{}

func (wst winSersTimeout) Error() string {
	return "a timeout has occured"
}

func (wst winSersTimeout) Timeout() bool {
	return true
}
