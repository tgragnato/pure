package sni

import (
	"log"
	"net"
	"time"
)

func establishFlow(clientConn net.Conn) {
	defer clientConn.Close()

	clientHello, clientReader, err := handleClientHello(clientConn)
	if err != nil {
		return
	}

	if !checkDomain(clientHello.ServerName) {
		return
	}

	backendConn, err := net.DialTimeout("tcp", net.JoinHostPort(clientHello.ServerName, "443"), time.Second)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer backendConn.Close()

	copyLoop(clientReader, clientConn, backendConn)
}
