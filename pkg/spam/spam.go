package spam

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"text/template"
)

type Spam struct {
	url       string
	method    string
	userAgent string
	proxy     bool
	body      io.Reader
}

func (s *Spam) Call() {
	req, err := http.NewRequestWithContext(context.Background(), s.method, s.url, s.body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Set("User-Agent", s.userAgent)
	if s.method != http.MethodGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}

	var resp *http.Response
	Counter++
	if s.proxy {
		resp, err = httpClient.Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	Counter--
	if err != nil {
		return
	}
	resp.Body.Close()
}

func (s *Spam) random() {
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

func (s *Spam) template() {
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

func (s *Spam) insertIcloud() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789."
	length := rand.Intn(20) + 5
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.Intn(len(characters))]
	}
	s.url += "/rev/" + string(addition) + "@icloud.com/4423/18709008/kjFsTG.XhOZV6/0"
}

func (s *Spam) insertGmail() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789."
	length := rand.Intn(20) + 5
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.Intn(len(characters))]
	}
	s.url += "/rev/" + string(addition) + "@gmail.com/4423/18709008/kjFsTG.XhOZV6/0"
}

func (s *Spam) randomPost() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	name_length := rand.Intn(20) + 3
	name := make([]byte, name_length)
	for i := 0; i < name_length; i++ {
		name[i] = characters[rand.Intn(len(characters))]
	}

	surname_length := rand.Intn(20) + 3
	surname := make([]byte, surname_length)
	for i := 0; i < surname_length; i++ {
		surname[i] = characters[rand.Intn(len(characters))]
	}

	fullname := string(name) + " " + string(surname)
	email := string(name) + "." + string(surname) + "@icloud.com"
	payload := "{\"name\": \"" + string(fullname) + "\", \"email\": \"" + email + "\"}"

	s.body = bytes.NewBufferString(payload)
}
