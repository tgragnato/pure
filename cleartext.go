package main

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

var dnsClient = &dns.Client{
	Dialer: &net.Dialer{
		Timeout: 500 * time.Millisecond,
		LocalAddr: &net.UDPAddr{
			IP:   net.ParseIP("192.168.1.31"),
			Port: 35353,
			Zone: "",
		},
	},
	SingleInflight: true,
}

func InitCleartext() {
	ticker := time.Tick(time.Hour / 100000)
	i := 0
	for {
		select {
		case <-ticker:
			for _, cc := range [3]string{"us", "eu", "asia"} {
				for _, ud := range [2]string{"upload", "download"} {
					qname := strconv.FormatInt(int64(i), 10)
					for t := 5 - len(qname); t > 0; t-- {
						qname = "0" + qname
					}
					qname = "gcs-" + cc + "-" + qname + ".content-storage-" + ud + ".googleapis.com."
					go func() {
						ip4, err4 := Cleartext(qname, false)
						if err4 == nil {
							cache4.Set(qname, ip4)
						}
					}()
					go func() {
						ip6, err6 := Cleartext(qname, true)
						if err6 == nil {
							cache6.Set(qname, ip6)
						}
					}()
				}
			}
			i++
			if i >= 100000 {
				i = 0
			}
		}
	}
}

func Cleartext(qName string, ipv6 bool) ([]net.IP, error) {
	ips := []net.IP{}

	m := new(dns.Msg)
	if ipv6 {
		m.SetQuestion(qName, dns.TypeAAAA)
	} else {
		m.SetQuestion(qName, dns.TypeA)
	}
	m.SetEdns0(4096, true)

	rm, _, err := dnsClient.Exchange(m, "192.168.1.254:53")
	if err != nil {
		return ips, errors.New("Error during DNS exchange")
	}

	if rm.Rcode != dns.RcodeSuccess {
		return ips, errors.New("Error code in DNS response")
	}

	for _, ansa := range rm.Answer {
		switch ansb := ansa.(type) {
		case *dns.A:
			ips = append(ips, ansb.A)
		case *dns.AAAA:
			ips = append(ips, ansb.AAAA)
		}
	}

	if len(ips) == 0 {
		return ips, errors.New("No IP addresses in response")
	}

	return ips, nil
}
