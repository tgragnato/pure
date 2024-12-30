package spam

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
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
	length := rand.IntN(10) + 52
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}

	randomUrl := fmt.Sprintf("/%s/%s", addition[:length-30], addition[length-30:])
	if rand.IntN(2) == 1 {
		s.url += "/track"
	}
	s.url += randomUrl
}

func (s *Spam) template() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition) + "_"

	length = 138
	addition = make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition) + "_"

	length = 32
	addition = make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition) + "_"

	length = 10
	addition = make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition)
}

func (s *Spam) insertIcloud() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789."
	length := rand.IntN(20) + 5
	addition := make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += "/rev/" + string(addition) + "@icloud.com/"

	s.url += strconv.Itoa(rand.IntN(9999)) + "/" + strconv.Itoa(rand.IntN(99999999)) + "/"

	characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length = 6
	addition = make([]byte, length)
	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition) + "."

	for i := 0; i < length; i++ {
		addition[i] = characters[rand.IntN(len(characters))]
	}
	s.url += string(addition) + "/0"
}

func (s *Spam) randomPost() {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	name_length := rand.IntN(20) + 3
	name := make([]byte, name_length)
	for i := 0; i < name_length; i++ {
		name[i] = characters[rand.IntN(len(characters))]
	}

	surname_length := rand.IntN(20) + 3
	surname := make([]byte, surname_length)
	for i := 0; i < surname_length; i++ {
		surname[i] = characters[rand.IntN(len(characters))]
	}

	fullname := string(name) + " " + string(surname)
	email := string(name) + "." + string(surname) + "@icloud.com"
	payload := "{\"name\": \"" + string(fullname) + "\", \"email\": \"" + email + "\"}"

	s.body = bytes.NewBufferString(payload)
}
