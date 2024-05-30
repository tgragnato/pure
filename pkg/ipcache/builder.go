package ipcache

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewCache(duration time.Duration, v6 bool, dsn string) *Cache {
	db, err := sql.Open("pgx", dsn)
	if err != nil || (db != nil && db.Ping() != nil) {
		db = nil
	}

	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		ipv6:  v6,
		cln:   false,
		db:    db,
	}

	cache.Prefetch()
	go cache.cleanupTimer()

	return cache
}
