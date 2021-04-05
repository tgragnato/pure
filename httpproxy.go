package main

import (
	"io"
	"net/http"
	"strings"
)

func redirectScheme(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 301)
}

func handleDirPort(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.RawPath, "/tor") {
		original := r.URL.Host
		r.URL.Scheme = "http"
		r.URL.Host = "127.0.0.1:9030"

		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Redirect(w, r, "https://"+original+r.URL.RequestURI(), 302)
			return
		}
		defer resp.Body.Close()

		respH := w.Header()
		for hk := range resp.Header {
			respH[hk] = resp.Header[hk]
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	} else {
		redirectScheme(w, r)
	}
}

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {

	if strings.HasSuffix(r.Host, ".apple.com") && r.Host != "ocsp.apple.com" {

		if !checkDomain(r.Host) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 302)
			return
		}

		r.URL.Scheme = "http"
		r.URL.Host = r.Host

		resp, err := http.DefaultTransport.RoundTrip(r)
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

	} else if r.Host == "updates-http.cdn-apple.com" {

		httpTransport := http.DefaultTransport.(*http.Transport).Clone()
		httpTransport.Proxy = http.ProxyURL(proxyurl)

		r.URL.Scheme = "http"
		r.URL.Host = r.Host

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

	} else {
		redirectScheme(w, r)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RemoteAddr, "172.16.31.") {
		go IncHTTP(r.Host)
		handleHTTPForward(w, r)
	} else {
		handleDirPort(w, r)
	}
}
