package main

import (
	"net"
	"strings"
)

var unrouteables [15]*net.IPNet

func initCheck() {
	cidrstrings := [15]string{
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
	}

	for i := range cidrstrings {
		_, unrouteables[i], _ = net.ParseCIDR(cidrstrings[i])
	}
}

func checkDomain(domain string) bool {
	ip, err := net.LookupIP(domain)
	if err != nil {
		return false
	}
	for x := range ip {
		if !ip[x].IsGlobalUnicast() {
			return false
		}
		for y := range unrouteables {
			if unrouteables[y].Contains(ip[x]) {
				return false
			}
		}
	}
	return true
}

func backend(local bool, sni string) string {
	if local {
		if strings.HasSuffix(sni, "tgragnato.it") && sni != "status.tgragnato.it" {
			return "127.0.0.1:8080"
		} else {
			if checkDomain(sni) {
				return net.JoinHostPort(sni, "443")
			} else {
				return ""
			}
		}
	} else {
		if strings.HasSuffix(sni, "tgragnato.it") {
			return "127.0.0.1:8080"
		} else {
			return "127.0.0.1:9001"
		}
	}
}
