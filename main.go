package main

import (
	"net"
	"net/http"
)

func main() {
	initCheck()
	daemon := make(chan struct{})

	go func() {
		handler := http.DefaultServeMux
		handler.HandleFunc("/", handleHTTP)
		http.ListenAndServe(":80", handler)
		daemon <- struct{}{}
	}()

	go func() {
		listener, err := net.Listen("tcp", ":443")
		if err != nil {
			daemon <- struct{}{}
		}
		for {
			flow, err := listener.Accept()
			if err != nil {
				continue
			}
			go establishFlow(flow)
		}
		//daemon <- struct{}{}
	}()

	<-daemon
	<-daemon
}
