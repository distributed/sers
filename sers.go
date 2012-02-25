package sers

import (
	"fmt"
	"io"
)

type SerialPort interface {
	io.Reader
	io.Writer
	io.Closer
	SetBaudRate(br int) error
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
