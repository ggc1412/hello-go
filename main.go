package main

import (
	"errors"
	"fmt"
	"net/http"
)

type respResult struct {
	url        string
	statusCode int
}

var errRequestFailed = errors.New("Request failed")

func main() {
	var results = make(map[string]int)
	c := make(chan respResult)

	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://academy.nomadcoders.co/",
	}

	for _, url := range urls {
		go hitUrl(url, c)
	}

	for i := 0; i < len(urls); i++ {
		respResult := <-c
		results[respResult.url] = respResult.statusCode
	}

	for url, code := range results {
		fmt.Println(url, code)
	}
}

func hitUrl(url string, c chan<- respResult) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(errRequestFailed)
		resp.StatusCode = 0
	}
	c <- respResult{url: url, statusCode: resp.StatusCode}
}
