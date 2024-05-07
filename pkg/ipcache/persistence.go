package ipcache

import (
	"log"
	"net"
	"strings"
)

func (cache *Cache) SetPersistent(key string, data []net.IP) {
	if cache.db == nil ||
		len(data) == 0 ||
		strings.HasSuffix(key, "googlevideo.com.") {
		return
	}

	serialize := ""
	for _, ip := range data {
		if net.ParseIP("0.0.0.0").Equal(ip) ||
			net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").Equal(ip) {
			return
		}
		serialize += ip.String() + ","
	}

	query := `
		INSERT INTO a (key, value, discovered_on, last_used)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, last_used = CURRENT_TIMESTAMP
		;
	`
	if cache.ipv6 {
		query = `
		INSERT INTO aaaa (key, value, discovered_on, last_used)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, last_used = CURRENT_TIMESTAMP
		;
	`
	}

	_, err := cache.db.Exec(query, key, serialize)
	if err != nil {
		log.Println(err.Error())
	}
}

func (cache *Cache) GetPersistent(key string) (data []net.IP, found bool) {
	if cache.db == nil {
		return []net.IP{}, false
	}

	query := `
		SELECT value
		FROM a
		WHERE key = $1
	`
	if cache.ipv6 {
		query = `
		SELECT value
		FROM aaaa
		WHERE key = $1
	`
	}

	row := cache.db.QueryRow(query, key)
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

func (cache *Cache) DeletePersistent(key string) {
	if cache.db == nil {
		return
	}

	query := `
		DELETE FROM a
		WHERE key = $1
	`
	if cache.ipv6 {
		query = `
		DELETE FROM aaaa
		WHERE key = $1
	`
	}

	_, err := cache.db.Exec(query, key)
	if err != nil {
		log.Println(err.Error())
	}
}
