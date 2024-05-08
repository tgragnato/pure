package errcache

import (
	"net"
	"strings"
	"time"

	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/dohot"
)

func (errCache *ErrCache) cleanup(keys []string) {
	for i := range keys {
		if _, exist := errCache.cache6.Get(keys[i]); errCache.ipv6 && exist {
			go delete(errCache.items, keys[i])
			continue
		}
		if _, exist := errCache.cache4.Get(keys[i]); !errCache.ipv6 && exist {
			go delete(errCache.items, keys[i])
			continue
		}

		ips, cnames, ttl, err := dohot.DoH(keys[i], errCache.ipv6)
		if err == nil {
			for _, cname := range cnames {
				if !checks.CheckDomain(cname) &&
					!strings.HasSuffix(cname, "cloudfront.net.") &&
					!strings.HasSuffix(cname, "s3.amazonaws.com.") {
					go errCache.cache4.Set(keys[i], []net.IP{net.ParseIP("0.0.0.0")}, 0)
					go errCache.cache6.Set(keys[i], []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")}, 0)
					go delete(errCache.items, keys[i])
					continue
				}
			}

			if !errCache.geoChecks.CheckIPs(ips) {
				go errCache.cache4.Set(keys[i], []net.IP{net.ParseIP("0.0.0.0")}, 0)
				go errCache.cache6.Set(keys[i], []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")}, 0)
				go delete(errCache.items, keys[i])
				continue
			}

			if errCache.ipv6 {
				go errCache.cache6.Set(keys[i], ips, ttl)
			} else {
				go errCache.cache4.Set(keys[i], ips, ttl)
			}
			go delete(errCache.items, keys[i])
		} else {
			go func() {
				errCache.Lock()
				errCache.items[keys[i]].inc(errCache.duration)
				errCache.Unlock()
			}()
		}

		time.Sleep(time.Second)
	}
}

func (errCache *ErrCache) cleanupTimer() {
	for range time.NewTicker(errCache.duration).C {
		if errCache.cln {
			continue
		}
		errCache.cln = true
		cleanupkeys := []string{}
		errCache.RLock()
		for key := range errCache.items {
			if !errCache.items[key].expired() {
				continue
			}
			cleanupkeys = append(cleanupkeys, key)
		}
		errCache.RUnlock()
		errCache.cleanup(cleanupkeys)
		errCache.cln = false
	}
}
