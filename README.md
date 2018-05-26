# Headerless

Simple proxy for evnironments where you can do HTTP
requests but it's hard or imposible to use (like Mikrotik Router OS).

Can be used [directly](#directly), with [docker](#docker) or with [Amazon Lambda](#amazon-lambda).

## Usage

Pass any values as query and they will be translated into headers.
Except next:
* `headerless_token` - authorization, _required_.
* `headerless_url` - url to request, _required_.
* `headerless_method` - http method used to query URL, _optional_. By defaut `GET`.
* `headerless_body` - body to be sent to URL, _optional_.

For example:
```
http://127.0.0.1:8000/?headerless_token=secret&headerless_url=http://httpbin.org/anything&test_header=test&User-Agent=me&headerless_method=POST&headerless_body=123
```
Will be translated to:
```
POST /anything HTTP/1.1

Accept-Encoding: gzip
Connection: close
Content-Length: 3
Host: httpbin.org
Test-Header: test
User-Agent: me

123
```

## Directly

Build with:
```go
go get && go build
```

Run with:
```
HEADERLESS_TOKEN=secret ./headerless
```

## Docker
Using prepader image:
```
docker run -d -p 8000:8000 -e HEADERLESS_TOKEN=secret mik9/headerless
```

Or build youresf:
```
docker build -t headerless .
docker run -d -p 8000:8000 -e HEADERLESS_TOKEN=secret headerless
```

## Amazon Lambda

Build for Lambda:
```
GOOS=linux GOARCH=amd64 go build --tags LAMBDA -o lambda
```

Deploy to Lambda with any way you like, for example:
```
zip lambda.zip ./lambda
aws lambda update-function-code \
    --function-name YOUR_LAMBDA \
    --zip-file fileb://lambda.zip
```
