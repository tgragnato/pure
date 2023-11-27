package ipcache

import (
	"net"
	"time"
)

type Item struct {
	data    []net.IP
	expires *time.Time
}

func (item *Item) expired() bool {
	if item.data == nil || item.expires == nil || len(item.data) == 0 {
		return true
	}
	if net.ParseIP("0.0.0.0").Equal(item.data[0]) ||
		net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").Equal(item.data[0]) {
		return false
	}
	return item.expires.Before(time.Now().Add(900 * time.Second))
}
