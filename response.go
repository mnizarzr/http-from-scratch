package main

import (
	"fmt"
	"net"
)

type Response struct {
	Protocol   string
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

func WriteResponse(conn net.Conn, res *Response) {
	_, err := fmt.Fprintf(conn, "%s %d %s\r\n", res.Protocol, res.StatusCode, res.StatusText)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
	for key, value := range res.Headers {
		_, err := fmt.Fprintf(conn, "%s: %s\r\n", key, value)
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
		}
	}
	_, err = fmt.Fprintf(conn, "\r\n%s", res.Body) // line break from as the end of Headers then Body
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}

func MakeResponse(res Response) Response {
	if res.Protocol == "" {
		res.Protocol = "HTTP/1.1" // default
	}

	switch res.StatusCode {
	case 0:
		res.StatusCode = 200
		res.StatusText = "OK"
	case 200:
		res.StatusText = "OK"
	case 201:
		res.StatusText = "Created"
	case 404:
		res.StatusText = "Not Found"
	case 500:
		res.StatusText = "Internal Server Error"
	}

	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}

	if res.Body != "" {
		if _, ok := res.Headers["Content-Type"]; !ok {
			res.Headers["Content-Type"] = "text/plain"
		}
		res.Headers["Content-Length"] = fmt.Sprintf("%d", len(res.Body))
	}

	return res
}
