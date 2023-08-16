package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 123})
	if err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
		return
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
}
