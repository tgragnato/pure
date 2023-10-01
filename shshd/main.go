package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
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

func main() {

	handler := http.DefaultServeMux
	handler.HandleFunc("/", handleHTTPForward)

	err := http.ListenAndServe("[::1]:80", handler)
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}
