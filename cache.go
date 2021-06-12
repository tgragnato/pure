package main

import (
	"math"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	cache4 = NewCache(3600*time.Second, false)
	cache6 = NewCache(3600*time.Second, true)
)

type Item struct {
	data    []net.IP
	expires *time.Time
}

func (item *Item) expired() bool {
	if item.data == nil || item.expires == nil || len(item.data) == 0 {
		return true
	}
	if net.ParseIP("0.0.0.0").Equal(item.data[0]) ||
		net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").Equal(item.data[0]) {
		return false
	}
	return item.expires.Before(time.Now().Add(900 * time.Second))
}

type Cache struct {
	sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
	ipv6  bool
}

func (cache *Cache) Set(key string, data []net.IP) {
	cache.Lock()
	expiration := time.Now().Add(cache.ttl)
	cache.items[key] = &Item{
		data:    data,
		expires: &expiration,
	}
	cache.Unlock()
}

func (cache *Cache) Get(key string) (data []net.IP, found bool) {
	cache.RLock()
	item, exists := cache.items[key]
	if exists {
		data = item.data
		found = true
	} else {
		data = []net.IP{}
		found = false
	}
	cache.RUnlock()
	return
}

func (cache *Cache) Cleanup(keys []string) {
	errors := 0
	for i := range keys {
		ips, _, err := DoH(keys[i], "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion", cache.ipv6)
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
		go cache.Set(keys[i], ips)
		if errors > 0 {
			errors--
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (cache *Cache) CleanupTimer() {
	duration := cache.ttl / 12
	if duration < time.Minute {
		duration = time.Minute
	}
	ticker := time.Tick(duration)
	for {
		select {
		case <-ticker:
			cleanupkeys := []string{}
			cache.RLock()
			for key := range cache.items {
				if strings.HasSuffix(key, "googleapis.com") {
					continue
				}
				if !cache.items[key].expired() {
					continue
				}
				cleanupkeys = append(cleanupkeys, key)
			}
			cache.RUnlock()
			cache.Cleanup(cleanupkeys)
		}
	}
}

func NewCache(duration time.Duration, v6 bool) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		ipv6:  v6,
	}
	go cache.CleanupTimer()
	return cache
}
