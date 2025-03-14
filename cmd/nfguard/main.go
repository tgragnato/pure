package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/miekg/dns"
	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/dnshandlers"
	"github.com/tgragnato/pure/pkg/http"
	"github.com/tgragnato/pure/pkg/sntp"
	"github.com/tgragnato/pure/pkg/spam"
)

func main() {

	var (
		iface4 string
		iface6 string
		dsn    string
		pyro   string
	)

	flag.StringVar(&iface4, "interfaceIPv4", "127.0.0.1", "Set here the IPv4 of the interface to bind to")
	flag.StringVar(&iface6, "interfaceIPv6", "[::1]", "Set here the IPv6 of the interface to bind to")
	flag.StringVar(&dsn, "dsn", "postgres://nfguard:nfguard@localhost:5432/nfguard?sslmode=disable", "Set here the DSN for the PostgreSQL database")
	flag.StringVar(&pyro, "pyroscope", "http://localhost:4040", "Set here the address of the Pyroscope server")
	flag.Parse()

	if pyro != "" {
		runtime.SetMutexProfileFraction(5)
		runtime.SetBlockProfileRate(5)

		profile, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: "nfguard.tgragnato.it",
			ServerAddress:   pyro,
			Logger:          nil,
			Tags:            map[string]string{"hostname": os.Getenv("HOSTNAME")},
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
		})
		if err == nil && profile != nil {
			defer func() {
				if err := profile.Stop(); err != nil {
					log.Println("failed to stop profiler:", err.Error())
				}
			}()
		}
	}

	geoChecks := checks.NewGeoChecks()

	handler, err := dnshandlers.MakeDnsHandlers(dsn, iface4, iface6, geoChecks)
	if err != nil {
		log.Fatalf("Failed to create DNS handlers: %s\n", err.Error())
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	err = sntp.Listen()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}

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
			Addr: "[::1]:53",
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
			Addr: "[::1]:53",
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

	go func() {
		udp := &dns.Server{
			Addr: iface4 + ":53",
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
			Addr: iface4 + ":53",
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

	go func() {
		udp := &dns.Server{
			Addr: iface6 + ":53",
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
			Addr: iface6 + ":53",
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

	http.Listen([]string{iface4, iface6}, dsn, geoChecks)

	httpWorker := make(chan spam.Spam, 1)
	stopped := false

	for !stopped {
		select {

		case <-signalCh:
			stopped = true

		case s := <-httpWorker:
			go s.Call()

		case httpWorker <- spam.MakeSpam():
			if spam.Counter < 10000 {
				time.Sleep(100 * time.Millisecond)
			} else if spam.Counter > 100000 {
				time.Sleep(time.Second)
			} else {
				time.Sleep(time.Duration(spam.Counter/100) * time.Millisecond)
			}
		}
	}
}
