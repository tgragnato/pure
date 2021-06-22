package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	httpSocks = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyurl),
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Minute,
				KeepAlive: time.Millisecond,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       20,
			IdleConnTimeout:       5 * time.Minute,
			ResponseHeaderTimeout: 2 * time.Second,
		},
		Timeout: 5 * time.Minute,
	}
)

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {
	host := strings.TrimSuffix(r.Host, ":80")
	go analytics.IncHTTP(host)

	if host == "static.ess.apple.com" && r.URL.Path == "/connectivity.txt" {
		fmt.Fprint(w, "AV was here!")
		return
	}

	is_apple := strings.HasSuffix(host, ".apple.com") && host != "ocsp.apple.com"
	is_updates := r.Host == "updates-http.cdn-apple.com"

	if !is_apple && !is_updates {
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 301)
		return
	}

	if !checkDomain(host) {
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 302)
		return
	}

	r.URL.Scheme = "http"
	r.URL.Host = r.Host
	r.RequestURI = ""

	resp, err := httpSocks.Do(r)
	if err != nil {
		log.Printf("Error doing HTTP request http://%s%s", host, r.URL.RequestURI())
		log.Printf("   Printing error: %s", err.Error())
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 302)
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
			case <-time.Tick(10 * time.Millisecond):
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
