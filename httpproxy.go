package main

import (
        "io"
        "net/http"
        //"net/url"
        "strings"
)

func main() {
        handler := http.DefaultServeMux

        handler.HandleFunc("/", handleFunc)

        s := &http.Server{
                Addr:    "127.0.0.1:9080",
                Handler: handler,
        }

        s.ListenAndServe()
}

func handleFunc(w http.ResponseWriter, r *http.Request) {

        if (strings.HasSuffix(r.Host, ".apple.com") && r.Host != "ocsp.apple.com" || r.Host == "updates-http.cdn-apple.com") {

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
                        http.Redirect(w, r, "https://" + r.Host + r.URL.RequestURI(), 302)
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

                http.Redirect(w, r, "https://" + r.Host + r.URL.RequestURI(), 301)

        }
}
