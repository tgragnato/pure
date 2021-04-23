package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Hits struct {
	dns  uint
	http uint
	tls  uint
}

type Express struct {
	domain  string
	counter uint
}

var (
	analytics = map[string]Hits{}
	express   = map[string]uint{}
)

func IncExpress(dName string) {
	split := strings.Split(dName, ".")
	var truncated string
	if len(split) > 1 {
		truncated = split[len(split)-2] + "." + split[len(split)-1]
	} else {
		truncated = split[len(split)-1]
	}

	_, exist := express[truncated]
	if exist {
		express[truncated]++
	} else {
		express[truncated] = 1
	}
}

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
	IncExpress(dName)
}

func IncTLS(dName string) {
	hit, exist := analytics[dName]
	if exist {
		hit.tls += 1
		analytics[dName] = hit
	} else {
		analytics[dName] = Hits{dns: 0, http: 0, tls: 1}
	}
	IncExpress(dName)
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	fmt.Fprint(w, "<html><head><title>Proxy analytics</title></head><body><table style='float: left;'>"+
		"<thead><th>fqdn</th><th>dns</th><th>http</th><th>tls</th><th>total</th></thead>"+
		"<tbody>")
	for key := range analytics {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			key, analytics[key].dns, analytics[key].http, analytics[key].tls,
			analytics[key].dns+analytics[key].http+analytics[key].tls)
	}
	fmt.Fprint(w, "</tbody></table>")

	slice := []Express{}
	for key := range express {
		slice = append(slice, Express{domain: key, counter: express[key]})
	}
	sort.Slice(slice, func(i int, j int) bool {
		return slice[i].counter > slice[j].counter
	})

	fmt.Fprint(w, "<table style='float: right;'><thead><th>domain</th><th>count</th></thead><tbody>")
	for i := range slice {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td></tr>", slice[i].domain, slice[i].counter)
	}
	fmt.Fprint(w, "</tbody></table></html>")

}
