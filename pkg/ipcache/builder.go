package ipcache

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewCache(duration time.Duration, v6 bool, dsn string) *Cache {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err.Error())
		db = nil
	}

	if db != nil && db.Ping() != nil {
		db = nil
	}

	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
		ipv6:  v6,
		cln:   false,
		db:    db,
	}

	go cache.cleanupTimer()

	return cache
}
