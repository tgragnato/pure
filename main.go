package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

func main() {
	rand.Seed(time.Now().Unix())
	initCheck()
	initProxy()

	go func() {
		muxAnalytics := http.NewServeMux()
		muxAnalytics.HandleFunc("/", handleAnalytics)
		err := http.ListenAndServe("127.0.0.1:1080", muxAnalytics)
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		dns.HandleFunc(".", handleDnsRequest)
		udp := &dns.Server{Addr: ":53", Net: "udp"}
		err := udp.ListenAndServe()
		defer udp.Shutdown()
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		dns.HandleFunc(".", handleDnsRequest)
		tcp := &dns.Server{Addr: ":53", Net: "tcp"}
		err := tcp.ListenAndServe()
		defer tcp.Shutdown()
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 123})
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
		for {
			request := make([]byte, 512)
			rlen, remote, err := listener.ReadFromUDP(request[0:])
			if err != nil {
				continue
			}
			if rlen > 0 && validFormat(request) {
				go listener.WriteTo(generate(request), remote)
			}
		}
	}()

	go func() {
		handler := http.DefaultServeMux
		handler.HandleFunc("/", handleHTTPForward)
		err := http.ListenAndServe("172.16.31.0:80", handler)
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		IPv6handler := http.NewServeMux()
		IPv6handler.HandleFunc("/", handleHTTPForward)
		err := http.ListenAndServe("[fd76:abcd:ef90::]:80", IPv6handler)
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", "[fd76:abcd:ef90::]:443")
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
		}
		for {
			flow, err := listener.Accept()
			if err != nil {
				continue
			}
			go establishFlow(flow)
		}
	}()

	listener, err := net.Listen("tcp", "172.16.31.0:443")
	if err != nil {
		log.Printf("Failed to start server: %s\n ", err.Error())
	}
	for {
		flow, err := listener.Accept()
		if err != nil {
			continue
		}
		go establishFlow(flow)
	}
}
