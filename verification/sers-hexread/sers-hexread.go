package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/distributed/sers"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

var forwardStdin = flag.Bool("forwardstdin", false, "set to forward data read from stdin to serial file")

func Main() error {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		return fmt.Errorf("please specify at least a serial file name")
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

	ctx := context.Background()
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if *forwardStdin {
		eg.Go(func() error {
			in := os.Stdin

			var buf [128]byte
			for {
				n, err := in.Read(buf[:])
				fmt.Printf("read %d bytes (% 02x) from stdin, forwarding to serial port...\n", n, buf[:n])
				if n > 0 {
					_, err := f.Write(buf[:n])
					if err != nil {
						return err
					}
					fmt.Printf("wrote %d bytes to serial port.\n", n)
				}
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	eg.Go(func() error {
		buf := [128]byte{}
		for {
			n, err := f.Read(buf[:])
			fmt.Println(n, err)
			if n > 0 {
				fmt.Printf("got %d bytes: % 02x\n", n, buf[:n])
			}
			if err != nil {
				return err
			}
		}
	})

	return eg.Wait()
}
