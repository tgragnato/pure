package sni

import (
	"fmt"
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
		Deadline:      time.Now().Add(time.Minute),
		DualStack:     true,
		FallbackDelay: time.Second,
		KeepAlive:     100 * time.Millisecond,
	}
	d.SetMultipathTCP(true)

	backendConn, err := d.Dial("tcp", getHostPort(clientHello.ServerName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer backendConn.Close()

	copyLoop(clientReader, clientConn, backendConn)
}
