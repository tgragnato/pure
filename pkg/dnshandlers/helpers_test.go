package dnshandlers

import (
	"net"
	"reflect"
	"testing"

	"github.com/miekg/dns"
)

func Test_retNull(t *testing.T) {
	t.Parallel()

	m := new(dns.Msg)
	retNull(m, "example.com.")

	if len(m.Answer) != 2 {
		t.Fatalf("expected 2 RRs, got %d", len(m.Answer))
	}
	if m.Answer[0].Header().Rrtype != dns.TypeA {
		t.Errorf("expected first type to be A")
	}
	if m.Answer[1].Header().Rrtype != dns.TypeAAAA {
		t.Errorf("expected second type to be AAAA")
	}

	if m.Answer[0].Header().Name != "example.com." {
		t.Errorf("expected name to be example.com")
	}
}

func Test_addIP(t *testing.T) {
	t.Parallel()

	m := new(dns.Msg)
	qName := "example.com."
	ip := []net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("203.0.113.1"),
	}

	addIP(m, qName, ip)

	expectedRRs := []dns.RR{
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			A: net.ParseIP("192.0.2.1"),
		},
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			A: net.ParseIP("203.0.113.1"),
		},
	}

	if len(m.Answer) != len(expectedRRs) {
		t.Fatalf("expected %d RRs, got %d", len(expectedRRs), len(m.Answer))
	}

	for i := range m.Answer {
		if !reflect.DeepEqual(m.Answer[i], expectedRRs[i]) {
			t.Errorf("expected RR %d to be %v, got %v", i, expectedRRs[i], m.Answer[i])
		}
	}
}

func Test_addIPv6(t *testing.T) {
	t.Parallel()

	m := new(dns.Msg)
	qName := "example.com."
	ip := []net.IP{
		net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7335"),
	}

	addIPv6(m, qName, ip)

	expectedRRs := []dns.RR{
		&dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			AAAA: net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		},
		&dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			AAAA: net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7335"),
		},
	}

	if len(m.Answer) != len(expectedRRs) {
		t.Fatalf("expected %d RRs, got %d", len(expectedRRs), len(m.Answer))
	}

	for i := range m.Answer {
		if !reflect.DeepEqual(m.Answer[i], expectedRRs[i]) {
			t.Errorf("expected RR %d to be %v, got %v", i, expectedRRs[i], m.Answer[i])
		}
	}
}

func Test_addHTTPS(t *testing.T) {
	t.Parallel()

	m := new(dns.Msg)
	qName := "example.com."
	hintIPv4 := []net.IP{net.ParseIP("192.0.2.1")}
	hintIPv6 := []net.IP{net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")}

	addHTTPS(m, qName, hintIPv4, hintIPv6)

	expectedRR := &dns.SVCB{
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
			&dns.SVCBIPv4Hint{
				Hint: hintIPv4,
			},
			&dns.SVCBIPv6Hint{
				Hint: hintIPv6,
			},
		},
	}

	if len(m.Answer) != 1 {
		t.Fatalf("expected 1 RR, got %d", len(m.Answer))
	}

	if !reflect.DeepEqual(m.Answer[0], expectedRR) {
		t.Errorf("expected RR to be %v, got %v", expectedRR, m.Answer[0])
	}
}
