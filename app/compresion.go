package main

import (
	"strconv"
	"strings"
)

func canCompress(req Request) bool {
	if _, exist := req.Header["Accept-Encoding"]; exist {
		if strings.Contains(req.Header["Accept-Encoding"], "gzip") {
			return true
		}
	}

	return false
}

func compress(req Request, res *Response) {
	if res.Header == nil {
		res.Header = make(map[string]string)
	}

	if strings.Contains(req.Header["Accept-Encoding"], "gzip") {
		res.Header["Content-Encoding"] = "gzip"

		res.Body = string(gzipCompress(res.Body))
		res.Header["Content-Length"] = strconv.Itoa(len(res.Body))
	}
}

func gzipCompress(body string) []byte {
	return []byte(body)
}
