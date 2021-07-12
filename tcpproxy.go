package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func EstablishFlow(clientConn net.Conn) {
	defer clientConn.Close()
	destIP, destPort := OriginalAddress(clientConn.(*net.TCPConn))

	switch destPort {
	case 443, 993:

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
		if destPort == 993 && !strings.HasSuffix(clientHello.ServerName, "imap.mail.me.com") {
			return
		}

		var backendConn net.Conn
		if destPort == 443 && strings.HasSuffix(clientHello.ServerName, "tgragnato.it") {
			backendConn, err = net.DialTimeout("tcp", "127.0.0.1:8080", 10*time.Second)
		} else if checkDomain(clientHello.ServerName) {
			backendConn, err = perhost.Dial("tcp", net.JoinHostPort(clientHello.ServerName, strconv.Itoa(destPort)))
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
		ctx, cancel := context.WithCancel(context.Background())
		go SafeCopy(clientConn, backendConn, &wg, ctx, cancel)
		go SafeCopy(backendConn, clientReader, &wg, ctx, cancel)
		wg.Wait()

	case 5900:

		backendConn, err := net.DialTimeout("tcp", "172.16.31.2:5900", 10*time.Second)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()

		var wg sync.WaitGroup
		wg.Add(2)
		ctx, cancel := context.WithCancel(context.Background())
		go SafeCopy(clientConn, backendConn, &wg, ctx, cancel)
		go SafeCopy(backendConn, clientConn, &wg, ctx, cancel)
		wg.Wait()

	case 5222, 5223:

		if !checkIPs([]net.IP{destIP}) {
			return
		}

		backendConn, err := socks5.Dial("tcp", net.JoinHostPort(destIP.String(), strconv.Itoa(destPort)))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()

		var wg sync.WaitGroup
		wg.Add(2)
		ctx, cancel := context.WithCancel(context.Background())
		go SafeCopy(clientConn, backendConn, &wg, ctx, cancel)
		go SafeCopy(backendConn, clientConn, &wg, ctx, cancel)
		wg.Wait()

	default:
		return
	}

}
