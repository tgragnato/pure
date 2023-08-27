package main

import (
	"net"

	"gitlab.torproject.org/tpo/anti-censorship/geoip"
)

var (
	geo, errGeo  = geoip.New("/usr/share/tor/geoip", "/usr/share/tor/geoip6")
	unrouteables [28]*net.IPNet
)

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

		if geo != nil {
			country, ok := geo.GetCountryByAddr(ips[x])
			if !ok {
				return false
			}
			switch country {
			case "CN", "HK", "MO", "RU", "BY":
				return false
			}
		}
	}

	return true
}

func initCheck() {
	cidrStrings := []string{
		"0.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"255.255.255.255/32",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::/128",
		"::1/128",
		//"::ffff:0:0/96",
		//"::ffff:0:0:0/96",
		"64:ff9b::/96",
		"64:ff9b:1::/48",
		"100::/64",
		"2001:0000::/32",
		"2001:20::/28",
		"2001:db8::/32",
		"2002::/16",
		"fc00::/7",
		"fe80::/10",
		"ff00::/8",
	}

	for i := range cidrStrings {
		_, unrouteables[i], _ = net.ParseCIDR(cidrStrings[i])
	}
}
