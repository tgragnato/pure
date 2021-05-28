package main

import (
	"log"
	"math"
	"net"
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
	if item.data == nil || item.expires == nil {
		return true
	}
	if net.ParseIP("0.0.0.0").Equal(item.data[0]) ||
		net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").Equal(item.data[0]) {
		return false
	}
	return item.expires.Before(time.Now().Add(900 * time.Second))
}

type Cache struct {
	ttl   time.Duration
	items map[string]*Item
	ipv6  bool
	fetch bool
}

func (cache *Cache) Set(key string, data []net.IP) {
	expiration := time.Now().Add(cache.ttl)
	cache.items[key] = &Item{
		data:    data,
		expires: &expiration,
	}
}

func (cache *Cache) Get(key string) (data []net.IP, found bool) {
	item, exists := cache.items[key]
	if exists {
		data = item.data
		found = true
	} else {
		data = []net.IP{}
		found = false
	}
	return
}

func (cache *Cache) cleanup() {
	if cache.fetch {
		return
	}
	cache.fetch = true

	IPvX := "IPv4"
	if cache.ipv6 {
		IPvX = "IPv6"
	}
	log.Printf("Info: starting cache cleanup, %d items in cache (%s)", len(cache.items), IPvX)
	errors := 0
	counter := 0

	for key := range cache.items {

		if !cache.items[key].expired() {
			continue
		}

		ips, _, err := DoH(key, "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion", cache.ipv6)
		if err != nil {
			log.Printf("Error during cleanup for query name: %s (%s)", key, IPvX)
			log.Printf("   Printing error: %s", err.Error())

			var exp_backoff time.Duration
			if errors < 7 {
				exp_backoff = 100 * time.Millisecond
			} else if errors < 16 {
				exp_backoff = time.Duration(int64(math.Exp2(float64(errors)))) * time.Millisecond
			} else if errors < 20 {
				exp_backoff = time.Minute
			} else {
				log.Printf("Too many consecutive errors during cleanup, aborting %s", IPvX)
				cache.fetch = false
				return
			}

			time.Sleep(exp_backoff)
			errors++
			continue
		}

		go cache.Set(key, ips)
		if errors > 0 {
			errors--
		}
		counter++
		time.Sleep(100 * time.Millisecond)
	}

	cache.fetch = false
	log.Printf("Info: cache cleanup completed, %d items updated (%s)", counter, IPvX)
}

func (cache *Cache) startCleanupTimer() {
	duration := cache.ttl / 12
	if duration < time.Minute {
		duration = time.Minute
	}
	ticker := time.Tick(duration)
	go (func() {
		for {
			select {
			case <-ticker:
				cache.cleanup()
			}
		}
	})()
}

func NewCache(duration time.Duration, v6 bool) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		ipv6:  v6,
		fetch: false,
	}
	cache.startCleanupTimer()
	return cache
}
