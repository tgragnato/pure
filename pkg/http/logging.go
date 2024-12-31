package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/tgragnato/pure/pkg/checks"
)

type log struct {
	Epoch   int64  `json:"epoch"`
	Remote  string `json:"remote"`
	Country string `json:"country"`
	Proto   string `json:"proto"`
	Host    string `json:"host"`
	Method  string `json:"method"`
	Request string `json:"request"`
	Status  int    `json:"status"`
	Bytes   int64  `json:"bytes"`
}

type logWriter struct {
	geoChecks *checks.GeoChecks
	db        *sql.DB
}

func (w *logWriter) writeLog(l log) {
	if w.db == nil {
		return
	}

	if _, err := w.db.Exec(
		`
			INSERT INTO http (epoch, remote, country, proto, host, method, request, status, bytes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		l.Epoch, l.Remote, l.Country, l.Proto, l.Host, l.Method, l.Request, l.Status, l.Bytes,
	); err != nil {
		logJson, errJson := json.Marshal(l)
		if errJson != nil {
			fmt.Println("Error marshalling log:", errJson.Error())
		} else {
			fmt.Printf("Error (%s) writing log: %s\n", err.Error(), string(logJson))
		}
	}
}

func loggingMiddleware(next http.Handler, writer *logWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		epoch := time.Now().Unix()

		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		remote, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remote = ""
		}

		country := ""
		if writer.geoChecks != nil && remote != "" {
			country = writer.geoChecks.GetCountry(net.ParseIP(remote))
		}

		host := r.Host
		if host == "" {
			host = r.URL.Host
		}

		path := r.RequestURI
		if path == "" {
			path = r.URL.RawPath
		}

		go writer.writeLog(log{
			Epoch:   epoch,
			Remote:  remote,
			Country: country,
			Proto:   r.Proto,
			Host:    host,
			Method:  r.Method,
			Request: path,
			Status:  wrapped.status,
			Bytes:   wrapped.size,
		})
	})
}
