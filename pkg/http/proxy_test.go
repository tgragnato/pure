package http

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Fatal("Failed to write response:", err)
		}
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	rec := httptest.NewRecorder()

	gzipMiddleware(handler).ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Expected Content-Encoding to be gzip")
	}

	if rec.Header().Get("Vary") != "Accept-Encoding" {
		t.Error("Expected Vary header to be Accept-Encoding")
	}

	reader, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal("Failed to create gzip reader:", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal("Failed to read gzipped content:", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content' but got '%s'", string(content))
	}
}

func TestFlateMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Fatal("Failed to write response:", err)
		}
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	rec := httptest.NewRecorder()

	flateMiddleware(handler).ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "deflate" {
		t.Error("Expected Content-Encoding to be deflate")
	}

	if rec.Header().Get("Vary") != "Accept-Encoding" {
		t.Error("Expected Vary header to be Accept-Encoding")
	}

	reader := flate.NewReader(rec.Body)
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal("Failed to read deflated content:", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content' but got '%s'", string(content))
	}
}

func TestCompressMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Fatal("Failed to write response:", err)
		}
	})

	tests := []struct {
		name           string
		acceptEncoding string
		wantEncoding   string
		wantContent    string
		wantError      bool
	}{
		{
			name:           "gzip encoding",
			acceptEncoding: "gzip",
			wantEncoding:   "gzip",
			wantContent:    "test content",
		},
		{
			name:           "deflate encoding",
			acceptEncoding: "deflate",
			wantEncoding:   "deflate",
			wantContent:    "test content",
		},
		{
			name:           "no encoding",
			acceptEncoding: "",
			wantContent:    "test content",
		},
		{
			name:           "unsupported encoding",
			acceptEncoding: "br",
			wantContent:    "test content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "http://example.com", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}
			rec := httptest.NewRecorder()

			compressMiddleware(handler).ServeHTTP(rec, req)

			var content []byte
			var err error

			switch tt.wantEncoding {
			case "gzip":
				if rec.Header().Get("Content-Encoding") != "gzip" {
					t.Error("Expected Content-Encoding to be gzip")
				}
				var reader *gzip.Reader
				reader, err = gzip.NewReader(rec.Body)
				if err != nil {
					t.Fatal("Failed to create gzip reader:", err)
				}
				defer reader.Close()
				content, err = io.ReadAll(reader)
			case "deflate":
				if rec.Header().Get("Content-Encoding") != "deflate" {
					t.Error("Expected Content-Encoding to be deflate")
				}
				reader := flate.NewReader(rec.Body)
				defer reader.Close()
				content, err = io.ReadAll(reader)
			default:
				content, err = io.ReadAll(rec.Body)
			}

			if err != nil && !tt.wantError {
				t.Fatal("Failed to read content:", err)
			}

			if string(content) != tt.wantContent {
				t.Errorf("Expected '%s' but got '%s'", tt.wantContent, string(content))
			}
		})
	}
}

func TestHeadersMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test content")); err != nil {
			t.Fatal("Failed to write response:", err)
		}
	})

	tests := []struct {
		name              string
		existingHeaders   map[string]string
		expectedHeaders   map[string]string
		expectedCookieLen uint64
	}{
		{
			name: "default headers",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://tgragnato.it",
				"Cache-Control":               "no-store, no-cache, must-revalidate, post-check=0, pre-check=0",
				"Content-Security-Policy":     "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' blob:; img-src 'self' data:; worker-src 'self' blob:; object-src 'none'; upgrade-insecure-requests;",
				"Expect-Ct":                   "max-age=86400, enforce",
				"Expires":                     "Thu, 19 Nov 1981 08:52:00 GMT",
				"Permissions-Policy":          "interest-cohort=(), accelerometer=(), autoplay=(), camera=(), clipboard-read=(), clipboard-write=(), document-domain=(), encrypted-media=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=(), gamepad=(), vibrate=(), vr=(), xr-spatial-tracking=()",
				"Pragma":                      "no-cache",
				"Referrer-Policy":             "strict-origin-when-cross-origin",
				"Server":                      "Apache/2.4.25 (RedStar4.0) OpenSSL/1.0.1e-fips PHP/5.6.2",
				"Strict-Transport-Security":   "max-age=31536000; includeSubDomains; preload",
				"X-Content-Type-Options":      "nosniff",
				"X-Frame-Options":             "DENY",
				"X-Powered-By":                "PHP/5.6.2",
				"X-XSS-Protection":            "1; mode=block",
			},
			expectedCookieLen: 26,
		},
		{
			name: "with existing cookie header",
			existingHeaders: map[string]string{
				"Set-Cookie": "existingcookie=value",
			},
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://tgragnato.it",
				"Cache-Control":               "no-store, no-cache, must-revalidate, post-check=0, pre-check=0",
				"Content-Security-Policy":     "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' blob:; img-src 'self' data:; worker-src 'self' blob:; object-src 'none'; upgrade-insecure-requests;",
				"Expect-Ct":                   "max-age=86400, enforce",
				"Expires":                     "Thu, 19 Nov 1981 08:52:00 GMT",
				"Permissions-Policy":          "interest-cohort=(), accelerometer=(), autoplay=(), camera=(), clipboard-read=(), clipboard-write=(), document-domain=(), encrypted-media=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=(), gamepad=(), vibrate=(), vr=(), xr-spatial-tracking=()",
				"Pragma":                      "no-cache",
				"Referrer-Policy":             "strict-origin-when-cross-origin",
				"Server":                      "Apache/2.4.25 (RedStar4.0) OpenSSL/1.0.1e-fips PHP/5.6.2",
				"Strict-Transport-Security":   "max-age=31536000; includeSubDomains; preload",
				"X-Content-Type-Options":      "nosniff",
				"X-Frame-Options":             "DENY",
				"X-Powered-By":                "PHP/5.6.2",
				"X-XSS-Protection":            "1; mode=block",
				"Set-Cookie":                  "existingcookie=value",
			},
			expectedCookieLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "http://example.com", nil)
			rec := httptest.NewRecorder()

			for k, v := range tt.existingHeaders {
				rec.Header().Set(k, v)
			}

			headersMiddleware(handler).ServeHTTP(rec, req)

			for k, v := range tt.expectedHeaders {
				if got := rec.Header().Get(k); got != v {
					t.Errorf("Expected header %s to be %s, got %s", k, v, got)
				}
			}

			if tt.expectedCookieLen > 0 {
				cookie := rec.Header().Get("Set-Cookie")
				if !strings.HasPrefix(cookie, "PHPSESSID=") {
					t.Error("Expected Set-Cookie header to start with PHPSESSID=")
				}
				sessID := strings.TrimPrefix(strings.Split(cookie, ";")[0], "PHPSESSID=")
				if uint64(len(sessID)) != tt.expectedCookieLen {
					t.Errorf("Expected session ID length to be %d, got %d", tt.expectedCookieLen, len(sessID))
				}
			}
		})
	}
}

func TestStripProxiedHeaders(t *testing.T) {
	t.Parallel()

	headers := map[string][]string{
		"Expect-Ct":                 {"max-age=86400, enforce"},
		"Permissions-Policy":        {"interest-cohort=()"},
		"Referrer-Policy":           {"strict-origin-when-cross-origin"},
		"Server":                    {"Apache/2.4.25"},
		"Strict-Transport-Security": {"max-age=31536000"},
		"X-Content-Type-Options":    {"nosniff"},
		"X-Frame-Options":           {"DENY"},
		"X-XSS-Protection":          {"1; mode=block"},
		"Content-Type":              {"text/html"},
	}

	resp := &http.Response{
		Header: headers,
	}

	err := stripProxiedHeaders(resp)
	if err != nil {
		t.Errorf("stripProxiedHeaders() returned error: %v", err)
	}

	for header := range headers {
		if header == "Content-Type" {
			if resp.Header.Get(header) == "" {
				t.Errorf("Expected header %s to be preserved", header)
			}
			continue
		}
		if resp.Header.Get(header) != "" {
			t.Errorf("Expected header %s to be stripped", header)
		}
	}
}

