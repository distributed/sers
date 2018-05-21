package main

import (
	"log"
	"github.com/distributed/sers"
	"flag"
	"time"
	"fmt"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

var ontime = flag.Duration("on",500*time.Millisecond,"on period")
var offtime = flag.Duration("off", 500*time.Millisecond, "off period")

func Main() error {
	flag.Parse()

	args:=flag.Args()
	if len(args) < 1 {
		return fmt.Errorf("please provide a serial file name")
	} else if len(args) > 1 {
		return fmt.Errorf("extraneous arguments")
	}

	fn:=args[0]

	sp,err:=sers.Open(fn)
	if err != nil {
		return err
	}

	for {
		fmt.Printf("setting break\n")
		err=sp.SetBreak(true)
		if err != nil {
			return err
		}

		time.Sleep(*ontime)

		fmt.Printf("clearing break\n")
		err = sp.SetBreak(false)
		if err != nil {
			return err
		}

		time.Sleep(*offtime)
	}

	return nil
}
