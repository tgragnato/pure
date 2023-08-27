package main

import (
	"fmt"
	"log"
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

			if q.Qtype == dns.TypeAAAA {
				if data, found := cache6.Get(q.Name); found {
					addIPv6(m, q.Name, data)
					return
				}
			} else {
				if data, found := cache4.Get(q.Name); found {
					addIP(m, q.Name, data)
					return
				}
			}

			if !CheckDomain(q.Name) {
				retNull(m, q.Name)
				return
			}

			ips, cnames, ttl, err := DoH(q.Name, q.Qtype == dns.TypeAAAA)
			if err != nil {
				log.Println(q.Name + ": " + err.Error())
				return
			}

			for _, cname := range cnames {
				if !CheckDomain(cname) && !strings.HasSuffix(cname, "cloudfront.net.") {
					retNull(m, q.Name)
					go cache4.Set(q.Name, []net.IP{net.ParseIP("0.0.0.0")}, 0)
					go cache6.Set(q.Name, []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")}, 0)
					return
				}
			}

			if !checkIPs(ips) {
				retNull(m, q.Name)
				go cache4.Set(q.Name, []net.IP{net.ParseIP("0.0.0.0")}, 0)
				go cache6.Set(q.Name, []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")}, 0)
				return
			}

			if q.Qtype == dns.TypeAAAA {
				go cache6.Set(q.Name, ips, ttl)
				addIPv6(m, q.Name, ips)
			} else {
				go cache4.Set(q.Name, ips, ttl)
				addIP(m, q.Name, ips)
			}
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
