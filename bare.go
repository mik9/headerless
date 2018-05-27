// +build !LAMBDA

package main

import (
	"io"
	"net/http"
	"os"
	"fmt"
)

func setup(handler func(queryString map[string][]string) (int, io.ReadCloser, map[string][]string)) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code, body, headers := handler(r.URL.Query())
		w.WriteHeader(code)
		io.Copy(w, body)
		body.Close()
		if headers != nil {
			for k, v := range headers {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
		}
	})

	cert, certSet := os.LookupEnv("TLS_CERT")
	if certSet {
		fmt.Printf("Using TLS\n")
		key := os.Getenv("TLS_KEY")
		http.ListenAndServeTLS(":8000", cert, key, nil)
	} else {
		fmt.Printf("No TLS\n")
		http.ListenAndServe(":8000", nil)
	}
}
