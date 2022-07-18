package main

import (
	"flag"
	"fmt"
	"github.com/distributed/sers"
	"log"
	"os"
	"time"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

type DeadlineSerialPort interface {
	sers.SerialPort
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

func Main() error {
	flag.Parse()

	if flag.NArg() < 1 {
		return fmt.Errorf("please specify serial port file name")
	}

	sp, err := sers.Open(flag.Arg(0))
	if err != nil {
		return fmt.Errorf("error opening port: %w", err)
	}
	defer sp.Close()

	fmt.Printf("opened %q\n", flag.Arg(0))

	dsp, gotdsp := sp.(DeadlineSerialPort)

	if !gotdsp {
		return fmt.Errorf("serial port does not conform to deadline interface")
	}

	fmt.Printf("got a deadline capable serial port struct\n")
	realsoon := time.Now().Add(50 * time.Millisecond)
	err = dsp.SetReadDeadline(realsoon)
	if err != nil {
		return fmt.Errorf("unable to set read deadline: %w", err)
	}

	var buf [16]byte
	_, err = dsp.Read(buf[:])
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Printf("got timeout error as expected\n")
			return nil
		} else {
			return fmt.Errorf("got unexpected error: %w")
		}
	}
	return fmt.Errorf("did not receive read error, expected a timeout error")
}
