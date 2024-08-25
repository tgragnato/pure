package dnshandlers

import (
	"fmt"
	"net"
	"strings"

	"github.com/tgragnato/pure/pkg/checks"
)

func (d *DnsHandlers) setPersistent(key string, data []net.IP, ipv6 bool) {
	serialize := ""
	for _, ip := range data {
		serialize += ip.String() + ","
	}

	query := `
		INSERT INTO a (key, value, discovered_on, last_used)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, last_used = CURRENT_TIMESTAMP
		;
	`
	if ipv6 {
		query = `
		INSERT INTO aaaa (key, value, discovered_on, last_used)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, last_used = CURRENT_TIMESTAMP
		;
	`
	}

	_, err := d.db.Exec(query, key, serialize)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (d *DnsHandlers) getPersistent(key string, ipv6 bool) (data []net.IP, found bool) {
	query := `
		SELECT value
		FROM a
		WHERE key = $1
	`
	if ipv6 {
		query = `
		SELECT value
		FROM aaaa
		WHERE key = $1
	`
	}

	row := d.db.QueryRow(query, key)
	var serialize string
	err := row.Scan(&serialize)
	if err != nil {
		return []net.IP{}, false
	}

	ips := []net.IP{}
	for _, ip := range strings.Split(serialize, ",") {
		if parsed := net.ParseIP(ip); parsed != nil {
			ips = append(ips, parsed)
		}
	}

	return ips, true
}

func (d *DnsHandlers) cleanPersistent() {
	query := "SELECT key FROM a"
	rows, err := d.db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var domain string
		err = rows.Scan(&domain)
		if err != nil || checks.CheckDomain(domain) {
			continue
		}

		_, err := d.db.Exec("DELETE FROM a WHERE key = $1", domain)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if err = rows.Close(); err != nil {
		fmt.Println(err.Error())
	}

	query = "SELECT key FROM aaaa"
	rows, err = d.db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var domain string
		err = rows.Scan(&domain)
		if err != nil || checks.CheckDomain(domain) {
			continue
		}

		_, err = d.db.Exec("DELETE FROM aaaa WHERE key = $1", domain)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if err = rows.Close(); err != nil {
		fmt.Println(err.Error())
	}
}
