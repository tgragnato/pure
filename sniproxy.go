package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func EstablishFlow(clientConn net.Conn) {
	defer clientConn.Close()

	clientHello, clientReader, err := HandleClientHello(clientConn)
	if err != nil {
		return
	}

	if !CheckDomain(clientHello.ServerName) {
		return
	}
	go analytics.IncTLS(clientHello.ServerName)

	var backendConn net.Conn
	if socks5 != "" {
		dialer, err := proxy.SOCKS5("tcp", socks5, nil, proxy.Direct)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		backendConn, err = dialer.Dial("tcp", net.JoinHostPort(clientHello.ServerName, "443"))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
	} else {
		backendConn, err = net.Dial("tcp", net.JoinHostPort(clientHello.ServerName, "443"))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
	}
	defer backendConn.Close()

	CopyLoop(clientReader, clientConn, backendConn)
}

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

func SafeCopy(dst net.Conn, src io.Reader, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	r := &readerCtx{ctx: ctx, r: src}
	_, err := io.Copy(dst, r)
	dst.(*net.TCPConn).CloseWrite()
	if err != nil {
		cancel()
	}
}

func CopyLoop(clientR io.Reader, clientW net.Conn, backend net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go SafeCopy(clientW, backend, &wg, ctx, cancel)
	go SafeCopy(backend, clientR, &wg, ctx, cancel)
	wg.Wait()
}

func HandleClientHello(clientConn net.Conn) (clientHello *tls.ClientHelloInfo, clientReader io.Reader, err error) {
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
