package dnshandlers

import (
	"net"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func ipSlicesEqual(a, b []net.IP) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func TestSetPersistent(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name  string
		key   string
		ips   []net.IP
		ipv6  bool
		table string
	}{
		{
			name:  "IPv4 insertion",
			key:   "example.com",
			ips:   []net.IP{net.ParseIP("192.0.2.1"), net.ParseIP("192.0.2.2")},
			ipv6:  false,
			table: "a",
		},
		{
			name:  "IPv6 insertion",
			key:   "example.org",
			ips:   []net.IP{net.ParseIP("2001:db8::1"), net.ParseIP("2001:db8::2")},
			ipv6:  true,
			table: "aaaa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := ""
			for _, ip := range tt.ips {
				expected += ip.String() + ","
			}

			// Expect the INSERT/UPDATE query
			mock.ExpectExec("INSERT INTO "+tt.table).
				WithArgs(tt.key, expected).
				WillReturnResult(sqlmock.NewResult(1, 1))

			d := &DnsHandlers{db: db}
			d.setPersistent(tt.key, tt.ips, tt.ipv6)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetPersistent(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		key      string
		ipv6     bool
		table    string
		value    string
		expected []net.IP
		found    bool
	}{
		{
			name:     "IPv4 lookup found",
			key:      "example.com",
			ipv6:     false,
			table:    "a",
			value:    "192.0.2.1,192.0.2.2,",
			expected: []net.IP{net.ParseIP("192.0.2.1"), net.ParseIP("192.0.2.2")},
			found:    true,
		},
		{
			name:     "IPv6 lookup found",
			key:      "example.org",
			ipv6:     true,
			table:    "aaaa",
			value:    "2001:db8::1,2001:db8::2,",
			expected: []net.IP{net.ParseIP("2001:db8::1"), net.ParseIP("2001:db8::2")},
			found:    true,
		},
		{
			name:     "Record not found",
			key:      "notfound.com",
			ipv6:     false,
			table:    "a",
			value:    "",
			expected: []net.IP{},
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"value"})
			if tt.found {
				rows.AddRow(tt.value)
			}

			mock.ExpectQuery("SELECT value FROM " + tt.table).
				WithArgs(tt.key).
				WillReturnRows(rows)

			d := &DnsHandlers{db: db}
			ips, found := d.getPersistent(tt.key, tt.ipv6)

			if found != tt.found {
				t.Errorf("getPersistent() found = %v, want %v", found, tt.found)
			}

			if !ipSlicesEqual(ips, tt.expected) {
				t.Errorf("getPersistent() ips = %v, want %v", ips, tt.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCleanPersistent(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		domains map[string][]string // table -> domains
	}{
		{
			name: "Clean invalid domains",
			domains: map[string][]string{
				"a":    {"example.com", "INVALID.COM", "test.com"},
				"aaaa": {"example.org", "INVALID.ORG", "test.org"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectations for A records
			rowsA := sqlmock.NewRows([]string{"key"})
			for _, domain := range tt.domains["a"] {
				rowsA.AddRow(domain)
			}
			mock.ExpectQuery("SELECT key FROM a").WillReturnRows(rowsA)

			// Expect DELETE queries for invalid A records
			mock.ExpectExec("DELETE FROM a").WithArgs("INVALID.COM").
				WillReturnResult(sqlmock.NewResult(1, 1))

			// Set up expectations for AAAA records
			rowsAAAA := sqlmock.NewRows([]string{"key"})
			for _, domain := range tt.domains["aaaa"] {
				rowsAAAA.AddRow(domain)
			}
			mock.ExpectQuery("SELECT key FROM aaaa").WillReturnRows(rowsAAAA)

			// Expect DELETE queries for invalid AAAA records
			mock.ExpectExec("DELETE FROM aaaa").WithArgs("INVALID.ORG").
				WillReturnResult(sqlmock.NewResult(1, 1))

			d := &DnsHandlers{db: db}
			d.cleanPersistent()

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
