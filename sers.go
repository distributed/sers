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

type SerialPort interface {
	io.Reader
	io.Writer
	io.Closer
	SetMode(baudrate, databits, parity, stopbits, handshake int) error
	SetReadParams(minread int, timeout float64) error
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
