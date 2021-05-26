package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	httpSocks = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyurl),
		},
		Timeout: 120 * time.Second,
	}
	httpDirect = &http.Client{
		Timeout: 60 * time.Second,
	}
	uastrings = [2]string{
		"Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0",
	}
)

func httpProxy(w http.ResponseWriter, r *http.Request, socks bool, ua bool) {
	host := strings.TrimSuffix(r.Host, ":80")
	if !checkDomain(host) {
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 302)
		return
	}

	r.URL.Scheme = "http"
	r.URL.Host = r.Host
	r.RequestURI = ""

	if ua {
		index := rand.Int() % len(uastrings)
		r.Header.Set("User-Agent", uastrings[index])
	}

	var htc *http.Client
	if socks {
		htc = httpSocks
	} else {
		htc = httpDirect
	}

	resp, err := htc.Do(r)
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
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Millisecond):
				w.(http.Flusher).Flush()
			case <-done:
				return
			}
		}
	}()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	close(done)
}

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {

	host := strings.TrimSuffix(r.Host, ":80")
	if host == "static.ess.apple.com" && r.URL.Path == "/connectivity.txt" {
		fmt.Fprint(w, "AV was here!")
		return
	}
	go IncHTTP(host)

	if strings.HasPrefix(r.RemoteAddr, "172.16.31.0:") {
		httpProxy(w, r, false, true)
	} else if strings.HasPrefix(r.RemoteAddr, "[fd76:abcd:ef90::]:") {
		httpProxy(w, r, false, true)
	} else if strings.HasSuffix(host, ".apple.com") && host != "ocsp.apple.com" {
		httpProxy(w, r, false, false)
	} else if r.Host == "updates-http.cdn-apple.com" {
		httpProxy(w, r, true, false)
	} else {
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 301)
	}
}
