package dohot

import (
	"bytes"
	"crypto/tls"
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
			Timeout:       100 * time.Millisecond,
			DualStack:     true,
			FallbackDelay: 100 * time.Millisecond,
			KeepAlive:     30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			ClientAuth:             0,
			InsecureSkipVerify:     false,
			SessionTicketsDisabled: true,
			ClientSessionCache:     nil,
			MinVersion:             tls.VersionTLS13,
			MaxVersion:             tls.VersionTLS13,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.X25519,
			},
			DynamicRecordSizingDisabled: false,
			Renegotiation:               0,
		},
		TLSHandshakeTimeout:    10 * time.Second,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxIdleConns:           4,
		MaxIdleConnsPerHost:    4,
		MaxConnsPerHost:        4,
		IdleConnTimeout:        90 * time.Second,
		ResponseHeaderTimeout:  2 * time.Second,
		ExpectContinueTimeout:  2 * time.Second,
		MaxResponseHeaderBytes: 4096,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      true,
	}}
)

func DoH(qName string, ipv6 bool) ([]net.IP, []string, error) {
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
		return ips, cnames, errors.New("error packing request")
	}

	req, err := http.NewRequest("POST", "https://dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion:443/dns-query", bytes.NewReader(out))
	if err != nil {
		return ips, cnames, errors.New("invalid HTTP request")
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := httpClient.Do(req)
	if err != nil {
		return ips, cnames, errors.New("error doing HTTP request")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return ips, cnames, errors.New("error reading response body")
	}

	rm := new(dns.Msg)
	if err := rm.Unpack(buf); err != nil {
		return ips, cnames, errors.New("error unpacking response")
	}

	if rm.Rcode != dns.RcodeSuccess {
		return ips, cnames, errors.New("error code in DNS response")
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
		return ips, cnames, errors.New("no IP addresses in response")
	}

	return ips, cnames, nil
}
