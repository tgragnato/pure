package dnshandlers

import (
	"fmt"
	"net"

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

func addHTTPS(
	m *dns.Msg,
	qName string,
	hintIPv4 []net.IP,
	hintIPv6 []net.IP,
) {
	https := &dns.SVCB{
		Hdr: dns.RR_Header{
			Name:   qName,
			Rrtype: dns.TypeHTTPS,
			Class:  dns.ClassINET,
			Ttl:    86400,
		},
		Priority: 1,
		Target:   ".",
		Value: []dns.SVCBKeyValue{
			&dns.SVCBAlpn{
				Alpn: []string{"h3", "h2"},
			},
			&dns.SVCBMandatory{
				Code: []dns.SVCBKey{
					dns.SVCB_ALPN,
				},
			},
		},
	}

	if len(hintIPv4) != 0 {
		https.Value = append(https.Value, &dns.SVCBIPv4Hint{
			Hint: hintIPv4,
		})
	}

	if len(hintIPv6) != 0 {
		https.Value = append(https.Value, &dns.SVCBIPv6Hint{
			Hint: hintIPv6,
		})
	}

	m.Answer = append(m.Answer, https)
}
