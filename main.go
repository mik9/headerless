package main

import (
	"net/http"
	"strings"
	"io"
	"fmt"
	"os"
)

var TOKEN = os.Getenv("HEADERLESS_TOKEN")

func makeRequest(queryString map[string][]string) (int, io.Reader, map[string][]string) {
	var headers = map[string][]string{}
	var url interface{}
	var body interface{}
	var tokenCheckDone = false
	var method = "GET"

	for k, v := range queryString {
		switch k {
		case "headerless_url":
			url = v[0]
		case "headerless_body":
			body = v[0]
		case "headerless_method":
			method = v[0]
		case "headerless_token":
			if v[0] == TOKEN {
				tokenCheckDone = true
			}
		default:
			headers[k] = v
		}
	}

	if !tokenCheckDone {
		return http.StatusForbidden, strings.NewReader("Access denied"), nil
	}

	if url == nil {
		return http.StatusBadRequest, strings.NewReader("url not set"), nil
	}

	var bodyReader interface{ io.Reader }
	if body != nil {
		bodyReader = strings.NewReader(body.(string))
	} else {
		bodyReader = nil
	}

	fmt.Printf("method = %s, url = %s, body = %s\n", method, url, body)

	req, _ := http.NewRequest(method, url.(string), bodyReader)
	for k, v := range headers {
		for _, vv := range v {
			fmt.Printf("add header %s: %s\n", k, vv)
			req.Header.Add(k, vv)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, strings.NewReader(fmt.Sprintf("Error on request: %s", err.Error())), nil
	}

	return resp.StatusCode, resp.Body, resp.Header
}

func main() {
	setup(makeRequest)
}
