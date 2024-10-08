package dnshandlers

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestDnsHandlers_ParseQuery(t *testing.T) {
	t.Parallel()

	const domain = "example.com."

	d := &DnsHandlers{
		db:        newDb(t),
		geoChecks: nil,
		hintIPv4:  net.ParseIP("192.0.2.3").To4(),
		hintIPv6:  net.ParseIP("2001:db8::3"),
	}
	d.setPersistent(domain, []net.IP{
		net.ParseIP("192.0.2.1").To4(),
		net.ParseIP("203.0.113.1").To4(),
	}, false)
	d.setPersistent(domain, []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
	}, true)

	t.Run("Test query A", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  domain,
					Qtype: dns.TypeA,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Fatalf("Expected 2 answers, got %d", len(m.Answer))
		}
	})

	t.Run("Test query AAAA", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  domain,
					Qtype: dns.TypeAAAA,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Fatalf("Expected 2 answers, got %d", len(m.Answer))
		}
	})

	t.Run("Test prohibited domain query A", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  "testprohibited.tgragnato.it.",
					Qtype: dns.TypeA,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Fatalf("Expected 2 answers, got %d", len(m.Answer))
		}
		if m.Answer[0].Header().Rrtype != dns.TypeA {
			t.Fatalf("Expected A record, got %d", m.Answer[0].Header().Rrtype)
		}
		if m.Answer[0].String() != "testprohibited.tgragnato.it.	3600	IN	A	0.0.0.0" {
			t.Fatalf("Expected null record, got %s", m.Answer[0].String())
		}
	})

	t.Run("Test prohibited domain query AAAA", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  "testprohibited.tgragnato.it.",
					Qtype: dns.TypeAAAA,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Fatalf("Expected 2 answers, got %d", len(m.Answer))
		}
		if m.Answer[1].Header().Rrtype != dns.TypeAAAA {
			t.Fatalf("Expected A record, got %d", m.Answer[0].Header().Rrtype)
		}
		if m.Answer[1].String() != "testprohibited.tgragnato.it.	3600	IN	AAAA	::" {
			t.Fatalf("Expected null record, got %s", m.Answer[1].String())
		}
	})

	t.Run("Test HTTPS query", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  domain,
					Qtype: dns.TypeHTTPS,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 1 {
			t.Fatalf("Expected 1 answer, got %d", len(m.Answer))
		}
		if m.Answer[0].Header().Rrtype != dns.TypeHTTPS {
			t.Fatalf("Expected HTTPS record, got %d", m.Answer[0].Header().Rrtype)
		}

		if m.Answer[0].String() != "example.com.	86400	IN	HTTPS	1 . alpn=\"h3,h2\" mandatory=\"alpn\" ipv4hint=\"<nil>\" ipv6hint=\"<nil>\"" {
			t.Fatalf("Expected HTTPS record, got %s", m.Answer[0].String())
		}
	})

	t.Run("Test CNAME query", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  "cname.example.com.",
					Qtype: dns.TypeCNAME,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 0 {
			t.Fatalf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeNotImplemented {
			t.Fatalf("Expected Not Implemented, got %d", m.Rcode)
		}
	})

	t.Run("Test prohibited domain query TXT", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion.",
					Qtype: dns.TypeTXT,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 0 {
			t.Fatalf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeRefused {
			t.Fatalf("Expected Refused, got %d", m.Rcode)
		}
	})

	t.Run("Test uncached domain query MX", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  "tgragnato.it.",
					Qtype: dns.TypeMX,
				},
			},
		}
		d.ParseQuery(m)
		if len(m.Answer) != 0 {
			t.Fatalf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeRefused {
			t.Fatalf("Expected Refused, got %d", m.Rcode)
		}
	})
}
