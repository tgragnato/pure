package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"strings"
	"time"
)

func peekClientHello(reader io.Reader) (*tls.ClientHelloInfo, io.Reader, error) {
	peekedBytes := new(bytes.Buffer)
	hello, err := readClientHello(io.TeeReader(reader, peekedBytes))
	if err != nil {
		return nil, nil, err
	}
	return hello, io.MultiReader(peekedBytes, reader), nil
}

type readOnlyConn struct {
	reader io.Reader
}

func (conn readOnlyConn) Read(p []byte) (int, error)         { return conn.reader.Read(p) }
func (conn readOnlyConn) Write(p []byte) (int, error)        { return 0, io.ErrClosedPipe }
func (conn readOnlyConn) Close() error                       { return nil }
func (conn readOnlyConn) LocalAddr() net.Addr                { return nil }
func (conn readOnlyConn) RemoteAddr() net.Addr               { return nil }
func (conn readOnlyConn) SetDeadline(t time.Time) error      { return nil }
func (conn readOnlyConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn readOnlyConn) SetWriteDeadline(t time.Time) error { return nil }

func readClientHello(reader io.Reader) (*tls.ClientHelloInfo, error) {
	var hello *tls.ClientHelloInfo

	err := tls.Server(readOnlyConn{reader: reader}, &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello = new(tls.ClientHelloInfo)
			*hello = *argHello
			return nil, nil
		},
	}).Handshake()

	if hello == nil {
		return nil, err
	}

	return hello, nil
}

func establishFlow(clientConn net.Conn) {
	defer clientConn.Close()

	err := clientConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}

	clientHello, clientReader, err := peekClientHello(clientConn)
	if err != nil {
		return
	}

	err = clientConn.SetReadDeadline(time.Time{})
	if err != nil {
		return
	}

	saddr := clientConn.RemoteAddr().String()
	host := backend(strings.HasPrefix(saddr, "172.16.31."), clientHello.ServerName)
	if host == "" {
		return
	}
	backendConn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return
	}
	defer backendConn.Close()

	done := make(chan struct{})

	go func() {
		io.Copy(clientConn, backendConn)
		clientConn.(*net.TCPConn).CloseWrite()
		done <- struct{}{}
	}()
	go func() {
		io.Copy(backendConn, clientReader)
		backendConn.(*net.TCPConn).CloseWrite()
		done <- struct{}{}
	}()

	<-done
	<-done
}

func backend(local bool, sni string) string {
	if local {
		if strings.HasSuffix(sni, "tgragnato.it") && sni != "status.tgragnato.it" {
			return "127.0.0.1:8080"
		} else {
			if checkDomain(sni) {
				return net.JoinHostPort(sni, "443")
			} else {
				return ""
			}
		}
	} else {
		if strings.HasSuffix(sni, "tgragnato.it") {
			return "127.0.0.1:8080"
		} else if strings.HasSuffix(sni, "awsmppl.com") ||
			strings.HasSuffix(sni, "dnsupdate.info") ||
			strings.HasSuffix(sni, "nerdpol.ovh") ||
			strings.HasSuffix(sni, "nsupdate.info") ||
			strings.HasSuffix(sni, "urown.cloud") {
			return "127.0.0.1:8081"
		} else {
			return "127.0.0.1:9001"
		}
	}
}
