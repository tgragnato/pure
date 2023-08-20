package main

import (
	"log"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

var analytics = &Analytics{
	data:  map[string]uint{},
	exact: map[string]uint{},
}

type Analytics struct {
	sync.Mutex
	data  map[string]uint
	exact map[string]uint
}

type Domain struct {
	domain  string
	counter uint
}

type DomainList []Domain

func (a *Analytics) Inc(dName string) {
	split := strings.Split(dName, ".")
	var truncated string
	if len(split) > 1 {
		truncated = split[len(split)-2] + "." + split[len(split)-1]
	} else {
		truncated = split[len(split)-1]
	}

	a.Lock()
	_, exist := a.data[truncated]
	if exist {
		a.data[truncated]++
	} else {
		a.data[truncated] = 1
	}
	_, exist = a.exact[dName]
	if exist {
		a.exact[dName]++
	} else {
		a.exact[dName] = 1
	}
	a.Unlock()
}

func (a *Analytics) Report() {
	for {
		select {
		case <-time.NewTicker(time.Hour / 2).C:
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
		case <-time.NewTicker(6 * time.Hour).C:
			a.Lock()
			log.Printf("%v\n", a.exact)
			a.Unlock()
		}
	}
}

func (a *Analytics) PreloadDNS() {
	for {
		select {
		case <-time.NewTicker(15 * time.Minute).C:
			a.Lock()
			for domain := range a.exact {
				go net.LookupIP(domain)
			}
			a.Unlock()
		}
	}
}
