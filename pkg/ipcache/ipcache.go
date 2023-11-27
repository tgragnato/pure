package ipcache

import (
	"net"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
	ipv6  bool
	cln   bool
}

func (cache *Cache) Set(key string, data []net.IP, ttl uint32) {
	cache.Lock()
	duration := time.Duration(ttl) * time.Second
	if duration < cache.ttl {
		duration = cache.ttl
	}
	expiration := time.Now().Add(duration)
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
