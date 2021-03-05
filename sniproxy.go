package main

import (
        "bytes"
        "crypto/tls"
        "io"
        "net"
        "os"
        "time"
)

func main() {
        listener, err := net.Listen("tcp", "127.0.0.1:9081")
        if err != nil {
                os.Exit(1)
        }
        for {
                flow, err := listener.Accept()
                if err != nil {
                        continue
                }
                go establishFlow(flow)
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

func establishFlow(clientConn net.Conn) {
        defer clientConn.Close()

        err := clientConn.SetReadDeadline(time.Now().Add(10 * time.Second));
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

        ip, err := net.LookupIP(clientHello.ServerName)
        if err != nil {
                return
        }
        for i := range ip {
                if !ip[i].IsGlobalUnicast() {
                        return
                }
        }

        backendConn, err := net.DialTimeout("tcp", net.JoinHostPort(clientHello.ServerName, "443"), 10*time.Second)
        if err != nil {
                return
        }
        defer backendConn.Close()

        done := make(chan struct{})

        go func() {
                io.Copy(clientConn, backendConn)
                clientConn.(*net.TCPConn).CloseWrite()
                done<-struct{}{}
        }()
        go func() {
                io.Copy(backendConn, clientReader)
                backendConn.(*net.TCPConn).CloseWrite()
                done<-struct{}{}
        }()

        <-done
        <-done
}
