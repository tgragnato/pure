package http

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"time"
)

type responseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w responseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Accept-Encoding")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		gzr := responseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzr, r)
	})
}

func flateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Accept-Encoding")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "deflate")
		fl, err := flate.NewWriter(w, flate.BestCompression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fl.Close()
		flr := responseWriter{Writer: fl, ResponseWriter: w}
		next.ServeHTTP(flr, r)
	})
}

func compressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodings := strings.Split(r.Header.Get("Accept-Encoding"), ",")
		for index := range encodings {
			encodings[index] = strings.TrimSpace(encodings[index])
			encodings[index] = strings.ToLower(encodings[index])
		}
		if i := sort.SearchStrings(encodings, "gzip"); i < len(encodings) && encodings[i] == "gzip" {
			gzipMiddleware(next).ServeHTTP(w, r)
			return
		}
		if i := sort.SearchStrings(encodings, "deflate"); i < len(encodings) && encodings[i] == "deflate" {
			flateMiddleware(next).ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expect-Ct", "max-age=86400, enforce")
		w.Header().Set("Permissions-Policy", "interest-cohort=(), accelerometer=(), autoplay=(), camera=(), clipboard-read=(), clipboard-write=(), document-domain=(), encrypted-media=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=(), gamepad=(), vibrate=(), vr=(), xr-spatial-tracking=()")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Server", "Apache/2.4.25 (RedStar4.0) OpenSSL/1.0.1e-fips")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", "https://tgragnato.it")
		}

		if w.Header().Get("Content-Security-Policy") == "" {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; object-src 'none'; upgrade-insecure-requests;")
		}

		if w.Header().Get("Set-Cookie") == "" {
			characters := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
			selected := make([]rune, 26)
			for i := range selected {
				selected[i] = characters[rand.IntN(len(characters))]
			}
			w.Header().Set("Set-Cookie", "PHPSESSID="+string(selected)+"; path=/")
		}

		next.ServeHTTP(w, r)
	})
}

func stripProxiedHeaders(resp *http.Response) error {
	resp.Header.Del("Expect-Ct")
	resp.Header.Del("Permissions-Policy")
	resp.Header.Del("Referrer-Policy")
	resp.Header.Del("Server")
	resp.Header.Del("Strict-Transport-Security")
	resp.Header.Del("X-Content-Type-Options")
	resp.Header.Del("X-Frame-Options")
	resp.Header.Del("X-XSS-Protection")
	return nil
}

func apiGateway() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err != nil || !net.ParseIP(host).IsPrivate() {
			http.Redirect(w, r, "https://tgragnato.it"+r.URL.RequestURI(), http.StatusFound)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/transmission") {
			if strings.HasPrefix(r.URL.Path, "/transmission/web/") || r.URL.Path == "/transmission/rpc" {
				proxy := httputil.NewSingleHostReverseProxy(&url.URL{
					Scheme: "http",
					Host:   "127.0.0.1:9091",
				})
				proxy.Transport = &http.Transport{
					ExpectContinueTimeout: time.Second,
					ForceAttemptHTTP2:     false,
					IdleConnTimeout:       time.Minute,
					ResponseHeaderTimeout: time.Second,
				}
				proxy.ModifyResponse = stripProxiedHeaders
				proxy.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "https://api.tgragnato.it/transmission/web/", http.StatusFound)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/grafana") {
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "https",
				Host:   "[::1]:3000",
			})
			proxy.Transport = &http.Transport{
				ExpectContinueTimeout: time.Second,
				ForceAttemptHTTP2:     true,
				IdleConnTimeout:       time.Minute,
				ResponseHeaderTimeout: time.Second,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				TLSHandshakeTimeout:   time.Second,
			}
			proxy.ModifyResponse = stripProxiedHeaders
			proxy.ServeHTTP(w, r)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   "[::1]:8080",
		})
		proxy.Transport = &http.Transport{
			ExpectContinueTimeout: time.Second,
			ForceAttemptHTTP2:     false,
			IdleConnTimeout:       time.Minute,
			ResponseHeaderTimeout: time.Second,
		}
		proxy.ModifyResponse = stripProxiedHeaders
		proxy.ServeHTTP(w, r)
	})
}
