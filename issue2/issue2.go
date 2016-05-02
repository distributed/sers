package main

import (
	"log"
	"github.com/distributed/sers"
	"time"
	"flag"
)


var port = flag.String("port", "/dev/ttyAMA0", "")
var baudrate = flag.Int("baudrate", 9600, "")

func main() {
	flag.Parse()

	log.Printf("using %s with %d baud\n", *port, *baudrate)

	head, err := sers.Open(*port)
	if err != nil {
		log.Fatal(err)
	}

	err = head.SetMode(*baudrate, 8, sers.N, 1, sers.NO_HANDSHAKE)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err = head.Write([]byte{0x55})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 10)
	}
	head.Close()
}
