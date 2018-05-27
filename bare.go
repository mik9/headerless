// +build !LAMBDA

package main

import (
	"io"
	"net/http"
	"os"
	"fmt"
)

func setup(handler func(queryString map[string][]string) (int, io.ReadCloser, map[string][]string)) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	key, keySet := os.LookupEnv("TLS_KEY")

	if certSet != keySet {
		if !certSet {
			panic("TLS_KEY set but TLS_CERT not set")
		}
		if !keySet {
			panic("TLS_CERT set but TLS_KEY not set")
		}
	}

	server := http.Server{Addr: ":8000", Handler: handlerFunc}
	var err error
	if certSet && keySet {
		fmt.Printf("Using TLS\n")
		err = server.ListenAndServeTLS(cert, key)
	} else {
		fmt.Printf("No TLS\n")
		err = server.ListenAndServe()
	}

	if err != nil {
		panic(err)
	}
}
