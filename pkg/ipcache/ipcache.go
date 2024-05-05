package ipcache

import (
	"database/sql"
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
	db    *sql.DB
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
	cache.SetPersistent(key, data)
}

func (cache *Cache) Get(key string) (data []net.IP, found bool) {
	cache.RLock()
	item, exists := cache.items[key]
	if exists {
		data = item.data
		found = true
	} else {
		data, found = cache.GetPersistent(key)
		if found {
			go cache.Set(key, data, 0)
		}
	}
	cache.RUnlock()
	return
}
