package main

import (
	"net"
	"net/http"
)

func main() {
	daemon := make(chan struct{})

	go func() {

		handler := http.DefaultServeMux
		handler.HandleFunc("/", handleFunc)

		s := &http.Server{
			Addr:    "127.0.0.1:9080",
			Handler: handler,
		}

		s.ListenAndServe()

		daemon <- struct{}{}
	}()

	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:9081")
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
		daemon <- struct{}{}
	}()

	<-daemon
	<-daemon
}
