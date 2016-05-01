package sers_test

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/distributed/sers"
)

// This program opens a serial port, configurable by changing portname
// below and configures it for 57600 baud, 8 data bits, no parity bit,
// 1 stop bit, no handshaking. It then reads up to 128 bytes from the serial
// port, with a timeout of 1 second and prints the received
// bytes to stdout.
func Example() {
	portname := "/dev/ttyUSB0"
	rb, err := readFirstBytesFromPort(portname)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("got %d bytes from %s:\n%s", len(rb), portname,
		hex.Dump(rb))
}

func readFirstBytesFromPort(fn string) ([]byte, error) {
	sp, err := sers.Open(fn)
	if err != nil {
		return nil, err
	}
	defer sp.Close()

	// 57600 baud, 8 data bits, no parity bit, 1 stop bit
	// no handshake. non-standard baud rates are possible.
	err = sp.SetMode(57600, 8, sers.N, 1, sers.NO_HANDSHAKE)
	if err != nil {
		return nil, err
	}

	// setting:
	// minread = 0: minimal buffering on read, return characters as early as possible
	// timeout = 1.0: time out if after 1.0 seconds nothing is received
	err = sp.SetReadParams(0, 1.0)
	if err != nil {
		return nil, err
	}

	var rb [128]byte
	n, err := sp.Read(rb[:])

	return rb[:n], err
}
