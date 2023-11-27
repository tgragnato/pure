package ipcache

import (
	"math"
	"time"

	"github.com/tgragnato/pure/pkg/dohot"
)

func (cache *Cache) cleanup(keys []string) {
	errors := 0
	for i := range keys {
		ips, _, ttl, err := dohot.DoH(keys[i], cache.ipv6)
		if err != nil {
			var exp_backoff time.Duration
			if errors < 7 {
				exp_backoff = 100 * time.Millisecond
			} else if errors < 16 {
				exp_backoff = time.Duration(int64(math.Exp2(float64(errors)))) * time.Millisecond
			} else if errors < 20 {
				exp_backoff = time.Minute
			} else {
				return
			}
			time.Sleep(exp_backoff)
			errors++
			continue
		}
		go cache.Set(keys[i], ips, ttl)
		if errors > 0 {
			errors--
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (cache *Cache) cleanupTimer() {
	duration := cache.ttl / 12
	if duration < time.Minute {
		duration = time.Minute
	}
	for range time.NewTicker(duration).C {
		if cache.cln {
			continue
		}
		cache.cln = true
		cleanupkeys := []string{}
		cache.RLock()
		for key := range cache.items {
			if !cache.items[key].expired() {
				continue
			}
			cleanupkeys = append(cleanupkeys, key)
		}
		cache.RUnlock()
		cache.cleanup(cleanupkeys)
		cache.cln = false
	}
}
