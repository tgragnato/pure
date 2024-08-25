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

			if q.Qtype == dns.TypeAAAA {
				if data, found := d.getPersistent(q.Name, true); found {
					addIPv6(m, q.Name, data)
					return
				}
			} else {
				if data, found := d.getPersistent(q.Name, false); found {
					addIP(m, q.Name, data)
					return
				}
			}

			if !checks.CheckDomain(q.Name) {
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
				go d.setPersistent(q.Name, ips, true)
				addIPv6(m, q.Name, ips)
			} else {
				go d.setPersistent(q.Name, ips, false)
				addIP(m, q.Name, ips)
			}

		case dns.TypeHTTPS:
			addHTTPS(m, q.Name, d.hintIPv4, d.hintIPv6)

		case dns.TypeMX, dns.TypeTXT, dns.TypeSOA, dns.TypeNS, dns.TypeSVCB, dns.TypeSRV:
			_, found6 := d.getPersistent(q.Name, true)
			_, found4 := d.getPersistent(q.Name, false)
			isApple := strings.HasPrefix(q.Name, ".apple.com.") || strings.HasPrefix(q.Name, ".icloud.com.")
			isAppleAkamai := strings.HasPrefix(q.Name, "apple.com.akadns.net.")
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
