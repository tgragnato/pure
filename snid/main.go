package main

import (
	"log"
	"net"
)

func main() {

	if errGeo != nil {
		log.Fatalf("Failed to read the geoIP database: %s\n", errGeo.Error())
	}

	go analytics.Report()

	listener, err := net.Listen("tcp6", "[::1]:443")
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
	for {
		flow, err := listener.Accept()
		if err != nil {
			continue
		}
		go EstablishFlow(flow)
	}

}
