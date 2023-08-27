package main

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

var analytics = &Analytics{
	data: map[string]uint{},
}

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
