package errcache

import "time"

type ErrItem struct {
	counter uint
	expires *time.Time
}

func (item *ErrItem) expired() bool {
	if item.counter == 0 || item.expires == nil {
		return true
	}

	return item.expires.Before(time.Now())
}

func (item *ErrItem) inc(period time.Duration) {
	item.counter += 1
	expiration := time.Now().Add(period * time.Duration(item.counter))
	item.expires = &expiration
}
