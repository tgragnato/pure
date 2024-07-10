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

	d := &net.Dialer{
		Timeout:       time.Second,
		Deadline:      time.Now().Add(time.Second),
		DualStack:     true,
		FallbackDelay: time.Second / 2,
		KeepAlive:     100 * time.Millisecond,
	}
	d.SetMultipathTCP(true)

	backendConn, err := d.Dial("tcp", getHostPort(clientHello.ServerName))
	if err != nil {
		return
	}
	defer backendConn.Close()

	copyLoop(clientReader, clientConn, backendConn)
}
