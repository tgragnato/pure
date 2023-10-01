package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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
	case "ocsp.digicert.com", "r3.o.lencr.org", "ocsp2.globalsign.com", "ocsp.sectigo.com", "ocsp.usertrust.com", "ocsp.godaddy.com":
		r.URL.Scheme = "http"
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0")
		r.Header.Del("X-Apple-Request-UUID")
	case "static.ess.apple.com", "certs.apple.com":
		r.URL.Scheme = "https"
	case "ocsp.pki.goog", "ocsp.r2m01.amazontrust.com", "ocsp.rootca1.amazontrust.com":
		r.URL.Scheme = "https"
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/116.0")
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
