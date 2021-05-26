package main

import (
	"net"
	"sync"
	"time"
)

var (
	cache4 = NewCache(3600*time.Second, false)
	cache6 = NewCache(3600*time.Second, true)
)

type Item struct {
	sync.RWMutex
	data    []net.IP
	expires *time.Time
}

func (item *Item) touch(duration time.Duration) {
	item.Lock()
	expiration := time.Now().Add(duration)
	item.expires = &expiration
	item.Unlock()
}

func (item *Item) expired() bool {
	var value bool
	item.RLock()
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	if net.ParseIP("0.0.0.0").Equal(item.data[0]) {
		value = false
	}
	if net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").Equal(item.data[0]) {
		value = false
	}
	item.RUnlock()
	return value
}

type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
	ipv6  bool
}

func (cache *Cache) Set(key string, data net.IP) {
	cache.mutex.Lock()
	if _, keyexists := cache.items[key]; keyexists {
		duplicated := false
		for i := range cache.items[key].data {
			if cache.items[key].data[i].Equal(data) {
				duplicated = true
			}
		}
		if !duplicated {
			item := &Item{data: append(cache.items[key].data, data)}
			item.touch(cache.ttl)
			cache.items[key] = item
		}
	} else {
		item := &Item{data: []net.IP{data}}
		item.touch(cache.ttl)
		cache.items[key] = item
	}
	cache.mutex.Unlock()
}

func (cache *Cache) Get(key string) (data []net.IP, found bool) {
	cache.mutex.RLock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = []net.IP{}
		found = false
	} else {
		data = item.data
		found = true
	}
	cache.mutex.RUnlock()
	return
}

func (cache *Cache) cleanup() {
	for key := range cache.items {
		if !cache.items[key].expired() {
			continue
		}
		ips, _, err := DoH(key, "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion", cache.ipv6)
		if err != nil {
			continue
		}
		cache.mutex.Lock()
		item := &Item{data: ips}
		item.touch(cache.ttl)
		cache.items[key] = item
		cache.mutex.Unlock()
	}
}

func (cache *Cache) startCleanupTimer() {
	duration := cache.ttl / 4
	if duration < time.Second {
		duration = time.Second
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
	}
	cache.startCleanupTimer()
	return cache
}
