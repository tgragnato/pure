package http

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
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
	status int
	size   int64
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (size int, err error) {
	if w.Writer == nil {
		size, err = w.ResponseWriter.Write(b)
	} else {
		size, err = w.Writer.Write(b)
	}
	w.size += int64(size)
	return size, err
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not support Hijack")
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
		gzr := &responseWriter{Writer: gz, ResponseWriter: w}
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
		flr := &responseWriter{Writer: fl, ResponseWriter: w}
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
		w.Header().Set("Access-Control-Allow-Origin", "https://tgragnato.it")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' blob:; img-src 'self' data:; worker-src 'self' blob:; object-src 'none'; upgrade-insecure-requests;")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		w.Header().Set("Expect-Ct", "max-age=86400, enforce")
		w.Header().Set("Expires", "Thu, 19 Nov 1981 08:52:00 GMT")
		w.Header().Set("Permissions-Policy", "interest-cohort=(), accelerometer=(), autoplay=(), camera=(), clipboard-read=(), clipboard-write=(), document-domain=(), encrypted-media=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=(), gamepad=(), vibrate=(), vr=(), xr-spatial-tracking=()")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Server", "Apache/2.4.25 (RedStar4.0) OpenSSL/1.0.1e-fips PHP/5.6.2")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Powered-By", "PHP/5.6.2")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

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
	resp.Header.Del("Access-Control-Allow-Origin")
	resp.Header.Del("Cache-Control")
	resp.Header.Del("Content-Security-Policy")
	resp.Header.Del("Cross-Origin-Embedder-Policy")
	resp.Header.Del("Cross-Origin-Opener-Policy")
	resp.Header.Del("Cross-Origin-Resource-Policy")
	resp.Header.Del("Expect-Ct")
	resp.Header.Del("Expires")
	resp.Header.Del("Permissions-Policy")
	resp.Header.Del("Pragma")
	resp.Header.Del("Referrer-Policy")
	resp.Header.Del("Server")
	resp.Header.Del("Strict-Transport-Security")
	resp.Header.Del("X-Content-Type-Options")
	resp.Header.Del("X-Frame-Options")
	resp.Header.Del("X-Powered-By")
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
