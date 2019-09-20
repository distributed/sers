// +build linux

package sers

// Copyright 2012 Michael Meier. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/*
#include <termios.h>

 extern int setbaudrate(int fd, int br);
 extern int clearnonblocking(int fd);
 extern int getbaudrate(int fd, int *br);
*/
import "C"
import "fmt"

type termiosPlatformData struct{}

func (bp *baseport) SetBaudRate(br int) error {
	// setting baud rate via new struct termios2 method
	_, err := C.setbaudrate(C.int(bp.fd), C.int(br))
	if err != nil {
		return err
	}

	return nil
}

func (bp *baseport) getBaudrate() (int, error) {
	var br C.int
	ret, errnoerr := C.getbaudrate(C.int(bp.fd), &br)
	if errnoerr != nil {
		return 0, fmt.Errorf("error getting baud rate: %v", errnoerr)
	}
	if ret != 0 {
		return 0, fmt.Errorf("error getting baud rate")
	}

	return int(br), nil
}
