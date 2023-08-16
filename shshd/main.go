package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	iface4 := ""
	iface6 := ""
	flag.StringVar(&iface4, "interfaceIPv4", "", "Set here the IPv4 of the interface to bind to")
	flag.StringVar(&iface6, "interfaceIPv6", "", "Set here the IPv6 of the interface to bind to")
	flag.Parse()

	handler := http.DefaultServeMux
	handler.HandleFunc("/", handleHTTPForward)

	if iface6 != "" {
		go func() {
			err := http.ListenAndServe(iface6+":80", handler)
			if err != nil {
				log.Fatalf("Failed to start server: %s\n", err.Error())
			}
		}()
	}

	err := http.ListenAndServe(iface4+":80", handler)
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}
