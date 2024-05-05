package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miekg/dns"
	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/dnshandlers"
	"github.com/tgragnato/pure/pkg/errcache"
	"github.com/tgragnato/pure/pkg/ipcache"
)

func main() {

	var (
		hint4    string
		hint6    string
		bindAddr string
		dsn      string
	)

	flag.StringVar(&hint4, "hintIPv4", "", "Set here the IPv4 of the HTTPS hint")
	flag.StringVar(&hint6, "hintIPv6", "", "Set here the IPv6 of the HTTPS hint")
	flag.StringVar(&bindAddr, "bindAddr", "127.0.0.1:53", "Set here the address to bind to")
	flag.StringVar(&dsn, "dsn", "postgres://dnsd:dnsd@localhost:5432/dnsd?sslmode=disable", "Set here the DSN for the PostgreSQL database")
	flag.Parse()

	var (
		cache4    = ipcache.NewCache(3600*time.Second, false, dsn)
		cache6    = ipcache.NewCache(3600*time.Second, true, dsn)
		geoChecks = checks.NewGeoChecks()
		errCache4 = errcache.NewErrCache(time.Minute, false, cache4, cache6, geoChecks)
		errCache6 = errcache.NewErrCache(time.Minute, true, cache4, cache6, geoChecks)
	)

	hintIPv4 := net.ParseIP(hint4).To4()
	if hintIPv4 == nil {
		log.Fatalln("Failed to parse IPv4 hint")
	}

	hintIPv6 := net.ParseIP(hint6)
	if hintIPv6 == nil {
		log.Fatalln("Failed to parse IPv6 hint")
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			dnshandlers.ParseQuery(
				m,
				cache4,
				cache6,
				errCache4,
				errCache6,
				geoChecks,
				hintIPv4,
				hintIPv6,
			)
		}

		w.WriteMsg(m)
	})

	go func() {
		udp := &dns.Server{
			Addr: bindAddr,
			Net:  "udp",
		}
		err := udp.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		} else {
			defer udp.Shutdown()
		}
	}()

	go func() {
		tcp := &dns.Server{
			Addr: bindAddr,
			Net:  "tcp",
		}
		err := tcp.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		} else {
			defer tcp.Shutdown()
		}
	}()

	<-signalCh
}
