package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/miekg/dns"
)

func DoH(qName string, trr string) ([]net.IP, []string, error) {
	var (
		ips    []net.IP
		cnames []string
	)

	m := new(dns.Msg)
	m.SetQuestion(qName, dns.TypeA)
	m.SetEdns0(4096, true)
	out, err := m.Pack()
	if err != nil {
		return ips, cnames, errors.New("Error packing request")
	}

	req, err := http.NewRequest("POST", "https://"+trr+":443/dns-query", bytes.NewReader(out))
	if err != nil {
		return ips, cnames, errors.New("Invalid HTTP request")
	}
	req.Header.Set("content-type", "application/dns-message")

	httpClient := &http.Client{Transport: http.DefaultTransport}
	if strings.HasSuffix(trr, ".onion") {
		proxyurl, err := url.Parse("socks5://127.0.0.1:9050")
		if err != nil {
			return ips, cnames, errors.New("Invalid socks5 proxy")
		}
		httpTransport := http.DefaultTransport.(*http.Transport).Clone()
		httpTransport.Proxy = http.ProxyURL(proxyurl)
		httpClient = &http.Client{Transport: httpTransport}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return ips, cnames, errors.New("Error doing HTTP request")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ips, cnames, errors.New("Error reading response body")
	}

	rm := new(dns.Msg)
	if err := rm.Unpack(buf); err != nil {
		return ips, cnames, errors.New("Error unpacking response")
	}

	if rm.Rcode != dns.RcodeSuccess {
		return ips, cnames, errors.New("Error code in DNS response")
	}

	for _, ansa := range rm.Answer {
		switch ansb := ansa.(type) {
		case *dns.A:
			ips = append(ips, ansb.A)
		case *dns.CNAME:
			cnames = append(cnames, ansb.Target)
		case *dns.DNAME:
			cnames = append(cnames, ansb.Target)
		}
	}

	if len(ips) == 0 {
		return ips, cnames, errors.New("No IP addresses in response")
	}

	return ips, cnames, nil
}

func retNull(m *dns.Msg, qName string) {
	newRR, err := dns.NewRR(fmt.Sprintf("%s A %s", qName, "0.0.0.0"))
	if err == nil {
		m.Answer = nil
		m.Answer = append(m.Answer, newRR)
	}
}

func addIP(m *dns.Msg, qName string, ip net.IP) {
	rr, err := dns.NewRR(fmt.Sprintf("%s A %s", qName, ip.String()))
	if err == nil {
		m.Answer = append(m.Answer, rr)
	}
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			ips, _, err := DoH(q.Name, "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion")
			if err != nil {
				retNull(m, q.Name)
			} else {
				for _, ip := range ips {
					if !checkIP(ip) {
						retNull(m, q.Name)
						return
					} else {
						addIP(m, q.Name, ip)
					}
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
