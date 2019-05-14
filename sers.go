// Copyright 2012 Michael Meier. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package sers offers serial port access. It is a stated goal of this
// package to allow for non-standard bit rates as the may be useful
// in a wide range of embedded projects.
package sers

import (
	"fmt"
	"io"
	"time"
)

const (
	N = 0 // no parity
	E = 1 // even parity
	O = 2 // odd parity
)

const (
	NO_HANDSHAKE     = 0
	RTSCTS_HANDSHAKE = 1
)

// Serialport represents a serial port and offers configuration of baud
// rate, frame format, handshaking and read paramters as well as setting and
// clearing break conditions.
type SerialPort interface {
	io.Reader
	io.Writer
	io.Closer

	// SetMode sets the frame format and handshaking configuration.
	// baudrate may be freely chosen, the driver is allowed to reject
	// unachievable baud rates. databits may be any number of data bits
	// supported by the driver. parity is one of (N|O|E) for none, odd
	// or even parity. handshake is either NO_HANDSHAKE or
	// RTSCTS_HANDSHAKE.
	SetMode(baudrate, databits, parity, stopbits, handshake int) error

	// SetReadParams sets the minimum number of bits to read and a read
	// timeout in seconds. These parameters roughly correspond to the
	// UNIX termios concepts of VMIN and VTIME.
	SetReadParams(minread int, timeout float64) error

	// SetBreak turns on the generation of a break condition if on == true,
	// otherwise it clear the break condition.
	SetBreak(on bool) error

	// SetDeadline sets the read and write deadlines associated
	// with the connection. It is equivalent to calling both
	// SetReadDeadline and SetWriteDeadline.
	//
	// A deadline is an absolute time after which I/O operations
	// fail with a timeout (see type Error) instead of
	// blocking. The deadline applies to all future and pending
	// I/O, not just the immediately following call to ReadFrom or
	// WriteTo. After a deadline has been exceeded, the connection
	// can be refreshed by setting a deadline in the future.
	//
	// An idle timeout can be implemented by repeatedly extending
	// the deadline after successful ReadFrom or WriteTo calls.
	//
	// A zero value for t means I/O operations will not time out.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for future ReadFrom calls
	// and any currently-blocked ReadFrom call.
	// A zero value for t means ReadFrom will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for future WriteTo calls
	// and any currently-blocked WriteTo call.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	// A zero value for t means WriteTo will not time out.
	SetWriteDeadline(t time.Time) error
}

type StringError string

func (se StringError) Error() string {
	return string(se)
}

type ParameterError struct {
	Parameter string
	Reason    string
}

func (pe *ParameterError) Error() string {
	return fmt.Sprintf("error in parameter '%s': %s")
}

type Error struct {
	Operation       string
	UnderlyingError error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Operation, e.UnderlyingError)
}
