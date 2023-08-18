package main

import (
	"log"
	"sync"
	"time"
)

var analytics = &Analytics{data: map[string]uint{}}

type Analytics struct {
	sync.Mutex
	data map[string]uint
}

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
			log.Println(analytics.data)
			a.Unlock()
		}
	}
}
