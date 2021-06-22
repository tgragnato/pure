package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func retNull(m *dns.Msg, qName string) {
	m.Answer = nil
	newRR, err := dns.NewRR(fmt.Sprintf("%s A %s", qName, "0.0.0.0"))
	if err == nil {
		m.Answer = append(m.Answer, newRR)
	}
	newRR, err = dns.NewRR(fmt.Sprintf("%s AAAA %s", qName, "0000:0000:0000:0000:0000:0000:0000:0000"))
	if err == nil {
		m.Answer = append(m.Answer, newRR)
	}
}

func addIP(m *dns.Msg, qName string, ip []net.IP) {
	for i := range ip {
		rr, err := dns.NewRR(fmt.Sprintf("%s A %s", qName, ip[i].String()))
		if err == nil {
			m.Answer = append(m.Answer, rr)
		}
	}
}

func addIPv6(m *dns.Msg, qName string, ip []net.IP) {
	for i := range ip {
		rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", qName, ip[i].String()))
		if err == nil {
			m.Answer = append(m.Answer, rr)
		}
	}
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA:

			qName := Cloaking(q.Name)
			if qName == "tgragnato.it." {
				addIPv6(m, q.Name, []net.IP{net.ParseIP("fd76:abcd:ef90::")})
				addIP(m, q.Name, []net.IP{net.ParseIP("172.16.31.0")})
				return
			}

			if q.Qtype == dns.TypeAAAA {
				if data, found := cache6.Get(qName); found {
					addIPv6(m, q.Name, data)
					return
				}
				if _, found := cache4.Get(qName); found && analytics.data[strings.ToLower(q.Name[:len(q.Name)-1])].dns > 5 {
					return
				}
			} else {
				if data, found := cache4.Get(qName); found {
					addIP(m, q.Name, data)
					return
				}
				if _, found := cache6.Get(qName); found && analytics.data[strings.ToLower(q.Name[:len(q.Name)-1])].dns > 5 {
					return
				}
			}

			go preload.Push(qName)

			if !checkQuery(qName) {
				retNull(m, q.Name)
				return
			}

			var (
				ips    []net.IP
				cnames []string
				err    error
			)
			bsl := rand.Float64() * float64(len(trr))
			for i := 0; i < len(trr); i++ {
				index := (int(bsl) + i) % len(trr)
				ips, cnames, err = DoH(qName, trr[index], q.Qtype == dns.TypeAAAA)
				if err != nil {
					if err.Error() == "No IP addresses in response" {
						return
					} else {
						continue
					}
				} else {
					break
				}
			}

			for _, cname := range cnames {
				if !checkQuery(cname) && !strings.HasSuffix(cname, "cloudfront.net.") {
					retNull(m, q.Name)
					go cache4.Set(qName, []net.IP{net.ParseIP("0.0.0.0")})
					go cache6.Set(qName, []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")})
					return
				}
			}

			if !checkIPs(ips) {
				retNull(m, q.Name)
				go cache4.Set(qName, []net.IP{net.ParseIP("0.0.0.0")})
				go cache6.Set(qName, []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")})
				return
			}

			if q.Qtype == dns.TypeAAAA {
				go cache6.Set(qName, ips)
				addIPv6(m, q.Name, ips)
			} else {
				go cache4.Set(qName, ips)
				addIP(m, q.Name, ips)
			}
			go analytics.IncDNS(strings.ToLower(q.Name[:len(q.Name)-1]))
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
