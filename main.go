package main

import (
	"net"
	"net/http"
	"os"
)

func main() {
	initCheck()

	go func() {
		handler := http.DefaultServeMux
		handler.HandleFunc("/", handleHTTP)
		http.ListenAndServe(":80", handler)
	}()

	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		os.Exit(1)
	}
	for {
		flow, err := listener.Accept()
		if err != nil {
			continue
		}
		go establishFlow(flow)
	}
}