func TestAPIGateway(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		path       string
		wantCode   int
		wantLoc    string
	}{
		{
			name:       "non-private IP redirect",
			remoteAddr: "8.8.8.8:1234",
			path:       "/test",
			wantCode:   http.StatusFound,
			wantLoc:    "https://tgragnato.it/test",
		},
		{
			name:       "private IP transmission base path",
			remoteAddr: "192.168.1.1:1234",
			path:       "/transmission",
			wantCode:   http.StatusFound,
			wantLoc:    "https://api.tgragnato.it/transmission/web/",
		},
		{
			name:       "private IP transmission subpath",
			remoteAddr: "192.168.1.1:1234",
			path:       "/transmission/rpc",
			wantCode:   http.StatusBadGateway,
		},
		{
			name:       "private IP other path",
			remoteAddr: "192.168.1.1:1234",
			path:       "/other",
			wantCode:   http.StatusBadGateway,
		},
		{
			name:       "private IP grafana base path",
			remoteAddr: "192.168.1.1:1234",
			path:       "/grafana/",
			wantCode:   http.StatusBadGateway,
		},
		{
			name:       "invalid remote addr",
			remoteAddr: "invalid",
			path:       "/test",
			wantCode:   http.StatusFound,
			wantLoc:    "https://tgragnato.it/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "http://example.com"+tt.path, nil)
			req.RemoteAddr = tt.remoteAddr
			rec := httptest.NewRecorder()

			handler := apiGateway()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("Expected status code %d, got %d", tt.wantCode, rec.Code)
			}

			if tt.wantLoc != "" {
				if got := rec.Header().Get("Location"); got != tt.wantLoc {
					t.Errorf("Expected Location header %s, got %s", tt.wantLoc, got)
				}
			}
		})
	}
}

func TestResponseWriterWriteHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     int
		wantCode int
	}{
		{
			name:     "success code",
			code:     http.StatusOK,
			wantCode: http.StatusOK,
		},
		{
			name:     "error code",
			code:     http.StatusBadRequest,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "redirect code",
			code:     http.StatusFound,
			wantCode: http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			w := &responseWriter{ResponseWriter: rec}
			w.WriteHeader(tt.code)

			if w.status != tt.wantCode {
				t.Errorf("responseWriter.status = %v, want %v", w.status, tt.wantCode)
			}
			if rec.Code != tt.wantCode {
				t.Errorf("ResponseWriter.Code = %v, want %v", rec.Code, tt.wantCode)
			}
		})
	}
}

func TestResponseWriterWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []byte
		customWrite bool
		wantSize    int
		wantErr     bool
	}{
		{
			name:        "default writer",
			input:       []byte("test content"),
			customWrite: false,
			wantSize:    12,
			wantErr:     false,
		},
		{
			name:        "custom writer",
			input:       []byte("custom content"),
			customWrite: true,
			wantSize:    14,
			wantErr:     false,
		},
		{
			name:        "empty input",
			input:       []byte{},
			customWrite: false,
			wantSize:    0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			var w *responseWriter
			if tt.customWrite {
				w = &responseWriter{
					ResponseWriter: rec,
					Writer:         rec.Body,
				}
			} else {
				w = &responseWriter{
					ResponseWriter: rec,
				}
			}

			size, err := w.Write(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("responseWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if size != tt.wantSize {
				t.Errorf("responseWriter.Write() size = %v, want %v", size, tt.wantSize)
			}
			if int64(size) != w.size {
				t.Errorf("responseWriter.size = %v, want %v", w.size, size)
			}
			if got := rec.Body.String(); got != string(tt.input) {
				t.Errorf("Written content = %v, want %v", got, string(tt.input))
			}
		})
	}
}
