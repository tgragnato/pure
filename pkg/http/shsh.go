package http

import (
	"fmt"
	"io"
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

func handleSHSHProtocol() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.TrimSuffix(r.Host, ":80")

		if host == "static.ess.apple.com" && r.URL.Path == "/connectivity.txt" {
			fmt.Fprint(w, "AV was here!")
			return
		}

		if strings.HasPrefix(host, ".tgragnato.it") {
			http.FileServer(http.Dir("/var/www")).ServeHTTP(w, r)
			return
		}

		if host != "updates-http.cdn-apple.com" &&
			host != "gs.apple.com" &&
			host != "static.ess.apple.com" &&
			host != "certs.apple.com" {
			http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusMovedPermanently)
			return
		}

		r.URL.Scheme = "http"
		r.URL.Host = r.Host
		r.RequestURI = ""

		resp, err := httpclient.Do(r)
		if err != nil {
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
				case <-time.NewTicker(time.Second / 3).C:
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
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
