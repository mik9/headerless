// +build LAMBDA

package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"context"
	"io"
	"io/ioutil"
)

type MyEvent struct {
	QueryStringParameters map[string]string `json:"queryStringParameters"`
}

type MyResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}

func mapToMultimap(m map[string]string) map[string][]string {
	mm := map[string][]string {}

	for k, v := range m {
		mm[k] = []string {v}
	}

	return mm
}

func multimapToMap(mm map[string][]string) map[string]string {
	m := map[string]string{}

	for k, v := range mm {
		m[k] = v[0]
	}

	return m
}

func setup(handler func(queryString map[string][]string) (int, io.Reader, map[string][]string)) {
	lambda.Start(func(ctx context.Context, event MyEvent) (MyResponse, error) {
		code, body, headers := handler(mapToMultimap(event.QueryStringParameters))

		b, err := ioutil.ReadAll(body)
		if err != nil {
			panic(err)
		}

		return MyResponse{string(b), multimapToMap(headers), code}, nil
	})
}
