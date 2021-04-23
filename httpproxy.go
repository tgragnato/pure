package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var uastrings = [3]string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:66.0) Gecko/20100101 Firefox/66.0",
	"Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0",
	"Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0",
}

func redirectScheme(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 301)
}

func httpProxy(w http.ResponseWriter, r *http.Request, socks bool, ua bool) {
	if !checkDomain(r.Host) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 302)
		return
	}

	r.URL.Scheme = "http"
	r.URL.Host = r.Host

	if ua {
		index := rand.Int() % len(uastrings)
		r.Header.Set("User-Agent", uastrings[index])
	}

	var httpTransport *http.Transport
	if socks {
		httpTransport = http.DefaultTransport.(*http.Transport).Clone()
		httpTransport.Proxy = http.ProxyURL(proxyurl)
	} else {
		httpTransport = http.DefaultTransport.(*http.Transport).Clone()
	}

	resp, err := httpTransport.RoundTrip(r)
	if err != nil {
		http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 302)
		return
	}
	defer resp.Body.Close()

	respH := w.Header()
	for hk := range resp.Header {
		respH[hk] = resp.Header[hk]
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {

	go IncHTTP(r.Host)

	if strings.HasPrefix(r.RemoteAddr, "172.16.31.0:") {
		httpProxy(w, r, false, true)
	} else if strings.HasSuffix(r.Host, ".apple.com") && r.Host != "ocsp.apple.com" {
		httpProxy(w, r, false, false)
	} else if r.Host == "updates-http.cdn-apple.com" {
		httpProxy(w, r, true, false)
	} else {
		redirectScheme(w, r)
	}
}
