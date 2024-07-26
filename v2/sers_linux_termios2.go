// +build linux
// +build 386 amd64 arm64 arm mips64 mips64le mips mipsle riscv64 s390x sparc64
// +build !sers_termios_enum

package sers

/*

This file uses struct termios2 to set baudrates. It uses struct
unix.Termios, which includes fields added by struct termios2, namely ispeed
and ospeed.

You can force using the old style method by setting the build tag sers_termios_enum.

The build tags on the second line serve to whitelist the platforms on which the
termios2 fields are defined. When adding a platform to the whitelist, care has
to be taken to blacklist the platform in sers_linux_termios1.go.

The following command, when run in golang.org/x/sys/unix will show which
platforms support struct termios2:

$ grep -e TCGETS -r * --include "*linux*"
zerrors_linux_386.go:	TCGETS                               = 0x5401
zerrors_linux_386.go:	TCGETS2                              = 0x802c542a
zerrors_linux_amd64.go:	TCGETS                               = 0x5401
zerrors_linux_amd64.go:	TCGETS2                              = 0x802c542a
zerrors_linux_arm64.go:	TCGETS                               = 0x5401
zerrors_linux_arm64.go:	TCGETS2                              = 0x802c542a
zerrors_linux_arm.go:	TCGETS                               = 0x5401
zerrors_linux_arm.go:	TCGETS2                              = 0x802c542a
zerrors_linux_mips64.go:	TCGETS                               = 0x540d
zerrors_linux_mips64.go:	TCGETS2                              = 0x4030542a
zerrors_linux_mips64le.go:	TCGETS                               = 0x540d
zerrors_linux_mips64le.go:	TCGETS2                              = 0x4030542a
zerrors_linux_mips.go:	TCGETS                               = 0x540d
zerrors_linux_mips.go:	TCGETS2                              = 0x4030542a
zerrors_linux_mipsle.go:	TCGETS                               = 0x540d
zerrors_linux_mipsle.go:	TCGETS2                              = 0x4030542a
zerrors_linux_ppc64.go:	TCGETS                               = 0x402c7413
zerrors_linux_ppc64le.go:	TCGETS                               = 0x402c7413
zerrors_linux_riscv64.go:	TCGETS                               = 0x5401
zerrors_linux_riscv64.go:	TCGETS2                              = 0x802c542a
zerrors_linux_s390x.go:	TCGETS                               = 0x5401
zerrors_linux_s390x.go:	TCGETS2                              = 0x802c542a
zerrors_linux_sparc64.go:	TCGETS                               = 0x40245408
zerrors_linux_sparc64.go:	TCGETS2                              = 0x402c540c

(run against 2837fb4f24fe)

*/

import (
	"golang.org/x/sys/unix"
)

// when updating types: they have to match in sers_linux_termio1.go
type termiosPlatformData struct{}
type cflagtype = uint32
type cctype = uint8

func (bp *baseport) getattr() (*unix.Termios, error) {
	return unix.IoctlGetTermios(bp.fd, unix.TCGETS2)
}

func (bp *baseport) setattr(tio *unix.Termios) error {
	return unix.IoctlSetTermios(bp.fd, unix.TCSETS2, tio)
}

func (bp *baseport) fillBaudrate(tio *unix.Termios, baudrate int) error {
	tio.Cflag &^= unix.CBAUD
	tio.Cflag |= unix.CBAUDEX
	tio.Ispeed = uint32(baudrate)
	tio.Ospeed = tio.Ispeed
	return nil
}

func (bp *baseport) extractBaudrate(tio *unix.Termios) (int, error) {
	return int(tio.Ispeed), nil
}
