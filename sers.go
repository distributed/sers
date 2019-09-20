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
}

func SetModeStruct(sp SerialPort, mode Mode) error {
	return sp.SetMode(mode.Baudrate, mode.DataBits, mode.Parity, mode.Stopbits, mode.Handshake)

}

type Mode struct {
	Baudrate  int
	DataBits  int
	Parity    int
	Stopbits  int
	Handshake int
}

func (m Mode) Valid() bool {
	if m.Baudrate < 0 {
		return false
	}
	if m.DataBits < 5 || m.DataBits > 8 {
		return false
	}
	if !(m.Parity == N || m.Parity == O || m.Parity == E) {
		return false
	}
	if !(m.Stopbits == 1 || m.Stopbits == 2) {
		return false
	}
	if !(m.Handshake == NO_HANDSHAKE || m.Handshake == RTSCTS_HANDSHAKE) {
		return false
	}

	return true
}

func (m Mode) String() string {
	if !m.Valid() {
		return fmt.Sprintf("invalid_mode(%d,%d,%d,%d,%d)",
			m.Baudrate,
			m.DataBits,
			m.Parity,
			m.Stopbits,
			m.Handshake)
	}

	parstring := ""
	switch m.Parity {
	case N:
		parstring = "n"
	case O:
		parstring = "o"
	case E:
		parstring = "e"
	default:
		panic("unhandled parity setting")
	}

	hsstring := ""
	switch m.Handshake {
	case NO_HANDSHAKE:
		hsstring = "none"
	case RTSCTS_HANDSHAKE:
		hsstring = "rtscts"
	default:
		panic("unhandled handshake setting")
	}

	return fmt.Sprintf("%d,%d%s%d,%s",
		m.Baudrate,
		m.DataBits,
		parstring,
		m.Stopbits,
		hsstring)
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
