package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

var (
	unrouteables [26]*net.IPNet
	blackhole    [27]net.IP
	blacklist    []string
	whitelist    []string
	prefixes     []string
	dbreader     *maxminddb.Reader
	asnreader    *maxminddb.Reader
)

func initCheck() {
	cidrstrings := [26]string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"0.0.0.0/8",
		"100.64.0.0/10",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"ff00::/8",
		//"::ffff:0:0/96",
		//"::ffff:0:0:0/96",
		"64:ff9b::/96",
		"100::/64",
		"2001::/32",
		"2001:20::/28",
		"2001:db8::/32",
		"2002::/16",
		"139.45.192.0/18",
	}

	blackholestrings := [27]string{
		"34.217.236.88",
		"35.190.64.11",
		"35.190.74.49",
		"35.190.74.222",
		"35.201.98.64",
		"35.201.103.212",
		"54.200.164.214",
		"72.52.178.23",
		"74.117.179.8",
		"88.85.66.220",
		"109.206.162.83",
		"109.206.162.85",
		"109.206.169.172",
		"162.252.21.21",
		"173.214.252.142",
		"173.214.252.167",
		"188.42.218.242",
		"188.42.224.45",
		"188.42.224.57",
		"188.42.224.69",
		"192.243.59.12",
		"192.243.59.13",
		"192.243.59.20",
		"216.21.13.14",
		"216.21.13.15",
		"216.127.41.28",
		"216.172.60.116",
	}

	for i := range cidrstrings {
		_, unrouteables[i], _ = net.ParseCIDR(cidrstrings[i])
	}

	for i := range blackholestrings {
		blackhole[i] = net.ParseIP(blackholestrings[i])
	}

	blacklist = populateCheck("/etc/proxy/blocked.names")
	whitelist = populateCheck("/etc/proxy/allowed.names")
	prefixes = populateCheck("/etc/proxy/prefixes.names")
	dbreader, _ = maxminddb.Open("/etc/proxy/GeoLite2-Country.mmdb")
	asnreader, _ = maxminddb.Open("/etc/proxy/GeoLite2-ASN.mmdb")
}

func populateCheck(conf string) []string {
	var dNames []string

	buf, err := os.Open(conf)
	if err != nil {
		log.Printf("Error opening file %s", conf)
		return dNames
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Printf("Error closing file %s : %s", conf, err.Error())
		}
	}()

	snl := bufio.NewScanner(buf)
	for snl.Scan() {

		if err := snl.Err(); err == nil {
			txt := snl.Text()
			if !strings.HasPrefix(txt, "#") && txt != "" {
				if !strings.HasSuffix(txt, ".") {
					txt = txt + "."
				}
				dNames = append(dNames, txt)
			}
		} else {
			log.Printf("Error reading newline in file %s : %s", conf, err.Error())
			break
		}
	}

	return dNames
}

func checkQuery(domain string) bool {
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}
	for i := range whitelist {
		if strings.HasSuffix(domain, whitelist[i]) {
			return true
		}
	}
	for i := range prefixes {
		if strings.HasPrefix(domain, prefixes[i]) {
			return false
		}
	}
	for i := range blacklist {
		if strings.HasSuffix(domain, blacklist[i]) {
			return false
		}
	}
	return true
}

func checkDomain(domain string) bool {
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}
	if !checkQuery(domain) {
		return false
	}
	ip4, found4 := cache4.Get(domain)
	ip6, found6 := cache6.Get(domain)
	if !found4 && !found6 {
		return false
	}
	if found4 && !found6 {
		return checkIPs(ip4)
	}
	if !found4 && found6 {
		return checkIPs(ip6)
	}
	return checkIPs(ip4) && checkIPs(ip6)
}

func checkIPs(ips []net.IP) bool {
	for x := range ips {
		if !ips[x].IsGlobalUnicast() {
			return false
		}
		for y := range unrouteables {
			if unrouteables[y].Contains(ips[x]) {
				return false
			}
		}
		for y := range blackhole {
			if blackhole[y].Equal(ips[x]) {
				return false
			}
		}
		if dbreader != nil {
			var record struct {
				Country struct {
					ISOCode string `maxminddb:"iso_code"`
				} `maxminddb:"country"`
			}
			err := dbreader.Lookup(ips[x], &record)
			if err == nil {
				switch record.Country.ISOCode {
				case "CN", "HK", "MO":
					return false
				}
			}
		}
	}
	return true
}

func IP2ASN(ip net.IP) uint {
	if asnreader == nil {
		return 0
	}
	var record struct {
		AutonomousSystemNumber uint `maxminddb:"autonomous_system_number"`
	}
	err := asnreader.Lookup(ip, &record)
	if err == nil {
		return record.AutonomousSystemNumber
	}
	return 0
}

func CheckCF(domain string) bool {
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}
	ip4, found4 := cache4.Get(domain)
	ip6, found6 := cache6.Get(domain)
	if !found4 && !found6 {
		return false
	}
	ret := true
	if found4 {
		for i := range ip4 {
			ret = ret && IP2ASN(ip4[i]) == 13335
		}
	}
	if found6 {
		for i := range ip6 {
			ret = ret && IP2ASN(ip6[i]) == 13335
		}
	}
	return ret
}
