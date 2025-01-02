package dnshandlers

import (
	"net"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miekg/dns"
)

func TestDnsHandlers_ParseQuery(t *testing.T) {
	t.Parallel()

	const domain = "example.com."

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	d := &DnsHandlers{
		db:        mockDb,
		geoChecks: nil,
		hintIPv4:  net.ParseIP("192.0.2.3").To4(),
		hintIPv6:  net.ParseIP("2001:db8::3"),
	}

	rowsA := sqlmock.NewRows([]string{"value"})
	rowsA.AddRow("192.0.2.1,203.0.113.1,")

	rowsAAAA := sqlmock.NewRows([]string{"value"})
	rowsAAAA.AddRow("2001:db8::1,2001:db8::2,")

	t.Run("Test query A", func(t *testing.T) {
		m := &dns.Msg{
			Question: []dns.Question{
				{
					Name:  domain,
					Qtype: dns.TypeA,
				},
			},
		}

		mock.ExpectQuery("SELECT value FROM a").
			WithArgs(domain).
			WillReturnRows(rowsA)

		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Fatalf("Expected 2 answers, got %d", len(m.Answer))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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

		mock.ExpectQuery("SELECT value FROM aaaa").
			WithArgs(domain).
			WillReturnRows(rowsAAAA)

		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Errorf("Expected 2 answers, got %d", len(m.Answer))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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

		mock.ExpectQuery("SELECT value FROM a").
			WithArgs("testprohibited.tgragnato.it.").
			WillReturnRows(sqlmock.NewRows([]string{"value"}))

		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Errorf("Expected 2 answers, got %d", len(m.Answer))
		}
		if m.Answer[0].Header().Rrtype != dns.TypeA {
			t.Errorf("Expected A record, got %d", m.Answer[0].Header().Rrtype)
		}
		if m.Answer[0].String() != "testprohibited.tgragnato.it.	3600	IN	A	0.0.0.0" {
			t.Errorf("Expected null record, got %s", m.Answer[0].String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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

		mock.ExpectQuery("SELECT value FROM aaaa").
			WithArgs("testprohibited.tgragnato.it.").
			WillReturnRows(sqlmock.NewRows([]string{"value"}))

		d.ParseQuery(m)
		if len(m.Answer) != 2 {
			t.Errorf("Expected 2 answers, got %d", len(m.Answer))
		}
		if m.Answer[1].Header().Rrtype != dns.TypeAAAA {
			t.Errorf("Expected A record, got %d", m.Answer[0].Header().Rrtype)
		}
		if m.Answer[1].String() != "testprohibited.tgragnato.it.	3600	IN	AAAA	::" {
			t.Errorf("Expected null record, got %s", m.Answer[1].String())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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

		mock.ExpectQuery("SELECT value FROM a").
			WithArgs(domain).
			WillReturnRows(rowsA)
		mock.ExpectQuery("SELECT value FROM aaaa").
			WithArgs(domain).
			WillReturnRows(rowsAAAA)

		d.ParseQuery(m)
		if len(m.Answer) != 1 {
			t.Errorf("Expected 1 answer, got %d", len(m.Answer))
		}
		if m.Answer[0].Header().Rrtype != dns.TypeHTTPS {
			t.Errorf("Expected HTTPS record, got %d", m.Answer[0].Header().Rrtype)
		}

		/*
			if m.Answer[0].String() != "example.com.	86400	IN	HTTPS	1 . alpn=\"h3,h2\" mandatory=\"alpn\" ipv4hint=\"192.0.2.1,203.0.113.1\" ipv6hint=\"2001:db8::1,2001:db8::2\"" {
				t.Errorf("Expected HTTPS record, got %s", m.Answer[0].String())
			}
		*/
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
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
			t.Errorf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeNotImplemented {
			t.Errorf("Expected Not Implemented, got %d", m.Rcode)
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
			t.Errorf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeRefused {
			t.Errorf("Expected Refused, got %d", m.Rcode)
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
			t.Errorf("Expected no answers, got %d", len(m.Answer))
		}
		if m.Rcode != dns.RcodeRefused {
			t.Errorf("Expected Refused, got %d", m.Rcode)
		}
	})
}
