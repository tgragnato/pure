package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
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
	chars = [64]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ",", "-",
	}
)

func handleHTTPForward(w http.ResponseWriter, r *http.Request) {
	host := strings.TrimSuffix(r.Host, ":80")
	r.URL.Scheme = "http"
	r.RequestURI = ""
	var httpclient *http.Client

	ipstring, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	ip := net.ParseIP(ipstring)
	if ip == nil {
		log.Print("Error: ip == nil")
		return
	}

	if !private4.Contains(ip) && !private6.Contains(ip) {

		if !strings.HasPrefix(r.URL.RequestURI(), "/tor/") {
			es := ""
			for i := 0; i < 26; i++ {
				es += chars[rand.Int()%len(chars)]
			}
			w.Header().Set("Server", "Apache/2.2.15 (RedStar4.0)")
			w.Header().Set("X-Powered-By", "PHP/5.3.5")
			http.SetCookie(w, &http.Cookie{
				Name:  "PHPSESSID",
				Value: es,
				Path:  "/",
			})
			w.Header().Set("Cache-Control", "max-age=1, private, must-revalidate")
			w.Header().Set("Expires", "Thu, 19 Nov 1981 08:52:00 GMT")
			w.Header().Set("Pragma", "no-cache")
			http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), 301)
			return
		}

		r.URL.Host = "127.0.0.1:9030"
		httpclient = http.DefaultClient

	} else {
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
		r.URL.Host = r.Host
		httpclient = httpSocks
	}

	resp, err := httpclient.Do(r)
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
