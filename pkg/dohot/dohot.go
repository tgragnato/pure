package dohot

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/miekg/dns"
)

var (
	proxy, _   = url.Parse("socks5://[::1]:9050")
	httpClient = &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   100 * time.Millisecond,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:      true,
		MaxIdleConns:           4,
		MaxIdleConnsPerHost:    4,
		MaxConnsPerHost:        4,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    10 * time.Second,
		ExpectContinueTimeout:  2 * time.Second,
		ResponseHeaderTimeout:  2 * time.Second,
		DisableKeepAlives:      false,
		DisableCompression:     true,
		MaxResponseHeaderBytes: 4096,
	}}
)

func DoH(qName string, ipv6 bool) ([]net.IP, []string, uint32, error) {
	var (
		ips    []net.IP
		cnames []string
		ttl    uint32
	)

	m := new(dns.Msg)
	if ipv6 {
		m.SetQuestion(qName, dns.TypeAAAA)
	} else {
		m.SetQuestion(qName, dns.TypeA)
	}
	m.SetEdns0(4096, true)
	out, err := m.Pack()
	if err != nil {
		return ips, cnames, ttl, errors.New("Error packing request")
	}

	req, err := http.NewRequest("POST", "https://dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion:443/dns-query", bytes.NewReader(out))
	if err != nil {
		return ips, cnames, ttl, errors.New("Invalid HTTP request")
	}
	req.Header.Set("content-type", "application/dns-message")

	resp, err := httpClient.Do(req)
	if err != nil {
		return ips, cnames, ttl, errors.New("Error doing HTTP request")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return ips, cnames, ttl, errors.New("Error reading response body")
	}

	rm := new(dns.Msg)
	if err := rm.Unpack(buf); err != nil {
		return ips, cnames, ttl, errors.New("Error unpacking response")
	}

	if rm.Rcode != dns.RcodeSuccess {
		return ips, cnames, ttl, errors.New("Error code in DNS response")
	}

	for _, ansa := range rm.Answer {
		switch ansb := ansa.(type) {
		case *dns.A:
			ips = append(ips, ansb.A)
			if ttl < ansb.Hdr.Ttl {
				ttl = ansb.Hdr.Ttl
			}
		case *dns.AAAA:
			ips = append(ips, ansb.AAAA)
			if ttl < ansb.Hdr.Ttl {
				ttl = ansb.Hdr.Ttl
			}
		case *dns.CNAME:
			cnames = append(cnames, ansb.Target)
		case *dns.DNAME:
			cnames = append(cnames, ansb.Target)
		}
	}

	if len(ips) == 0 {
		return ips, cnames, ttl, errors.New("No IP addresses in response")
	}

	return ips, cnames, ttl, nil
}
