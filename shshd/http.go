package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	proxy, _   = url.Parse("socks5://[::1]:9050")
	httpclient = &http.Client{
		Transport: &http.Transport{
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
			Proxy:                 http.ProxyURL(proxy),
		},
		Timeout: 5 * time.Minute,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {
	host := strings.TrimSuffix(r.Host, ":80")

	if host == "static.ess.apple.com" && r.URL.Path == "/connectivity.txt" {
		fmt.Fprint(w, "AV was here!")
		return
	}

	switch host {
	case "updates-http.cdn-apple.com", "gs.apple.com":
		r.URL.Scheme = "http"
	case "ocsp.digicert.com", "r3.o.lencr.org", "ocsp2.globalsign.com", "ocsp.sectigo.com", "ocsp.usertrust.com", "ocsp.godaddy.com", "ocsp.comodoca.com":
		r.URL.Scheme = "http"
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0")
		r.Header.Del("X-Apple-Request-UUID")
	case "static.ess.apple.com", "certs.apple.com":
		r.URL.Scheme = "https"
	case "ocsp.pki.goog", "ocsp.r2m01.amazontrust.com", "ocsp.rootca1.amazontrust.com":
		r.URL.Scheme = "https"
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/118.0")
		r.Header.Del("X-Apple-Request-UUID")
	default:
		http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusMovedPermanently)
		return
	}

	r.URL.Host = r.Host
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
