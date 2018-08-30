package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/distributed/sers"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}

func readpart(r io.ReadCloser) error {
	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("============================> close now\n")
		err := r.Close()
		fmt.Printf("close err %v\n", err)
	}()

	_, err := r.Read(make([]byte, 128))
	if err != nil {
		return err
	}

	return nil
}

func Main() error {
	flag.Parse()
	if len(flag.Args()) < 1 {
		return fmt.Errorf("please provide a serial file name")
	}
	fn := flag.Args()[0]

	f,err := sers.Open(fn)
	if err != nil {
		return err
	}

	return readpart(f)
}