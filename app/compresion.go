package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
		newBody, err := gzipCompress(res.Body)
		if err != nil {
			fmt.Println("Error compressing response body:", err)
			return
		}

		res.Header["Content-Encoding"] = "gzip"
		res.Body = string(newBody)
		res.Header["Content-Length"] = strconv.Itoa(len(newBody))
	}
}

func gzipCompress(body string) ([]byte, error) {
	var buff bytes.Buffer
	gz := gzip.NewWriter(&buff)

	_, err := gz.Write([]byte(body))
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
