package dohot

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

func Transparent(qName string, t uint16, apple bool) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(qName, t)
	m.SetEdns0(4096, true)

	out, err := m.Pack()
	if err != nil {
		return nil, errors.New("error packing request")
	}

	req, err := http.NewRequest("POST", "https://dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion:443/dns-query", bytes.NewReader(out))
	if apple {
		req, err = http.NewRequest("POST", "https://doh.dns.apple.com:443/dns-query", bytes.NewReader(out))
	}
	if err != nil {
		return nil, errors.New("invalid HTTP request")
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New("error doing HTTP request")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	rm := new(dns.Msg)
	if err := rm.Unpack(buf); err != nil {
		return nil, errors.New("error unpacking response")
	}

	if rm.Rcode != dns.RcodeSuccess {
		return nil, errors.New("error code in DNS response")
	}

	return rm.Answer, nil
}
