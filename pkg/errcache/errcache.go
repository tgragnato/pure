package errcache

import (
	"sync"
	"time"

	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/ipcache"
)

type ErrCache struct {
	sync.RWMutex
	items     map[string]*ErrItem
	duration  time.Duration
	ipv6      bool
	cln       bool
	cache4    *ipcache.Cache
	cache6    *ipcache.Cache
	geoChecks *checks.GeoChecks
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
