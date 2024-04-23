package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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

	fmt.Printf("Header text: %s\n", headerText)

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
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		handleConnection(l)
	}
}
