package main

import (
	"flag"
	"log"
	"net"
)

func main() {

	iface4 := ""
	iface6 := ""
	flag.StringVar(&iface4, "interfaceIPv4", "", "Set here the IPv4 of the interface to bind to")
	flag.StringVar(&iface6, "interfaceIPv6", "", "Set here the IPv6 of the interface to bind to")
	flag.Parse()

	if errGeo != nil {
		log.Fatalf("Failed to read the geoIP database: %s\n", errGeo.Error())
	}

	go analytics.Report()

	if iface6 != "" {
		go func() {
			listener, err := net.Listen("tcp6", iface6+":443")
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
		}()
	}

	listener, err := net.Listen("tcp4", iface4+":443")
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
