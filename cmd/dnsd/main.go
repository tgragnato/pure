package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
	"github.com/tgragnato/pure/pkg/dnshandlers"
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

	handler, err := dnshandlers.MakeDnsHandlers(dsn, hint4, hint6)
	if err != nil {
		log.Fatalf("Failed to create DNS handlers: %s\n", err.Error())
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			handler.ParseQuery(m)
		}

		err := w.WriteMsg(m)
		if err != nil {
			log.Printf("Failed to write message: %s\n", err.Error())
		}
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
			defer func() {
				err := udp.Shutdown()
				if err != nil {
					log.Fatalf("Failed to shutdown server: %s\n", err.Error())
				}
			}()
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
			defer func() {
				err := tcp.Shutdown()
				if err != nil {
					log.Fatalf("Failed to shutdown server: %s\n", err.Error())
				}
			}()
		}
	}()

	<-signalCh
}
