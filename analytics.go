package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
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

type Analytics struct {
	sync.Mutex
	data map[string]Hits
}

type SafeExpress struct {
	sync.Mutex
	data map[string]uint
}

var (
	analytics = &Analytics{data: map[string]Hits{}}
	express   = &SafeExpress{data: map[string]uint{}}
)

func (e *SafeExpress) IncExpress(dName string) {
	split := strings.Split(dName, ".")
	var truncated string
	if len(split) > 1 {
		truncated = split[len(split)-2] + "." + split[len(split)-1]
	} else {
		truncated = split[len(split)-1]
	}

	e.Lock()
	_, exist := e.data[truncated]
	if exist {
		e.data[truncated]++
	} else {
		e.data[truncated] = 1
	}
	e.Unlock()
}

func (a *Analytics) IncDNS(dName string) {
	a.Lock()
	hit, exist := a.data[dName]
	if exist {
		hit.dns++
		a.data[dName] = hit
	} else {
		a.data[dName] = Hits{dns: 1, http: 0, tls: 0}
	}
	a.Unlock()
}

func (a *Analytics) IncHTTP(dName string) {
	go express.IncExpress(dName)
	a.Lock()
	hit, exist := a.data[dName]
	if exist {
		hit.http++
		a.data[dName] = hit
	} else {
		a.data[dName] = Hits{dns: 0, http: 1, tls: 0}
	}
	a.Unlock()
}

func (a *Analytics) IncTLS(dName string) {
	go express.IncExpress(dName)
	a.Lock()
	hit, exist := a.data[dName]
	if exist {
		hit.tls++
		a.data[dName] = hit
	} else {
		a.data[dName] = Hits{dns: 0, http: 0, tls: 1}
	}
	a.Unlock()
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	fmt.Fprint(w, "<html><head><title>Proxy analytics</title></head><body><table style='float: left;'>"+
		"<thead><th>fqdn</th><th>dns</th><th>http</th><th>tls</th><th>total</th></thead>"+
		"<tbody>")
	for key := range analytics.data {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>",
			key, analytics.data[key].dns, analytics.data[key].http, analytics.data[key].tls,
			analytics.data[key].dns+analytics.data[key].http+analytics.data[key].tls)
	}
	fmt.Fprint(w, "</tbody></table>")

	slice := []Express{}
	for key := range express.data {
		slice = append(slice, Express{domain: key, counter: express.data[key]})
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
