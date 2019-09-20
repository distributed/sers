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
	"sync"
	"unsafe"
)

const (
	// using C.IOSSIOSPEED yields 0x80085402
	// which does not work. don't ask me why
	// this define is wrong in cgo.
	IOSSIOSPEED = 0x80045402
)

type termiosPlatformData struct {
	lock        sync.Mutex
	baudrateSet bool
	baudrate    int
}

func (bp *baseport) SetBaudRate(br int) error {
	var speed C.speed_t = C.speed_t(br)

	//fmt.Printf("C.IOSSIOSPEED %x\n", uint64(C.IOSSIOSPEED))
	//fmt.Printf("for file %v, fd %d\n", bp.f, bp.fd)

	ret, err := C.ioctl1(C.int(bp.fd), C.uint(IOSSIOSPEED), unsafe.Pointer(&speed))
	if ret == -1 || err != nil {
		return &Error{"setting baud rate: ioctl", err}
	}

	bp.platformData.lock.Lock()
	bp.platformData.baudrate = br
	bp.platformData.baudrateSet = true
	bp.platformData.lock.Unlock()

	return nil
}

type osxBaudrateRetrievalFailed struct{}

func (osxBaudrateRetrievalFailed) Error() string {
	return "sers: Cannot get current baud rate setting. You have to set the baudrate through SetMode before you can read it via GetMode. See documentation."
}

func (bp *baseport) getBaudrate() (int, error) {
	bp.platformData.lock.Lock()
	defer bp.platformData.lock.Unlock()

	if bp.platformData.baudrateSet {
		return bp.platformData.baudrate, nil
	}

	return -1, osxBaudrateRetrievalFailed{}
}
