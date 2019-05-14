package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/distributed/sers"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

var baudrate = flag.Int("baudrate", 0, "baudrate to be used, 0 means no change on opened file. if changed, format is 8N1")
var sendstring = flag.String("sendstring", "", "send this string before reading")
var timeout = flag.Duration("timeout", 1*time.Second, "deadline is time.Now.Add(thisParameter)")

func Main() error {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return fmt.Errorf("please provide a serial file name")
	} else if len(args) > 1 {
		return fmt.Errorf("superfluous arguments: %v", args[1:])
	}

	if *timeout < 0 {
		return fmt.Errorf("-timeout has to be nonnegative")
	}

	var sp sers.SerialPort
	var err error

	serfn := args[0]
	if !*takeover {
		sp, err = sers.Open(serfn)
		if err != nil {
			return err
		}
	} else {
		f, err := os.Open(serfn)
		if err != nil {
			return err
		}

		sp, err = takeoverFile(f)
		if err != nil {
			return err
		}
	}
	defer sp.Close()

	if *baudrate != 0 {
		err = sp.SetMode(*baudrate, 8, sers.N, 1, sers.NO_HANDSHAKE)
		if err != nil {
			return err
		}
	}

	if len(*sendstring) > 0 {
		_, err := io.WriteString(sp, *sendstring)
		if err != nil {
			return err
		}
	}

	err = sp.SetReadDeadline(time.Now().Add(*timeout))
	if err != nil {
		return err
	}
	start := time.Now()

	go func() {
		hardtimeout := (*timeout) + 500*time.Millisecond
		time.Sleep(hardtimeout)
		panic("program has not finished after 500 ms grace period over read deadline")
	}()

	var buf [128]byte
	for {
		n, err := sp.Read(buf[:])
		b := buf[:n]
		if n > 0 {
			fmt.Printf("data: % 02x\n", b)
		}
		if err != nil {
			if os.IsTimeout(err) {
				end := time.Now()
				fmt.Printf("got timeout after %v\n", end.Sub(start))
				return nil
			}

			return err
		}
	}

	return nil
}
