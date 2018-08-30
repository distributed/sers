// +build darwin

package sers

// Copyright 2012 Michael Meier. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/*#include <sys/ioctl.h>
#include <IOKit/serial/ioss.h>

#include <sys/types.h>
#include <unistd.h>
#include <fcntl.h>

 extern int fcntl1(int i, unsigned int r, unsigned int d);
 extern int ioctl1(int i, unsigned int r, void *d);*/
import "C"

import (
	"unsafe"
)

func (bp *baseport) SetBaudRate(br int) error {
	var speed C.speed_t = C.speed_t(br)

	//fmt.Printf("C.IOSSIOSPEED %x\n", uint64(C.IOSSIOSPEED))
	//fmt.Printf("for file %v, fd %d\n", bp.f, bp.fd)

	ret, err := C.ioctl1(C.int(bp.fd), C.uint(IOSSIOSPEED), unsafe.Pointer(&speed))
	if ret == -1 || err != nil {
		return &Error{"setting baud rate: ioctl", err}
	}

	return nil
}
