// +build linux
// +build sers_termios_enum !386,!amd64,!arm64,!arm,!mips64,!mips64le,!mips,!mipsle,!riscv64,!s390x,!sparc64

package sers

import (
	"golang.org/x/sys/unix"
)

/*

This files uses struct termios to set baudrates. It does not care whether the
struct termios2 fields ispeed and ospeed are present. Instead, it uses the
baudrate enums.

*/

// when updating types: they have to match in sers_linux_termio2.go
type termiosPlatformData struct{}
type cflagtype = uint32
type cctype = uint8

func (bp *baseport) getattr() (*unix.Termios, error) {
	return unix.IoctlGetTermios(bp.fd, unix.TCGETS)
}

func (bp *baseport) setattr(tio *unix.Termios) error {
	return unix.IoctlSetTermios(bp.fd, unix.TCSETS, tio)
}

func (bp *baseport) fillBaudrate(tio *unix.Termios, baudrate int) error {
	orval, err := lookupbaudrate(baudrate)
	if err != nil {
		return err
	}

	tio.Cflag &^= unix.CBAUD
	tio.Cflag |= orval

	return nil
}

func (bp *baseport) extractBaudrate(tio *unix.Termios) (int, error) {
	return speedtobaudrate(tio.Cflag)
}
