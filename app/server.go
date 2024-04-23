package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
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

func send(conn net.Conn, resp Response) error {
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

	if req.Path == "/" {
		resp := newResponse(200, "OK", map[string]string{"Content-Type": "text/plain"}, "")
		send(*conn, resp)
		return
	}

	if strings.HasPrefix(req.Path, "/echo") {
		respText := strings.Join(path[2:], "/")

		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = strconv.Itoa(len(respText))

		resp := newResponse(200, "OK", header, respText)
		send(*conn, resp)
	}

	if strings.HasPrefix(req.Path, "/user-agent") {
		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = strconv.Itoa(len(req.Header["User-Agent"]))

		resp := newResponse(200, "OK", header, req.Header["User-Agent"])
		send(*conn, resp)
	}

	if strings.HasPrefix(req.Path, "/files") {
		fmt.Println("dans files")

		filePath := strings.Join(path[2:], "/")

		filePath = fmt.Sprintf("%s/%s", *rootDir, filePath)

		file, err := os.ReadFile(filePath)
		if err != nil {
			resp := newResponse(404, "Not Found", map[string]string{"Content-Type": "text/plain"}, "Details: "+err.Error())
			send(*conn, resp)
			return
		}

		header := make(map[string]string)
		header["Content-Type"] = "application/octet-stream"
		header["Content-Length"] = strconv.Itoa(len(file))

		resp := newResponse(200, "OK", header, string(file))
		send(*conn, resp)
		return
	}

	// if the path is not found
	resp := newResponse(404, "Not Found", map[string]string{"Content-Type": "text/plain"}, "")
	send(*conn, resp)
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
