package ipcache

import (
	"log"
	"net"
	"strings"

	"github.com/tgragnato/pure/pkg/checks"
)

func (cache *Cache) Prefetch() {
	if cache.db == nil {
		return
	}

	query := `
		SELECT key, value
		FROM a
		ORDER BY last_used ASC
		;
	`
	if cache.ipv6 {
		query = `
		SELECT key, value
		FROM aaaa
		ORDER BY last_used ASC
		;
	`
	}

	rows, err := cache.db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer rows.Close()

	geo := checks.NewGeoChecks()

	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		data := strings.Split(value, ",")
		ips := make([]net.IP, 0, len(data))
		for _, ip := range data {
			ips = append(ips, net.ParseIP(ip))
		}

		if !checks.CheckDomain(key) || geo.CheckIPs(ips) {
			cache.DeletePersistent(key)
			continue
		}

		cache.Set(key, ips, 0)
	}
}
