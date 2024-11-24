package spam

import (
	"math/rand"
	"net/http"
)

func MakeSpam() Spam {
	switch rand.Intn(7) {

	case 0:
		return Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))] + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}

	case 1:
		s := Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}
		s.random()
		return s

	case 2:
		s := Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))] + paths[rand.Intn(len(paths))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}
		s.template()
		return s

	case 3:
		s := Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))],
			method:    http.MethodGet,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}
		s.insertIcloud()
		return s

	case 4:
		s := Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))],
			method:    http.MethodPost,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}
		s.randomPost()
		return s

	case 5:
		s := Spam{
			url:       proxiedPrefixes[rand.Intn(len(proxiedPrefixes))] + paths[rand.Intn(len(paths))],
			method:    http.MethodPost,
			userAgent: userAgents[rand.Intn(len(userAgents))],
			proxy:     true,
			body:      nil,
		}
		s.randomPost()
		return s

	}

	return Spam{}
}
