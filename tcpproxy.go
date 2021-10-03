package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	_, private4, _ = net.ParseCIDR("172.16.31.0/24")
	_, private6, _ = net.ParseCIDR("fd76:abcd:ef90::/120")
	cfonion        = [10]string{
		"cflarexljc3rw355ysrkrzwapozws6nre6xsy3n4yrj7taye3uiby3ad.onion",
		"cflarenuttlfuyn7imozr4atzvfbiw3ezgbdjdldmdx7srterayaozid.onion",
		"cflares35lvdlczhy3r6qbza5jjxbcplzvdveabhf7bsp7y4nzmn67yd.onion",
		"cflareusni3s7vwhq2f7gc4opsik7aa4t2ajedhzr42ez6uajaywh3qd.onion",
		"cflareki4v3lh674hq55k3n7xd4ibkwx3pnw67rr3gkpsonjmxbktxyd.onion",
		"cflarejlah424meosswvaeqzb54rtdetr4xva6mq2bm2hfcx5isaglid.onion",
		"cflaresuje2rb7w2u3w43pn4luxdi6o7oatv6r2zrfb5xvsugj35d2qd.onion",
		"cflareer7qekzp3zeyqvcfktxfrmncse4ilc7trbf6bp6yzdabxuload.onion",
		"cflareub6dtu7nvs3kqmoigcjdwap2azrkx5zohb2yk7gqjkwoyotwqd.onion",
		"cflare2nge4h4yqr3574crrd7k66lil3torzbisz6uciyuzqc2h2ykyd.onion",
	}
	onionbypass = populateCheck("/etc/proxy/bypass.names")
	names453    = ReadLine("/etc/proxy/453.names")
)

func EstablishFlow(clientConn net.Conn) {
	defer clientConn.Close()
	ip := clientConn.RemoteAddr().(*net.TCPAddr).IP
	if !private4.Contains(ip) && !private6.Contains(ip) {
		HandleReverse(clientConn, ip)
	} else {
		HandleForward(clientConn)
	}
}

func HandleReverse(clientConn net.Conn, remote net.IP) {
	clientHello, clientReader, err := HandleClientHello(clientConn)
	if err != nil {
		return
	}
	var backendConn net.Conn
	if strings.HasSuffix(clientHello.ServerName, "tgragnato.it") && IsAllowed(remote) {
		backendConn, err = net.DialTimeout("tcp", "127.0.0.1:8080", 10*time.Second)
	} else {
		backendConn, err = net.DialTimeout("tcp", "127.0.0.1:9001", 10*time.Second)
	}
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer backendConn.Close()
	CopyLoop(clientReader, clientConn, backendConn)
}

func HandleForward(clientConn net.Conn) {
	destIP, destPort := OriginalAddress(clientConn.(*net.TCPConn))

	switch destPort {
	case 443:
		clientHello, clientReader, err := HandleClientHello(clientConn)
		if err != nil {
			return
		}
		var backendConn net.Conn
		if strings.HasSuffix(clientHello.ServerName, "tgragnato.it") {
			backendConn, err = net.DialTimeout("tcp", "127.0.0.1:8080", 10*time.Second)
		} else if checkDomain(Cloaking(clientHello.ServerName)) {
			backendConn, err = CustomDialer(clientHello.ServerName)
		} else {
			return
		}
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()
		go analytics.IncTLS(clientHello.ServerName)
		CopyLoop(clientReader, clientConn, backendConn)

	case 445:
		if !net.ParseIP("172.16.31.0").Equal(destIP) {
			return
		}
		backendConn, err := net.DialTimeout("tcp", "127.0.0.1:5445", 10*time.Second)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()
		CopyLoop(clientConn, clientConn, backendConn)

	case 453:
		clientHello, clientReader, err := HandleClientHello(clientConn)
		if err != nil {
			return
		}
		var backendConn net.Conn
		if clientHello.ServerName == names453 {
			backendConn, err = socks5.Dial("tcp", net.JoinHostPort(names453, "453"))
		} else {
			return
		}
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()
		go analytics.IncTLS(clientHello.ServerName)
		CopyLoop(clientReader, clientConn, backendConn)

	case 993:
		clientHello, clientReader, err := HandleClientHello(clientConn)
		if err != nil {
			return
		}
		var backendConn net.Conn
		if strings.HasSuffix(clientHello.ServerName, "imap.mail.me.com") {
			backendConn, err = socks5.Dial("tcp", net.JoinHostPort(clientHello.ServerName, "993"))
		} else {
			return
		}
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()
		go analytics.IncTLS(clientHello.ServerName)
		CopyLoop(clientReader, clientConn, backendConn)

	case 5222, 5223:
		if !checkIPs([]net.IP{destIP}) {
			return
		}
		if destPort == 5222 && IP2ASN(destIP) != 32934 {
			return
		}
		if destPort == 5223 && IP2ASN(destIP) != 714 {
			return
		}
		backendConn, err := socks5.Dial("tcp", net.JoinHostPort(destIP.String(), strconv.Itoa(destPort)))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		defer backendConn.Close()
		CopyLoop(clientConn, clientConn, backendConn)

	default:
		return
	}
}

func CustomDialer(address string) (backendConn net.Conn, err error) {
	if address == "duckduckgo.com" || strings.HasSuffix(address, ".duckduckgo.com") {
		backendConn, err = socks5.Dial("tcp", "duckduckgogg42xjoc72x3sjasowoarfbgcmvfimaftt6twagswzczad.onion:443")
		if err == nil {
			return
		}
	}
	if CheckCF(address) {
		is_bypass := false
		for _, zone := range onionbypass {
			if strings.HasSuffix(address+".", zone) {
				is_bypass = true
				break
			}
		}
		if !is_bypass {
			bsl := rand.Float64() * float64(len(cfonion))
			for i := 0; i < len(cfonion); i++ {
				index := (int(bsl) + i) % len(cfonion)
				backendConn, err = socks5.Dial("tcp", net.JoinHostPort(cfonion[index], "443"))
				if err == nil {
					return
				}
			}
		}
	}
	return perhost.Dial("tcp", net.JoinHostPort(address, "443"))
}
