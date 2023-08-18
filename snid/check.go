package main

import (
	"net"

	"gitlab.torproject.org/tpo/anti-censorship/geoip"
)

var geo, errGeo = geoip.New("/usr/share/tor/geoip", "/usr/share/tor/geoip6")

func CheckDomain(domain string) bool {
	if net.ParseIP(domain) != nil {
		return false
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for x := range ips {
		if !ips[x].IsGlobalUnicast() {
			return false
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
