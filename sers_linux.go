// +build linux

package sers

// Copyright 2012 Michael Meier. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/*
#include <termios.h>

 extern int setbaudrate(int fd, int br);
 extern int clearnonblocking(int fd);
*/
import "C"

func (bp *baseport) SetBaudRate(br int) error {
	// setting baud rate via new struct termios2 method 
	_, err := C.setbaudrate(C.int(bp.fd), C.int(br))
	if err != nil {
		return err
	}

	return nil
}
