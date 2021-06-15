package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	proxyurl, _ = url.Parse("socks5://127.0.0.1:9050")
	perhost     *proxy.PerHost
)

func initProxy() {
	socks5, _ := proxy.FromURL(proxyurl, proxy.Direct)
	perhost = proxy.NewPerHost(socks5, proxy.Direct)

	conf := "/etc/proxy/bypass.names"
	buf, err := os.Open(conf)
	if err != nil {
		log.Printf("Error opening file %s", conf)
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Printf("Error closing file %s : %s", conf, err.Error())
		}
	}()

	snl := bufio.NewScanner(buf)
	for snl.Scan() {

		if err := snl.Err(); err == nil {
			txt := snl.Text()
			if !strings.HasPrefix(txt, "#") && txt != "" {
				perhost.AddZone(txt)
			}
		} else {
			log.Printf("Error reading newline in file %s : %s", conf, err.Error())
			break
		}
	}
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

func SafeCopy(dst net.Conn, src io.Reader, wg *sync.WaitGroup, errc chan error) {
	defer wg.Done()
	select {
	case <-errc:
		return
	default:
		_, err := io.Copy(dst, src)
		if err == nil {
			dst.(*net.TCPConn).CloseWrite()
		} else {
			errc <- err
		}
	}
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

	var backendConn net.Conn

	if strings.HasSuffix(clientHello.ServerName, "tgragnato.it") {
		backendConn, err = net.DialTimeout("tcp", "127.0.0.1:8080", 10*time.Second)
	} else if checkDomain(clientHello.ServerName) {
		backendConn, err = perhost.Dial("tcp", net.JoinHostPort(clientHello.ServerName, "443"))
	} else {
		return
	}

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer backendConn.Close()

	go analytics.IncTLS(clientHello.ServerName)

	var wg sync.WaitGroup
	wg.Add(2)
	errc := make(chan error, 1)

	go SafeCopy(clientConn, backendConn, &wg, errc)
	go SafeCopy(backendConn, clientReader, &wg, errc)

	wg.Wait()
}
