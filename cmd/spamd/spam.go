package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"text/template"
)

type spam struct {
	url       string
	method    string
	userAgent string
	proxy     bool
}

func (s *spam) call() {
	req, err := http.NewRequestWithContext(context.Background(), s.method, s.url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("User-Agent", s.userAgent)
	var resp *http.Response
	counter++
	if s.proxy {
		resp, err = httpClient.Do(req)
	} else {
		resp, err = directHttp.Do(req)
	}
	counter--
	if err != nil {
		return
	}
	resp.Body.Close()
}

func (s *spam) random() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(10) + 52
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.Intn(len(characters))]
	}

	randomUrl := fmt.Sprintf("/%s/%s", addition[:length-30], addition[length-30:])
	if rand.Intn(2) == 1 {
		s.url += "/track"
	}
	s.url += randomUrl
}

func (s *spam) template() {
	text := templates[rand.Intn(len(templates))]

	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789."
	length := rand.Intn(20) + 5
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.Intn(len(characters))]
	}

	tmpl, err := template.New("template").Parse(text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, struct{ Padding string }{string(addition)})
	if len(buf.Bytes()) == 0 {
		fmt.Println("empty rendered template")
		return
	}

	if rand.Intn(2) == 1 {
		s.url += "?up=" + base64.StdEncoding.EncodeToString(buf.Bytes())
	} else {
		s.url += "?upn=" + base64.StdEncoding.EncodeToString(buf.Bytes())
	}
	if rand.Intn(2) == 1 {
		s.url += "&v=0"
	}
}

func makeSpam() spam {
	switch rand.Intn(6) {

	case 0:
		return spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))] + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}

	case 1:
		s := spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}
		s.random()
		return s

	case 2:
		s := spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}
		s.template()
		return s

	case 3:
		return spam{
			url:       directPrefixes[rand.Intn(len(directPrefixes))] + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: userAgents[0],
			proxy:     false,
		}

	case 4:
		s := spam{
			url:       directPrefixes[rand.Intn(len(directPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[0],
			proxy:     false,
		}
		s.random()
		return s

	case 5:
		s := spam{
			url:       directPrefixes[rand.Intn(len(directPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[0],
			proxy:     false,
		}
		s.template()
		return s
	}

	return spam{}
}
