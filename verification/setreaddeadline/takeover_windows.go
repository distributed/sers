// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/distributed/sers"
)

var takeoverhardcoded = false
var takeover = &takeoverhardcoded

func takeoverFile(f *os.File) (sers.SerialPort, error) {
	return nil, fmt.Errorf("TakeOver not implemented on Windows")
}
