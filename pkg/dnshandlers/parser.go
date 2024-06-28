package dnshandlers

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/tgragnato/pure/pkg/checks"
	"github.com/tgragnato/pure/pkg/dohot"
	"github.com/tgragnato/pure/pkg/errcache"
	"github.com/tgragnato/pure/pkg/ipcache"
)

func ParseQuery(
	m *dns.Msg,
	cache4 *ipcache.Cache,
	cache6 *ipcache.Cache,
	errCache4 *errcache.ErrCache,
	errCache6 *errcache.ErrCache,
	geoChecks *checks.GeoChecks,
	hintIPv4 net.IP,
	hintIPv6 net.IP,
) {
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

			if !checks.CheckDomain(q.Name) {
				retNull(m, q.Name)
				return
			}

			ips, cnames, ttl, err := dohot.DoH(q.Name, q.Qtype == dns.TypeAAAA)
			if err != nil {
				if q.Qtype == dns.TypeAAAA {
					go errCache6.Add(q.Name)
				} else {
					go errCache4.Add(q.Name)
				}
				m.SetRcode(m, dns.RcodeNameError)
				return
			}

			for _, cname := range cnames {
				if !checks.CheckDomain(cname) &&
					!strings.HasSuffix(cname, "cloudfront.net.") &&
					!strings.HasSuffix(cname, "s3.amazonaws.com.") {
					retNull(m, q.Name)
					go cache4.Set(q.Name, []net.IP{net.ParseIP("0.0.0.0")}, 0)
					go cache6.Set(q.Name, []net.IP{net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000")}, 0)
					return
				}
			}

			if !geoChecks.CheckIPs(ips) {
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
			addHTTPS(m, q.Name, hintIPv4, hintIPv6)

		case dns.TypeMX, dns.TypeTXT, dns.TypeSOA, dns.TypeNS, dns.TypeSVCB, dns.TypeSRV:
			_, found6 := cache6.Get(q.Name)
			_, found4 := cache4.Get(q.Name)
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
