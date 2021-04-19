package main

import (
	"log"
	"net"
	"net/http"

	"github.com/miekg/dns"
)

func main() {
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
		server := &dns.Server{Addr: "127.0.0.1:1053", Net: "udp"}
		err := server.ListenAndServe()
		defer server.Shutdown()
		if err != nil {
			log.Printf("Failed to start server: %s\n ", err.Error())
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
