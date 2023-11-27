package sni

import (
	"net"
)

func checkDomain(domain string) bool {
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
	}

	return true
}
