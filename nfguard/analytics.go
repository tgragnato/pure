package main

import (
	"log"
	"sync"
	"time"
)

var analytics = &Analytics{
	Unmodified: 0,
	Modified:   0,
}

type Analytics struct {
	sync.RWMutex
	Unmodified uint64
	Modified   uint64
}

func (a *Analytics) IncUnmodified() {
	a.Lock()
	a.Unmodified += 1
	a.Unlock()
}

func (a *Analytics) IncModified() {
	a.Lock()
	a.Modified += 1
	a.Unlock()
}

func (a *Analytics) Report() {
	for range time.NewTicker(time.Hour).C {
		a.RLock()
		log.Printf("Unmodified: %d - Modified: %d \n", a.Unmodified, a.Modified)
		a.RUnlock()
	}
}
