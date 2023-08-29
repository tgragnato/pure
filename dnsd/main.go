package main

import (
	"flag"
	"log"
	"net"

	"github.com/miekg/dns"
)

func main() {

	hint4 := ""
	hint6 := ""
	flag.StringVar(&hint4, "hintIPv4", "", "Set here the IPv4 of the HTTPS hint")
	flag.StringVar(&hint6, "hintIPv6", "", "Set here the IPv6 of the HTTPS hint")
	flag.Parse()

	if hint4 == "" || hint6 == "" {
		log.Fatalln("Failed to acquire hints")
	}

	hintIPv4 = net.ParseIP(hint4).To4()
	if hintIPv4 == nil {
		log.Fatalln("Failed to parse IPv4")
	}

	hintIPv6 = net.ParseIP(hint6)
	if hintIPv6 == nil {
		log.Fatalln("Failed to parse IPv6")
	}

	dns.HandleFunc(".", handleDnsRequest)
	initCheck()

	go func() {
		udp := &dns.Server{
			Addr: "127.0.0.1:53",
			Net:  "udp",
		}
		err := udp.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
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
		log.Fatalf("Failed to start server: %s\n", err.Error())
	} else {
		defer tcp.Shutdown()
	}

}
