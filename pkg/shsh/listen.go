package shsh

import (
	"log"
	"net/http"
)

func Listen(iface4 string, iface6 string) {
	handler := http.DefaultServeMux
	handler.HandleFunc("/", handleHTTPForward)
	go func() {
		err := http.ListenAndServe(iface4+":80", handler)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	go func() {
		err := http.ListenAndServe(iface6+":80", handler)
		if err != nil {
			log.Println(err.Error())
		}
	}()
}
