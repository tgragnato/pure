package sni

import (
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

	backendConn, err := net.DialTimeout("tcp", getHostPort(clientHello.ServerName), time.Second)
	if err != nil {
		return
	}
	defer backendConn.Close()

	copyLoop(clientReader, clientConn, backendConn)
}
