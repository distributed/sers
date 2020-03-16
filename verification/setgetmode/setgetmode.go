package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/distributed/sers"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return fmt.Errorf("please provide a file name to act upon")
	} else if len(args) > 2 {
		return fmt.Errorf("extraneous arguments: %q", args[1:])
	}
	modestring := "9600,8o1,rtscts"
	doset := false
	if len(args) >= 2 {
		modestring = args[1]
		doset = true
	}

	setmode, err := sers.ParseModestring(modestring)
	if err != nil {
		return err
	}

	fn := args[0]

	sp, err := sers.Open(fn)
	if err != nil {
		return err
	}
	defer sp.Close()

	if doset {
		err = sers.SetModeStruct(sp, setmode)
		if err != nil {
			return err
		}
		fmt.Printf("set %q\n", setmode)
	}

	getmode, err := sp.GetMode()
	if err != nil {
		return err
	}

	fmt.Printf("read mode: %v\n", getmode)

	return nil
}
