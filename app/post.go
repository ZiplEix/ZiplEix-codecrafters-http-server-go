package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func post(conn *net.Conn, req Request, path []string) {
	if strings.HasPrefix(req.Path, "/files") {
		filePath := strings.Join(path[2:], "/")
		filePath = fmt.Sprintf("%s/%s", *rootDir, filePath)

		err := os.WriteFile(filePath, []byte(req.Body), 0644)
		if err != nil {
			resp := newResponse(500, "Internal server error", map[string]string{"Content-Type": "text/plain"}, "Details: "+err.Error())
			send(*conn, req, resp)
			return
		}

		header := make(map[string]string)
		header["Content-Type"] = "text/plain"
		header["Content-Length"] = "0"

		resp := newResponse(201, "Created", header, "")
		send(*conn, req, resp)
		return
	}

	// if the path is not found
	resp := newResponse(404, "Not Found", map[string]string{"Content-Type": "text/plain"}, "")
	send(*conn, req, resp)
}
