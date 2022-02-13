package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func CheckDomain(domain string) bool {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for x := range ips {
		if !ips[x].IsGlobalUnicast() {
			return false
		}

		if dbreader != nil {
			var record struct {
				Country struct {
					ISOCode string `maxminddb:"iso_code"`
				} `maxminddb:"country"`
			}
			err := dbreader.Lookup(ips[x], &record)
			if err == nil {
				switch record.Country.ISOCode {
				case "CN", "HK", "MO":
					return false
				}
			}
		}
	}

	return true
}

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {
	host := strings.TrimSuffix(r.Host, ":80")
	if !CheckDomain(host) {
		return
	}
	go analytics.IncHTTP(host)

	r.URL.Scheme = "http"
	r.URL.Host = host
	r.RequestURI = ""

	resp, err := httpclient.Do(r)
	if err != nil {
		log.Printf("Error doing HTTP request http://%s%s", host, r.URL.RequestURI())
		log.Printf("   Printing error: %s", err.Error())
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusFound)
		return
	}
	defer resp.Body.Close()

	respH := w.Header()
	for hk := range resp.Header {
		respH[hk] = resp.Header[hk]
	}

	done := make(chan bool)
	defer close(done)

	go func() {
		for {
			select {
			case <-time.Tick(time.Second / 3):
				f, ok := w.(http.Flusher)
				if ok {
					f.Flush()
				}
			case <-done:
				return
			}
		}
	}()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
