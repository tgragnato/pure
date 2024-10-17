package dnshandlers

import (
	"strings"

	"github.com/miekg/dns"
	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/dohot"
)

func (d *DnsHandlers) ParseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA, dns.TypeAAAA:

			qName := strings.ToLower(q.Name)

			if q.Qtype == dns.TypeAAAA {
				if data, found := d.getPersistent(qName, true); found {
					addIPv6(m, q.Name, data)
					return
				}
			} else {
				if data, found := d.getPersistent(qName, false); found {
					addIP(m, q.Name, data)
					return
				}
			}

			if !checks.CheckDomain(qName) {
				retNull(m, q.Name)
				return
			}

			ips, cnames, err := dohot.DoH(q.Name, q.Qtype == dns.TypeAAAA)
			if err != nil {
				m.SetRcode(m, dns.RcodeNameError)
				return
			}

			for _, cname := range cnames {
				if !checks.CheckDomain(cname) &&
					!strings.HasSuffix(cname, "cloudfront.net.") &&
					!strings.HasSuffix(cname, "s3.amazonaws.com.") {
					retNull(m, q.Name)
					return
				}
			}

			if !d.geoChecks.CheckIPs(ips) {
				retNull(m, q.Name)
				return
			}

			if q.Qtype == dns.TypeAAAA {
				go d.setPersistent(qName, ips, true)
				addIPv6(m, q.Name, ips)
			} else {
				go d.setPersistent(qName, ips, false)
				addIP(m, q.Name, ips)
			}

		case dns.TypeHTTPS:
			qName := strings.ToLower(q.Name)

			if !checks.CheckDomain(qName) {
				m.SetRcode(m, dns.RcodeRefused)
				return
			}

			hintIPv4, _ := d.getPersistent(qName, true)
			hintIPv6, _ := d.getPersistent(qName, false)
			addHTTPS(m, q.Name, hintIPv4, hintIPv6)

		case dns.TypeMX, dns.TypeTXT, dns.TypeSOA, dns.TypeNS, dns.TypeSVCB, dns.TypeSRV:
			qName := strings.ToLower(q.Name)

			if !checks.CheckDomain(qName) {
				m.SetRcode(m, dns.RcodeRefused)
				return
			}

			_, found6 := d.getPersistent(qName, true)
			_, found4 := d.getPersistent(qName, false)
			isApple := strings.HasPrefix(qName, ".apple.com.") || strings.HasPrefix(qName, ".icloud.com.")
			isAppleAkamai := strings.HasPrefix(qName, "apple.com.akadns.net.")
			if found4 || found6 || isApple || isAppleAkamai {
				answer, err := dohot.Transparent(q.Name, q.Qtype, isApple)
				if err != nil {
					m.SetRcode(m, dns.RcodeServerFailure)
				} else {
					m.Answer = answer
					m.SetRcode(m, dns.RcodeSuccess)
				}
			} else {
				m.SetRcode(m, dns.RcodeRefused)
			}

		default:
			m.SetRcode(m, dns.RcodeNotImplemented)
		}
	}
}
