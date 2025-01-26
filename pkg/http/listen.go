package http

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tgragnato/pure/pkg/checks"
	"golang.org/x/crypto/acme/autocert"
)

func Listen(hosts []string, dsn string, geoChecks *checks.GeoChecks) {
	db, err := sql.Open("pgx", dsn)
	if err != nil || (db != nil && db.Ping() != nil) {
		db = nil
	}
	writer := &logWriter{
		geoChecks: geoChecks,
		db:        db,
	}

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache("/var/www/acme"),
	}

	go func(manager *autocert.Manager, writer *logWriter) {
		httpMux := http.NewServeMux()
		httpMux.Handle("/", loggingMiddleware(headersMiddleware(handleSHSHProtocol()), writer))
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
	}(manager, writer)

	go func(manager *autocert.Manager, writer *logWriter) {
		httpsMux := http.NewServeMux()
		httpsMux.Handle("/", loggingMiddleware(headersMiddleware(apiGateway()), writer))
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
	}(manager, writer)

}
