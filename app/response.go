package main

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
