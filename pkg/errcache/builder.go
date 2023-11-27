package errcache

import (
	"time"

	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/ipcache"
)

func NewErrCache(
	duration time.Duration,
	v6 bool,
	cache4 *ipcache.Cache,
	cache6 *ipcache.Cache,
	geoChecks *checks.GeoChecks,
) *ErrCache {

	errCache := &ErrCache{
		duration:  duration,
		items:     map[string]*ErrItem{},
		ipv6:      v6,
		cln:       false,
		cache4:    cache4,
		cache6:    cache6,
		geoChecks: geoChecks,
	}

	go errCache.cleanupTimer()

	return errCache
}
