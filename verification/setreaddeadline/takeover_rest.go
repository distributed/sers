// +build !windows

package main

import (
	"flag"
	"github.com/distributed/sers"
	"os"
)

func takeoverFile(f *os.File) (sers.SerialPort, error) { return sers.TakeOver(f) }

var takeover = flag.Bool("takeover", false, "if set, use sers.Takeover() instead of sers.Open(). use to purposefully break this verification step")
