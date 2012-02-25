package sers

/*#include <stddef.h>
#include <stdlib.h>
#include <termios.h>
#include <sys/ioctl.h>
#include <IOKit/serial/ioss.h> 

#include <sys/types.h>
#include <unistd.h>
#include <fcntl.h>

 extern int ioctl1(int i, unsigned int r, void *d);
 extern int fcntl1(int i, unsigned int r, unsigned int d);
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	// using C.IOSSIOSPEED yields 0x80085402 
	// which does not work. don't ask me why
	// this define is wrong in cgo.
	IOSSIOSPEED = 0x80045402 
)

type baseport struct {
	f *os.File
}

func TakeOver(f *os.File) (SerialPort, error) {
	if f == nil {
		return nil, &ParameterError{"f", "needs to be non-nil"}
	}
	bp := &baseport{f}

	tio, err := bp.getattr()
	if err != nil {
		return nil, &Error{"putting fd in non-canonical mode", err}
	}

	C.cfmakeraw(tio)

	err = bp.setattr(tio)
	if err != nil {
		return nil, &Error{"putting fd in non-canonical mode", err}
	}

	return bp, nil
}

func (bp *baseport) Read(b []byte) (int, error) {
	return bp.f.Read(b)
}

func (b *baseport) Close() error {
	return b.f.Close()
}

func (bp *baseport) Write(b []byte) (int, error) {
	return bp.f.Write(b)
}

func (bp *baseport) getattr() (*C.struct_termios, error) {
	var tio C.struct_termios
	res, err := C.tcgetattr(C.int(bp.f.Fd()), (*C.struct_termios)(unsafe.Pointer(&tio)))
	if res != 0 || err != nil {
		return nil, err
	}

	return &tio, nil
}

func (bp *baseport) setattr(tio *C.struct_termios) error {
	res, err := C.tcsetattr(C.int(bp.f.Fd()), C.TCSANOW, (*C.struct_termios)(unsafe.Pointer(tio)))
	if res != 0 || err != nil {
		return err
	}

	return nil
}

func (bp *baseport) SetBaudRate(br int) error {
	var speed C.speed_t = C.speed_t(br)

	fmt.Printf("C.IOSSIOSPEED %x\n", uint64(C.IOSSIOSPEED))
	fmt.Printf("for file %v, fd %d\n", bp.f, bp.f.Fd())

	ret, err := C.ioctl1(C.int(bp.f.Fd()), C.uint(IOSSIOSPEED), unsafe.Pointer(&speed))
	if ret == -1 || err != nil {
		return err
	}

	return nil
}

func Open(fn string) (SerialPort, error) {
	// the order of system calls is taken from Apple's SerialPortSample
	// open the TTY device read/write, nonblocking, i.e. not waiting
	// for the CARRIER signal and without the TTY controlling the process
	f, err := os.OpenFile(fn, syscall.O_RDWR|
		syscall.O_NONBLOCK|
		syscall.O_NOCTTY, 0666)
	if err != nil {
		return nil, err
	}

	// clear non-blocking mode
	res, err := C.fcntl1(C.int(f.Fd()), C.F_SETFL, 0)
	if res < 0 || err != nil {
		f.Close()
		return nil, &Error{"putting fd into non-blocking mode", err}
	}

	s, err := TakeOver(f)
	if err != nil {
		return nil, err
	}

	return s, nil
}
