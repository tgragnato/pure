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

	if strings.HasSuffix(r.Host, ".apple.com") && r.Host != "ocsp.apple.com" || r.Host == "updates-http.cdn-apple.com" {

		if !checkDomain(r.Host) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), 302)
			return
		}

		/*proxyurl, err := url.Parse("socks5://127.0.0.1:9050")
		  if err != nil {
		          http.Error(w, "Could not parse proxy URL", 500)
		          return
		}*/
		httpTransport := http.DefaultTransport.(*http.Transport).Clone()
		//httpTransport.Proxy = http.ProxyURL(proxyurl)

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
		handleHTTPForward(w, r)
	} else {
		handleDirPort(w, r)
	}
}
