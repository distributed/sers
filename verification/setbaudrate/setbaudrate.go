package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/distributed/sers"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

var sendU = flag.Bool("sendU", false, "set to continuously send U (0x55)")

func Main() error {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		return fmt.Errorf("please provide a serial file name and a baud rate")
	}

	sfn := args[0]
	baudratestring := args[1]

	baudrate64, err := strconv.ParseInt(baudratestring, 0, 32)
	if err != nil {
		return fmt.Errorf("error parsing baud rate: %v", err)
	}

	baudrate := int(baudrate64)

	sp, err := sers.Open(sfn)
	if err != nil {
		return err
	}

	err = sp.SetMode(baudrate, 8, sers.N, 1, sers.NO_HANDSHAKE)
	if err != nil {
		return fmt.Errorf("error setting baud rate: %v", err)
	}

	fmt.Printf("set baud rate of %q to %d baud\n", sfn, baudrate)

	if *sendU {
		buf := []byte(strings.Repeat("U", 128))
		fmt.Printf("entering loop to send %d byte blocks of Us\n", len(buf))

		for {
			_, err = sp.Write(buf)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
