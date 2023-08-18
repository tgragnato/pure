package main

import (
	"log"
	"net"
	"sort"
	"sync"
	"time"
)

var analytics = &Analytics{data: map[string]uint{}}

type Analytics struct {
	sync.Mutex
	data map[string]uint
}

type Domain struct {
	domain  string
	counter uint
}

type DomainList []Domain

func (a *Analytics) Inc(dName string) {
	a.Lock()
	_, exist := a.data[dName]
	if exist {
		a.data[dName]++
	} else {
		a.data[dName] = 1
	}
	a.Unlock()
}

func (a *Analytics) Report() {
	for {
		select {
		case <-time.NewTicker(time.Hour).C:
			a.Lock()
			dl := make(DomainList, len(a.data))
			i := 0
			for domain, counter := range a.data {
				dl[i] = Domain{domain, counter}
				i++
			}
			a.Unlock()
			sort.Slice(dl, func(i, j int) bool {
				return dl[i].counter < dl[j].counter
			})
			for _, orderedLogItem := range dl {
				log.Printf("%s, %d\n", orderedLogItem.domain, orderedLogItem.counter)
			}
		}
	}
}

func (a *Analytics) PreloadDNS() {
	for {
		select {
		case <-time.NewTicker(15 * time.Minute).C:
			a.Lock()
			for domain := range a.data {
				go net.LookupIP(domain)
			}
			a.Unlock()
		}
	}
}
