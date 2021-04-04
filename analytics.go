package main

import (
	"fmt"
	"net/http"
)

type Hits struct {
	dns  uint
	http uint
	tls  uint
}

var analytics = map[string]Hits{}

func IncDNS(dName string) {
	hit, exist := analytics[dName]
	if exist {
		hit.dns += 1
		analytics[dName] = hit
	} else {
		analytics[dName] = Hits{dns: 1, http: 0, tls: 0}
	}
}

func IncHTTP(dName string) {
	hit, exist := analytics[dName]
	if exist {
		hit.http += 1
		analytics[dName] = hit
	} else {
		analytics[dName] = Hits{dns: 0, http: 1, tls: 0}
	}
}

func IncTLS(dName string) {
	hit, exist := analytics[dName]
	if exist {
		hit.tls += 1
		analytics[dName] = hit
	} else {
		analytics[dName] = Hits{dns: 0, http: 0, tls: 1}
	}
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	fmt.Fprint(w, "<html><head><title>Proxy analytics</title></head><body><table>"+
		"<thead><th>fqdn</th><th>dns</th><th>http</th><th>tls</th><th>total</th></thead>"+
		"<tbody>")
	for key := range analytics {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			key, analytics[key].dns, analytics[key].http, analytics[key].tls,
			analytics[key].dns+analytics[key].http+analytics[key].tls)
	}
	fmt.Fprint(w, "</tbody></table></html>")
}
