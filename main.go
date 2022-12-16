package main

import (
	"flag"
	"log"
	"log/syslog"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

var (
	analytics               = &Analytics{data: map[string]Hits{}}
	express                 = &SafeExpress{data: map[string]uint{}}
	disableSyslog    bool   = false
	disableAppleOnly bool   = false
	interfaceIP      string = "172.16.33.0"
	interfaceIPv6    string = ""
	socks5           string = ""
	httpclient              = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Minute,
				KeepAlive: time.Millisecond,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       20,
			IdleConnTimeout:       5 * time.Minute,
			ResponseHeaderTimeout: 2 * time.Second,
		},
		Timeout: 5 * time.Minute,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	countryPath string
	dbreader    *maxminddb.Reader
)

func main() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&countryPath, "countryPath", "/var/db/GeoIP/GeoLite2-Country.mmdb.", "The path of the GeoIP2 County DataBase")
	flag.BoolVar(&disableSyslog, "disableSyslog", false, "Set this to disable the log redirection to syslog")
	flag.StringVar(&interfaceIP, "interfaceIP", "172.16.33.0", "Set here the IP of the interface to bind to")
	flag.BoolVar(&disableAppleOnly, "disableAppleOnly", false, "Set this to disable the pass filter inside unencrypted HTTP for Apple only")
	flag.StringVar(&socks5, "socks5", "", "Set this to the address of the upstream socks5 proxy you want to use")
	flag.StringVar(&interfaceIPv6, "interfaceIPv6", "", "Set here the IPv6 of the interface to bind to")
	flag.Parse()

	if !disableSyslog {
		syslogger, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "proxy")
		if err != nil {
			log.Fatalf("Failed to use syslog: %s\n", err.Error())
		}
		log.SetOutput(syslogger)
	}

	if socks5 != "" {
		proxy, err := url.Parse("socks5://" + socks5)
		if err != nil {
			return
		}
		httpclient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	}

	dbreader, _ = maxminddb.Open(countryPath)
	if dbreader != nil {
		defer dbreader.Close()
		log.Printf("Error opening %s\n", countryPath)
	}

	handler := http.DefaultServeMux
	handler.HandleFunc("/", handleHTTPForward)
	handler.HandleFunc(interfaceIP+"/", handleAnalytics)

	go func() {
		err := http.ListenAndServe(interfaceIP+":80", handler)
		if err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		}
	}()

	if interfaceIPv6 != "" {
		go func() {
			err := http.ListenAndServe(interfaceIPv6+":80", handler)
			if err != nil {
				log.Fatalf("Failed to start server: %s\n", err.Error())
			}
		}()

		go func() {
			listener, err := net.Listen("tcp6", interfaceIPv6+":443")
			if err != nil {
				log.Fatalf("Failed to start server: %s\n", err.Error())
			}
			for {
				flow, err := listener.Accept()
				if err != nil {
					continue
				}
				go EstablishFlow(flow)
			}
		}()
	}

	listener, err := net.Listen("tcp", interfaceIP+":443")
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
	for {
		flow, err := listener.Accept()
		if err != nil {
			continue
		}
		go EstablishFlow(flow)
	}
}
