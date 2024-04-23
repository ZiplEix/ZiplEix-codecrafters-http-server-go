package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sanity-io/litter"
)

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

func recv(conn net.Conn) Request {
	data := make([]byte, 4096)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	req := string(data)
	request := parseRequest(req)

	litter.Dump(request)

	return request
}

func send(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// read the data from the connection
	req := recv(conn)

	// send the response
	if req.Path == "/" {
		send(conn, "HTTP/1.1 200 OK\r\n\r\n")
	} else {
		send(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}
