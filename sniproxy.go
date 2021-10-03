package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/net/proxy"
)

var (
	proxyurl, _ = url.Parse("socks5://127.0.0.1:9050")
	socks5, _   = proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, &net.Dialer{Timeout: time.Second})
	perhost     *proxy.PerHost
)

const SO_ORIGINAL_DST = 80

func initProxy() {
	cfsocks, _ := proxy.SOCKS5("tcp", "127.0.0.1:9040", nil, &net.Dialer{Timeout: 10 * time.Second})
	cfproxy := proxy.NewPerHost(cfsocks, proxy.Direct)
	SetBypass("/etc/proxy/fallback.names", cfproxy)
	perhost = proxy.NewPerHost(socks5, cfproxy)
	SetBypass("/etc/proxy/bypass.names", perhost)
}

func SetBypass(conf string, newproxy *proxy.PerHost) {
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
				newproxy.AddZone(txt)
			}
		} else {
			log.Printf("Error reading newline in file %s : %s", conf, err.Error())
			break
		}
	}
}

func ReadLine(conf string) string {
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
	snl.Scan()
	err = snl.Err()
	if err == nil {
		return snl.Text()
	}
	return ""
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

func IPAddress(ip []byte) net.IP {
	switch len(ip) {

	case 4:
		return net.IPv4(ip[0], ip[1], ip[2], ip[3])

	case 16:
		if bytes.Equal(ip[:12], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}) {
			return IPAddress(ip[12:16])
		}
		return net.IP([]byte{
			ip[0], ip[1], ip[2], ip[3],
			ip[4], ip[5], ip[6], ip[7],
			ip[8], ip[9], ip[10], ip[11],
			ip[12], ip[13], ip[14], ip[15],
		})

	default:
		return nil
	}
}

func OriginalAddress(conn *net.TCPConn) (destIP net.IP, destPort int) {
	f, err := conn.File()
	if err != nil {
		return
	}
	defer f.Close()
	fd := int(f.Fd())

	level := syscall.IPPROTO_IP
	if conn.RemoteAddr().String()[0] == '[' {
		level = syscall.IPPROTO_IPV6
	}
	addr, err := syscall.GetsockoptIPv6MTUInfo(fd, level, SO_ORIGINAL_DST)
	if err != nil {
		return
	}

	ip := (*[4]byte)(unsafe.Pointer(&addr.Addr.Flowinfo))[:4]
	if level == syscall.IPPROTO_IPV6 {
		ip = addr.Addr.Addr[:]
	}
	port := (*[2]byte)(unsafe.Pointer(&addr.Addr.Port))[:2]

	destIP = IPAddress(ip)
	destPort = int(port[0])*256 + int(port[1])
	return
}
