package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tgragnato/pure/pkg/checks"
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

	for _, host := range hosts {
		go func(host string, writer *logWriter) {
			httpMux := http.NewServeMux()
			httpMux.Handle("/", loggingMiddleware(headersMiddleware(compressMiddleware(handleSHSHProtocol())), writer))
			srv := &http.Server{
				Addr:              host + ":80",
				Handler:           httpMux,
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      time.Minute,
				IdleTimeout:       time.Minute,
			}
			if err := srv.ListenAndServe(); err != nil {
				fmt.Println(err.Error())
			}
		}(host, writer)
	}
}
