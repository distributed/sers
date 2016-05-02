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
	// setting aliased baud rate
	_, err := C.setbaudrate(C.int(bp.f.Fd()), C.int(br))
	if err != nil {
		return err
	}

	tio, err := bp.getattr()
	if err != nil {
		return err
	}

	// using aliased baudrate
	//C.cfsetspeed(tio, C.B38400)

	err = bp.setattr(tio)
	if err != nil {
		return err
	}

	return nil
}

func (bp *baseport) ClearNonBlocking() error {
	_, err := C.clearnonblocking(C.int(bp.f.Fd()))
	return err
}
