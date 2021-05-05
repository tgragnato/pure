package main

import (
	"net"
	"sync"
	"time"
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
	item.RUnlock()
	return value
}

type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
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
	cache.mutex.Lock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = []net.IP{}
		found = false
	} else {
		item.touch(cache.ttl)
		data = item.data
		found = true
	}
	cache.mutex.Unlock()
	return
}

func (cache *Cache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *Cache) startCleanupTimer() {
	duration := cache.ttl
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

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
	}
	cache.startCleanupTimer()
	return cache
}
