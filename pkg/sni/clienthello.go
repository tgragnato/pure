package sni

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"time"
)

func handleClientHello(clientConn net.Conn) (clientHello *tls.ClientHelloInfo, clientReader io.Reader, err error) {
	err = clientConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}

	clientHello, clientReader, err = peekClientHello(clientConn)
	if err != nil {
		return
	}

	err = clientConn.SetReadDeadline(time.Time{})
	if err != nil {
		return
	}

	return
}

func peekClientHello(reader io.Reader) (*tls.ClientHelloInfo, io.Reader, error) {
	peekedBytes := new(bytes.Buffer)
	hello, err := readClientHello(io.TeeReader(reader, peekedBytes))
	if err != nil {
		return nil, nil, err
	}
	return hello, io.MultiReader(peekedBytes, reader), nil
}
