package dnshandlers

import (
	"database/sql"
	"errors"
	"net"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tgragnato/pure/pkg/checks"
)

type DnsHandlers struct {
	db        *sql.DB
	geoChecks *checks.GeoChecks
	hintIPv4  net.IP
	hintIPv6  net.IP
	muIPv4    bool
	muIPv6    bool
}

func MakeDnsHandlers(dsn string, hint4 string, hint6 string, geoChecks *checks.GeoChecks) (*DnsHandlers, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil || (db != nil && db.Ping() != nil) {
		return nil, err
	}

	hintIPv4 := net.ParseIP(hint4).To4()
	if hintIPv4 == nil {
		return nil, errors.New("failed to parse IPv4 hint")
	}

	hint6 = strings.Trim(hint6, "[")
	hint6 = strings.Trim(hint6, "]")
	hintIPv6 := net.ParseIP(hint6)
	if hintIPv6 == nil {
		return nil, errors.New("failed to parse IPv6 hint")
	}

	d := &DnsHandlers{
		db:        db,
		geoChecks: geoChecks,
		hintIPv4:  hintIPv4,
		hintIPv6:  hintIPv6,
	}

	d.cleanPersistent()
	go d.crossPrefetch()
	go func() {
		for range time.NewTicker(time.Minute).C {
			go d.selfPrefetch(false)
			go d.selfPrefetch(true)
		}
	}()

	return d, nil
}
