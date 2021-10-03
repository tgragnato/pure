package main

import (
	"net"
	"strings"
)

var (
	cloudflare []*net.IPNet
	allowednet []*net.IPNet
)

func Cloaking(qName string) string {
	if !strings.HasSuffix(qName, ".") {
		qName = qName + "."
	}
	qName = strings.ToLower(qName)

	switch qName {
	case "www.youtube.com.":
		qName = "restrictmoderate.youtube.com."
	case "m.youtube.com.":
		qName = "restrictmoderate.youtube.com."
	case "www.youtube-nocookie.com.":
		qName = "restrictmoderate.youtube.com."
	}

	if strings.HasSuffix(qName, "tgragnato.it.") {
		qName = "tgragnato.it."
	}

	if strings.HasSuffix(qName, "github.io.") {
		qName = "github.io."
	}

	return qName
}

func InitCloaking() {
	cloudflarestrings := []string{
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}

	allowednetstrings := []string{
		"192.168.194.0/31",
		"127.0.0.1/8",
	}

	for i := range cloudflarestrings {
		_, cidr, _ := net.ParseCIDR(cloudflarestrings[i])
		cloudflare = append(cloudflare, cidr)
	}

	for i := range allowednetstrings {
		_, cidr, _ := net.ParseCIDR(allowednetstrings[i])
		allowednet = append(allowednet, cidr)
	}
}

func IsAllowed(ip net.IP) bool {
	for i := range cloudflare {
		if cloudflare[i].Contains(ip) {
			return true
		}
	}
	for i := range allowednet {
		if allowednet[i].Contains(ip) {
			return true
		}
	}
	return false
}
