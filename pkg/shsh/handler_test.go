package shsh

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHTTPForward(t *testing.T) {
	t.Parallel()

	// Assuming there's not a proxy server running in the test environment
	tests := []struct {
		host       string
		path       string
		expected   string
		statusCode int
	}{
		{
			host:       "static.ess.apple.com",
			path:       "/connectivity.txt",
			expected:   "AV was here!",
			statusCode: http.StatusOK,
		},
		{
			host:       "updates-http.cdn-apple.com",
			path:       "/",
			expected:   "<a href=\"https://updates-http.cdn-apple.com/\">Found</a>.\n\n",
			statusCode: http.StatusFound,
		},
		{
			host:       "www.example.com",
			path:       "/",
			expected:   "<a href=\"https://www.example.com/\">Moved Permanently</a>.\n\n",
			statusCode: http.StatusMovedPermanently,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Host = test.host

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleHTTPForward)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.statusCode {
			t.Errorf("For host %s and path %s, expected status code %d, but got %d", test.host, test.path, test.statusCode, rr.Code)
		}

		if rr.Body.String() != test.expected {
			t.Errorf("For host %s and path %s, expected response body %q, but got %q", test.host, test.path, test.expected, rr.Body.String())
		}
	}
}
