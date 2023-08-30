package main

import (
	"net"
	"strings"
	"sync"
	"time"
)

var (
	errCache4 = newErrCache(time.Minute, false)
	errCache6 = newErrCache(time.Minute, true)
)

type ErrItem struct {
	counter uint
	expires *time.Time
}

func (item *ErrItem) expired() bool {
	if item.counter == 0 || item.expires == nil {
		return true
	}

	return item.expires.Before(time.Now())
}

func (item *ErrItem) inc(period time.Duration) {
	item.counter += 1
	expiration := time.Now().Add(period * time.Duration(item.counter))
	item.expires = &expiration
}

type ErrCache struct {
	sync.RWMutex
	items    map[string]*ErrItem
	duration time.Duration
	ipv6     bool
	cln      bool
}

func (errCache *ErrCache) Add(key string) {
	errCache.Lock()
	expiration := time.Now().Add(errCache.duration)
	errCache.items[key] = &ErrItem{
		counter: 0,
		expires: &expiration,
	}
	errCache.Unlock()
}

func (errCache *ErrCache) Del(key string) {
	errCache.Lock()
	if _, exist := errCache.items[key]; exist {
		delete(errCache.items, key)
	}
	errCache.Unlock()
}

func (errCache *ErrCache) Exist(key string) bool {
	errCache.RLock()
	_, exist := errCache.items[key]
	errCache.RUnlock()
	return exist
}

func (errCache *ErrCache) Cleanup(keys []string) {
	for i := range keys {
		if _, exist := cache6.Get(keys[i]); errCache.ipv6 && exist {
			go errCache.Del(keys[i])
			continue
		}
		if _, exist := cache4.Get(keys[i]); !errCache.ipv6 && exist {
			go errCache.Del(keys[i])
			continue
		}

		ips, cnames, ttl, err := DoH(keys[i], errCache.ipv6)
		if err == nil {
			for _, cname := range cnames {
				if !checkDomain(cname) && !strings.HasSuffix(cname, "cloudfront.net.") {
					go cache4.Set(keys[i], []net.IP{net.ParseIP(nullIPv4)}, 0)
					go cache6.Set(keys[i], []net.IP{net.ParseIP(nullIPv6)}, 0)
					go errCache.Del(keys[i])
					continue
				}
			}

			if !checkIPs(ips) {
				go cache4.Set(keys[i], []net.IP{net.ParseIP(nullIPv4)}, 0)
				go cache6.Set(keys[i], []net.IP{net.ParseIP(nullIPv6)}, 0)
				go errCache.Del(keys[i])
				continue
			}

			if errCache.ipv6 {
				go cache6.Set(keys[i], ips, ttl)
			} else {
				go cache4.Set(keys[i], ips, ttl)
			}
			go errCache.Del(keys[i])
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

func (errCache *ErrCache) CleanupTimer() {
	for {
		select {
		case <-time.NewTicker(errCache.duration).C:
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
			errCache.Cleanup(cleanupkeys)
			errCache.cln = false
		}
	}
}

func newErrCache(duration time.Duration, v6 bool) *ErrCache {
	errCache := &ErrCache{
		duration: duration,
		items:    map[string]*ErrItem{},
		ipv6:     v6,
		cln:      false,
	}
	go errCache.CleanupTimer()
	return errCache
}
