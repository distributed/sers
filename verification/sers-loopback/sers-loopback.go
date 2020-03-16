package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/distributed/sers"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

//var forwardStdin = flag.Bool("forwardstdin", false, "set to forward data read from stdin to serial file")
var initSleep = flag.Duration("initsleep", 0, "inital sleep before sending data")
var sendSleep = flag.Duration("sendsleep", 500*time.Millisecond, "time to sleep after send")

func Main() error {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		return fmt.Errorf("please specify a serial file name and a mode string")
	} else if len(args) > 2 {
		return fmt.Errorf("extraneous arguments")
	}

	fn := args[0]
	modestring := ""
	if len(args) > 1 {
		modestring = args[1]
	}

	var mode sers.Mode
	var err error

	if len(modestring) > 0 {
		mode, err = sers.ParseModestring(modestring)
		if err != nil {
			return err
		}
	}

	f, err := sers.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(modestring) > 0 {
		err = sers.SetModeStruct(f, mode)
		if err != nil {
			return err
		}
	}

	log.SetFlags(log.Flags() | log.Lmicroseconds)

	ctx := context.Background()
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		time.Sleep(*initSleep)
		n := 0
		for {
			bytes := []byte{0x21, 0x22, 0x23, byte(n)}
			n++

			log.Printf("sending % 02x\n", bytes)
			_, err := f.Write(bytes)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(*sendSleep):
			}
		}

	})

	eg.Go(func() error {
		buf := [128]byte{}
		for {
			n, err := f.Read(buf[:])
			if n > 0 {
				log.Printf("got %d bytes: % 02x\n", n, buf[:n])
			}
			if err != nil {
				return err
			}
		}
	})

	return eg.Wait()
}
