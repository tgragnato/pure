package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var preload = NewPreload(300*time.Second, "/etc/proxy/preload.names")

type Preload struct {
	sync.Mutex
	data     map[string]bool
	conf     string
	duration time.Duration
}

func (p *Preload) Load() {
	domains := populateCheck(p.conf)

	p.Lock()
	for i := range domains {
		p.data[domains[i]] = true
	}
	p.Unlock()

	for i := range domains {
		if !checkQuery(domains[i]) {
			continue
		}

		ip4, cname4, err4 := DoH(domains[i], "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion", false)
		ip6, cname6, err6 := DoH(domains[i], "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion", true)
		time.Sleep(time.Second / 4)

		if err4 == nil && err6 == nil {
			ok := checkIPs(ip4) && checkIPs(ip6)
			for t := range cname4 {
				ok = ok && checkQuery(cname4[t])
			}
			for t := range cname6 {
				ok = ok && checkQuery(cname6[t])
			}

			if ok {
				go cache4.Set(domains[i], ip4)
				go cache6.Set(domains[i], ip6)
			}

		} else if err4 == nil && err6 != nil {
			ok := checkIPs(ip4)
			for t := range cname4 {
				ok = ok && checkQuery(cname4[t])
			}

			if ok {
				go cache4.Set(domains[i], ip4)
			}

		} else if err4 != nil && err6 == nil {
			ok := checkIPs(ip6)
			for t := range cname6 {
				ok = ok && checkQuery(cname6[t])
			}

			if ok {
				go cache6.Set(domains[i], ip6)
			}
		}
	}
}

func (p *Preload) Write() {
	f, err := os.OpenFile(p.conf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Error opening file %s during preload write", p.conf)
		log.Printf("   Printing error: %s", err.Error())
	}
	defer f.Close()

	for x := range p.data {
		if p.data[x] {
			continue
		}

		if _, err := f.WriteString(x + "\n"); err != nil {
			log.Printf("Error writing to file %s during preload write", p.conf)
			log.Printf("   Printing error: %s", err.Error())
			continue
		}

		p.data[x] = true
	}
}

func (p *Preload) Push(domain string) {
	if strings.HasSuffix(domain, "googleapis.com.") {
		return
	}
	p.Lock()
	if _, found := p.data[domain]; !found {
		p.data[domain] = false
	}
	p.Unlock()
}

func (p *Preload) PeriodicWrite() {
	ticker := time.Tick(p.duration)
	for {
		select {
		case <-ticker:
			p.Write()
		}
	}
}

func NewPreload(dur time.Duration, configuration string) *Preload {
	preload := &Preload{
		conf:     configuration,
		duration: dur,
		data:     map[string]bool{},
	}
	go preload.Load()
	go preload.PeriodicWrite()
	return preload
}
