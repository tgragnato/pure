package main

import (
	"flag"
	"log"
	"log/syslog"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

var (
	analytics               = &Analytics{data: map[string]Hits{}}
	express                 = &SafeExpress{data: map[string]uint{}}
	disableSyslog    bool   = false
	disableAppleOnly bool   = false
	interfaceIP      string = "172.16.33.0"
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
	flag.Parse()

	if !disableSyslog {
		syslogger, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "proxy")
		if err != nil {
			log.Fatalf("Failed to use syslog: %s\n", err.Error())
		}
		log.SetOutput(syslogger)
	}

	dbreader, _ = maxminddb.Open(countryPath)
	if dbreader != nil {
		defer dbreader.Close()
		log.Printf("Error opening %s\n", countryPath)
	}

	go func() {
		handler := http.DefaultServeMux
		handler.HandleFunc("/", handleHTTPForward)
		handler.HandleFunc(interfaceIP+"/", handleAnalytics)
		err := http.ListenAndServe(interfaceIP+":80", handler)
		if err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		}
	}()

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
