package sni

import (
	"net"
	"strings"
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

func getHostPort(sni string) string {
	if strings.HasSuffix(sni, "imap.mail.me.com") || sni == "imap.gmail.com" {
		return net.JoinHostPort(sni, "993")
	}

	return net.JoinHostPort(sni, "443")
}
