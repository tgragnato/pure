package dnshandlers

import (
	"fmt"

	"github.com/tgragnato/pure/pkg/dohot"
)

func (d *DnsHandlers) crossPrefetch() {
	query := `
		SELECT a.key
		FROM a
		LEFT JOIN aaaa ON aaaa.key = a.key
		WHERE aaaa.key IS NULL
		;
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			continue
		}

		ips, _, err := dohot.DoH(key, true)
		if err != nil {
			continue
		}
		d.setPersistent(key, ips, true)
	}
}

func (d *DnsHandlers) domainPrefetch() {
	query := `
		WITH domains AS (
			SELECT 
				a.key,
				reverse(
					split_part(reverse(a.key), '.', 1) || '.' ||
					split_part(reverse(a.key), '.', 2) || '.' ||
					split_part(reverse(a.key), '.', 3)
				) AS domain
			FROM a
		)
		SELECT domain
		FROM domains
		WHERE NOT EXISTS (SELECT 1 FROM a WHERE a.key = domain);
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			continue
		}

		ips, _, err := dohot.DoH(key, true)
		if err != nil {
			continue
		}
		d.setPersistent(key, ips, true)
	}
}

func (d *DnsHandlers) selfPrefetch(ipv6 bool) {
	var query string

	if ipv6 {

		if d.muIPv6 {
			return
		}

		query = `
			SELECT key
			FROM aaaa
			ORDER BY last_used ASC
			LIMIT 6000
			;
		`

		d.muIPv6 = true

	} else {

		if d.muIPv4 {
			return
		}

		query = `
			SELECT key
			FROM a
			ORDER BY last_used ASC
			LIMIT 6000
			;
		`

		d.muIPv4 = true
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var domain string
		err = rows.Scan(&domain)
		if err != nil {
			continue
		}

		ips, _, err := dohot.DoH(domain, ipv6)
		if err != nil {
			continue
		}

		d.setPersistent(domain, ips, ipv6)
	}

	if err = rows.Close(); err != nil {
		fmt.Println(err.Error())
	}

	if ipv6 {
		d.muIPv6 = false
	} else {
		d.muIPv4 = false
	}
}
