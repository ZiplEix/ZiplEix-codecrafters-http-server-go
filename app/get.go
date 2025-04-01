package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func get(conn *net.Conn, req Request, path []string) {
	if req.Path == "/" {
		resp := newResponse(200, "OK", map[string]string{"Content-Type": "text/plain"}, "")
		send(*conn, req, resp)
		return
	}

	if strings.HasPrefix(req.Path, "/echo") {
		respText := strings.Join(path[2:], "/")

		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = strconv.Itoa(len(respText))

		resp := newResponse(200, "OK", header, respText)
		send(*conn, req, resp)
	}

	if strings.HasPrefix(req.Path, "/user-agent") {
		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = strconv.Itoa(len(req.Header["User-Agent"]))

		resp := newResponse(200, "OK", header, req.Header["User-Agent"])
		send(*conn, req, resp)
	}

	if strings.HasPrefix(req.Path, "/files") {
		filePath := strings.Join(path[2:], "/")

		filePath = fmt.Sprintf("%s/%s", *rootDir, filePath)

		file, err := os.ReadFile(filePath)
		if err != nil {
			resp := newResponse(404, "Not Found", map[string]string{"Content-Type": "text/plain"}, "Details: "+err.Error())
			send(*conn, req, resp)
			return
		}

		header := make(map[string]string)
		header["Content-Type"] = "application/octet-stream"
		header["Content-Length"] = strconv.Itoa(len(file))

		resp := newResponse(200, "OK", header, string(file))
		send(*conn, req, resp)
		return
	}

	// if the path is not found
	resp := newResponse(404, "Not Found", map[string]string{"Content-Type": "text/plain"}, "")
	send(*conn, req, resp)
}
