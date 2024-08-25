package dnshandlers

import (
	"database/sql"
	"net"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newDb(t *testing.T) *sql.DB {
	dbConn, err := sql.Open("sqlite3", "file:"+t.Name()+"?cache=shared&mode=memory")
	if err != nil {
		return nil
	}

	_, err = dbConn.Exec(`
		CREATE TABLE a (
			key TEXT PRIMARY KEY,
			value TEXT,
			discovered_on TIMESTAMP,
			last_used TIMESTAMP
		);
	`)
	if err != nil {
		return nil
	}

	_, err = dbConn.Exec(`
		CREATE TABLE aaaa (
			key TEXT PRIMARY KEY,
			value TEXT,
			discovered_on TIMESTAMP,
			last_used TIMESTAMP
		);
	`)
	if err != nil {
		return nil
	}

	return dbConn
}

func TestSetGetPersistent_IPv4(t *testing.T) {
	t.Parallel()

	d := &DnsHandlers{db: newDb(t)}
	key := "example.com"
	data := []net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("203.0.113.1"),
	}
	d.setPersistent(key, data, false)

	ips, found := d.getPersistent(key, false)
	if !found {
		t.Errorf("Expected to find persistent data for key %s, but not found", key)
	}
	if !reflect.DeepEqual(ips, data) {
		t.Errorf("For key %s, expected %v, but got %v", key, data, ips)
	}
}

func TestSetGetPersistent_IPv6(t *testing.T) {
	t.Parallel()

	d := &DnsHandlers{db: newDb(t)}
	key := "example.net"
	data := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
	}
	d.setPersistent(key, data, true)

	ips, found := d.getPersistent(key, true)
	if !found {
		t.Errorf("Expected to find persistent data for key %s, but not found", key)
	}
	if !reflect.DeepEqual(ips, data) {
		t.Errorf("For key %s, expected %v, but got %v", key, data, ips)
	}
}
