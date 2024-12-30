package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func Listen(hosts []string) {

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache("/var/www/acme"),
	}

	go func(manager *autocert.Manager) {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", handleSHSHProtocol)
		srv := &http.Server{
			Addr:              ":80",
			Handler:           manager.HTTPHandler(httpMux),
			ReadTimeout:       time.Second,
			ReadHeaderTimeout: time.Second,
			WriteTimeout:      time.Minute,
			IdleTimeout:       time.Minute,
		}
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}(manager)

	go func(manager *autocert.Manager) {
		httpsMux := http.NewServeMux()
		httpsMux.Handle("/", headersMiddleware(compressMiddleware(apiGateway())))
		srv := &http.Server{
			Addr:              ":443",
			Handler:           httpsMux,
			ReadTimeout:       time.Second,
			ReadHeaderTimeout: time.Second,
			WriteTimeout:      time.Minute,
			IdleTimeout:       time.Minute,
			TLSConfig: &tls.Config{
				GetCertificate:         manager.GetCertificate,
				NextProtos:             []string{"h2", "http/1.1"},
				SessionTicketsDisabled: true,
				MinVersion:             tls.VersionTLS13,
				CurvePreferences:       []tls.CurveID{tls.X25519, tls.CurveP521},
			},
		}
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			fmt.Println(err.Error())
		}
	}(manager)

}
