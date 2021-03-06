package main

import (
	"net/http"
	"strings"
	"io"
	"fmt"
	"os"
	"io/ioutil"
)

var TOKEN = os.Getenv("HEADERLESS_TOKEN")

func makeRequest(queryString map[string][]string) (int, io.ReadCloser, map[string][]string) {
	var headers = map[string][]string{}
	var url interface{}
	var body interface{}
	var tokenCheckDone = false
	var method = "GET"

	for k, v := range queryString {
		switch k {
		case withPrefix("url"):
			url = v[0]
		case withPrefix("body"):
			body = v[0]
		case withPrefix("method"):
			method = v[0]
		case withPrefix("token"):
			if v[0] == TOKEN {
				tokenCheckDone = true
			}
		default:
			headers[k] = v
		}
	}

	if !tokenCheckDone {
		return http.StatusForbidden, makeErrorReadCloser("Access denied"), nil
	}

	if url == nil {
		return http.StatusBadRequest, makeErrorReadCloser(withPrefix("url") + " not set"), nil
	}

	var bodyReader interface{ io.Reader }
	if body != nil {
		bodyReader = strings.NewReader(body.(string))
	} else {
		bodyReader = nil
	}

	fmt.Printf("method = %s, url = %s, body = %s\n", method, url, body)

	req, err := http.NewRequest(method, url.(string), bodyReader)
	if err != nil {
		return http.StatusBadRequest, makeErrorReadCloser(fmt.Sprintf("Cannot prepare request: %s", err.Error())), nil
	}
	for k, v := range headers {
		for _, vv := range v {
			fmt.Printf("add header %s: %s\n", k, vv)
			req.Header.Add(k, vv)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, makeErrorReadCloser(fmt.Sprintf("Error on request: %s", err.Error())), nil
	}

	return resp.StatusCode, resp.Body, resp.Header
}

func withPrefix(name string) string {
	return "headerless_" + name
}

func makeErrorReadCloser(text string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(text))
}

func main() {
	setup(makeRequest)
}
