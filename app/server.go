package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	rootDir *string
)

func recv(conn net.Conn) (Request, error) {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		return Request{}, fmt.Errorf("Error reading data: %s", err.Error())
	}

	// sanitize the data
	data = []byte(strings.Trim(string(data), "\x00"))

	req := string(data)
	request := parseRequest(req)

	return request, nil
}

func formatHeader(header map[string]string) string {
	headerText := ""
	for key, value := range header {
		headerText += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return headerText
}

func send(conn net.Conn, req Request, resp Response) error {
	if canCompress(req) {
		compress(req, &resp)
	}

	data := fmt.Sprintf("%s %d %s\r\n%s\r\n%s", resp.Protocol, resp.ReturnCode, resp.Status, formatHeader(resp.Header), resp.Body)

	_, err := conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("Error sending data: %s", err.Error())
	}

	return nil
}

func processRequest(conn *net.Conn) {
	req, err := recv(*conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	path := strings.Split(req.Path, "/")

	if req.Method == "GET" {
		get(conn, req, path)
	} else if req.Method == "POST" {
		post(conn, req, path)
	} else {
		resp := newResponse(405, "Method Not Allowed", map[string]string{"Content-Type": "text/plain"}, "")
		send(*conn, req, resp)
	}
}

func handleConnection(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()
		processRequest(&conn)
	}()
}

func main() {

	rootDir = flag.String("directory", "", "Directory to serve files from")

	if rootDir == nil || *rootDir == "" {
		*rootDir = "./"
	}

	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		handleConnection(l)
	}
}
