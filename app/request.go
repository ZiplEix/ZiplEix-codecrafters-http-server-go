package main

import "strings"

type Request struct {
	Method   string
	Path     string
	Protocol string
	Header   map[string]string
	Body     string
}

func parseRequest(data string) Request {
	lines := strings.Split(data, "\r\n")
	requestLine := strings.Split(lines[0], " ")

	method := requestLine[0]
	path := requestLine[1]
	protocol := requestLine[2]

	header := make(map[string]string)
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			body := strings.Join(lines[i+1:], "\r\n")
			return Request{method, path, protocol, header, body}
		}

		headerLine := strings.Split(lines[i], ": ")
		header[headerLine[0]] = headerLine[1]
	}

	return Request{method, path, protocol, header, ""}
}
