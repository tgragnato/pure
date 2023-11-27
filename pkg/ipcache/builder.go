package ipcache

import (
	"time"
)

func NewCache(duration time.Duration, v6 bool) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		ipv6:  v6,
		cln:   false,
	}

	go cache.cleanupTimer()

	return cache
}
