package sers

import (
	"io"
)


type SerialPort interface {
	io.Reader
	io.Writer
	io.Closer
}
