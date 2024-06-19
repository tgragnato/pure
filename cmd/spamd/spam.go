package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
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
	counter++
	req.Header.Set("User-Agent", s.userAgent)
	var resp *http.Response
	if s.proxy {
		resp, err = httpClient.Do(req)
	} else {
		resp, err = directHttp.Do(req)
	}
	if err != nil {
		counter--
		return
	}
	resp.Body.Close()
	counter--
}

func (s *spam) random() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(10) + 60
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.Intn(len(characters))]
	}

	randomUrl := fmt.Sprintf("/%s/%s#%s", addition[:length-30], addition[length-30:length-22], addition[length-22:])
	if rand.Intn(2) == 1 {
		s.url += "/track"
	}
	s.url += randomUrl
}

func makeSpam() spam {
	switch rand.Intn(4) {

	case 0:
		return spam{
			url:       domains[rand.Intn(len(domains))] + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}

	case 1:
		return spam{
			url:       "http://162.62.230.106" + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
			proxy:     false,
		}

	case 2:
		s := spam{
			url:       freshDomains[rand.Intn(len(freshDomains))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}
		s.random()
		return s

	case 3:
		s := spam{
			url:       freshDomains[rand.Intn(len(freshDomains))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
		}
		if rand.Intn(2) == 1 {
			s.random()
		} else {
			s.url += paths[rand.Intn(len(paths))]
		}
		return s
	}

	return spam{}
}
