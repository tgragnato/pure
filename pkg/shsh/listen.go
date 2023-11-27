package shsh

import (
	"net/http"
)

func Listen(iface4 string, iface6 string) {
	handler := http.DefaultServeMux
	handler.HandleFunc("/", handleHTTPForward)
	go http.ListenAndServe(iface4+":80", handler)
	go http.ListenAndServe(iface6+":80", handler)
}
