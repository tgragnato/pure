package main

import (
	"strings"
)

func Cloaking(qName string) string {
	qName = strings.ToLower(qName)

	switch qName {
	case "youtube.com.":
		qName = "restrictmoderate.youtube.com."
	case "www.youtube.com.":
		qName = "restrictmoderate.youtube.com."
	case "m.youtube.com.":
		qName = "restrictmoderate.youtube.com."
	case "youtubei.googleapis.com.":
		qName = "restrictmoderate.youtube.com."
	case "youtube.googleapis.com.":
		qName = "restrictmoderate.youtube.com."
	case "www.youtube-nocookie.com.":
		qName = "restrictmoderate.youtube.com."
	case "consent.youtube.com.":
		qName = "consent.youtube.com."
	}

	if strings.HasSuffix(qName, "tgragnato.it.") {
		qName = "tgragnato.it."
	}

	if strings.HasSuffix(qName, "github.io.") {
		qName = "github.io."
	}

	return qName
}
