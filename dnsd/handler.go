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
				if found := errCache6.Exist(q.Name); found {
					m.SetRcode(m, dns.RcodeNameError)
					return
				}
			} else {
				if data, found := cache4.Get(q.Name); found {
					addIP(m, q.Name, data)
					return
				}
				if found := errCache4.Exist(q.Name); found {
					m.SetRcode(m, dns.RcodeNameError)
					return
				}
			}

			if !checkDomain(q.Name) {
				retNull(m, q.Name)
				return
			}

			ips, cnames, ttl, err := DoH(q.Name, q.Qtype == dns.TypeAAAA)
			if err != nil {
				log.Println(q.Name + ": " + err.Error())
				if q.Qtype == dns.TypeAAAA {
					go errCache6.Add(q.Name)
				} else {
					go errCache4.Add(q.Name)
				}
				m.SetRcode(m, dns.RcodeNameError)
				return
			} else {
				log.Println(q.Name + ": First resolution")
			}

			for _, cname := range cnames {
				if !checkDomain(cname) && !strings.HasSuffix(cname, "cloudfront.net.") {
					retNull(m, q.Name)
					go cache4.Set(q.Name, []net.IP{net.ParseIP(nullIPv4)}, 0)
					go cache6.Set(q.Name, []net.IP{net.ParseIP(nullIPv6)}, 0)
					return
				}
			}

			if !checkIPs(ips) {
				retNull(m, q.Name)
				go cache4.Set(q.Name, []net.IP{net.ParseIP(nullIPv4)}, 0)
				go cache6.Set(q.Name, []net.IP{net.ParseIP(nullIPv6)}, 0)
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
