package main

import (
	"log"

	"github.com/miekg/dns"
)

func main() {

	dns.HandleFunc(".", handleDnsRequest)
	initCheck()

	go func() {
		udp := &dns.Server{
			Addr: "127.0.0.1:53",
			Net:  "udp",
		}
		err := udp.ListenAndServe()
		if err != nil {
			log.Printf("Failed to start server: %s\n", err.Error())
		} else {
			defer udp.Shutdown()
		}
	}()

	tcp := &dns.Server{
		Addr: "127.0.0.1:53",
		Net:  "tcp",
	}
	err := tcp.ListenAndServe()
	if err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
	} else {
		defer tcp.Shutdown()
	}

}
