package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

var (
	unrouteables [25]*net.IPNet
	blackhole    []net.IP
	blacklist    []string
	whitelist    []string
	prefixes     []string
)

func initCheck() {
	cidrstrings := [25]string{
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
	}

	blackholestrings := []string{
		"35.190.64.11",
		"35.190.74.222",
		"35.201.98.64",
		"35.201.103.212",
		"74.117.179.8",
		"88.85.66.220",
		"139.45.196.140",
		"139.45.196.145",
		"139.45.196.200",
		"139.45.196.206",
		"139.45.197.14",
		"139.45.197.64",
		"139.45.197.66",
		"139.45.197.67",
		"139.45.197.81",
		"139.45.197.82",
		"139.45.197.90",
		"139.45.197.105",
		"139.45.197.108",
		"139.45.197.116",
		"139.45.197.178",
		"139.45.197.235",
		"139.45.197.236",
		"139.45.197.237",
		"139.45.197.238",
		"139.45.197.239",
		"139.45.197.243",
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
		blackhole = append(blackhole, net.ParseIP(blackholestrings[i]))
	}

	blacklist = populateCheck("/etc/proxy/blocked.names")
	whitelist = populateCheck("/etc/proxy/allowed.names")
	prefixes = populateCheck("/etc/proxy/prefixes.names")

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
	if !checkQuery(domain) {
		return false
	}
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}
	return checkIPs(ips)
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
	}
	return true
}
