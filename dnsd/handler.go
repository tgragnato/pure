package main

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

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
				m.SetRcode(m, dns.RcodeNXRrset)
				return
			} else {
				log.Println(q.Name + ": First resolution")
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

		case dns.TypeHTTPS:
			addHTTPS(m, q.Name)

		default:
			m.SetRcode(m, dns.RcodeNotImplemented)
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
