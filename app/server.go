package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method   string
	Path     string
	Protocol string
	Header   map[string]string
	Body     string
}

type Response struct {
	Protocol   string
	ReturnCode int
	Status     string
	Header     map[string]string
	Body       string
}

func newResponse(returnCode int, status string, header map[string]string, body string) Response {
	return Response{"HTTP/1.1", returnCode, status, header, body}
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

	// litter.Dump(request)

	return request
}

func formatHeader(header map[string]string) string {
	headerText := ""
	for key, value := range header {
		headerText += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return headerText
}

func send(conn net.Conn, resp Response) {
	data := fmt.Sprintf("%s %d %s\r\n%s\r\n%s", resp.Protocol, resp.ReturnCode, resp.Status, formatHeader(resp.Header), resp.Body)

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

	// get the mast part of the path
	path := strings.Split(req.Path, "/")

	if req.Path == "/" {
		resp := newResponse(200, "OK", map[string]string{"Content-Type": "text/plain"}, "")
		send(conn, resp)
		return
	}

	if len(path) < 3 {
		resp := newResponse(404, "Not Found, path to short", map[string]string{"Content-Type": "text/plain"}, "Not Found")
		send(conn, resp)
		return
	} else if path[1] != "echo" {
		resp := newResponse(404, "Not Found, no echo", map[string]string{"Content-Type": "text/plain"}, "Not Found")
		send(conn, resp)
		return
	} else {
		respText := path[2]

		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = strconv.Itoa(len(respText))

		resp := newResponse(200, "OK", header, respText)
		send(conn, resp)
	}
}
