package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

func ReadRequest(conn net.Conn) (*Request, bool) {
	var closeNext = false
	reader := bufio.NewReader(conn)

	// first line = GET /Path HTTP/1.1
	firstLine, err := reader.ReadString('\n')
	if err == io.EOF {
		return nil, false
	}
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}
	parts := strings.Split(firstLine, " ")
	if len(parts) != 3 {
		fmt.Println("Error invalid request: ", firstLine)
		os.Exit(1)
	}

	headers := make(map[string]string)
	// Headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request: ", err.Error())
			os.Exit(1)
		}
		line = strings.TrimSpace(line)

		// if line break or empty line then header end or doesnt exist
		if line == "\r\n" || line == "" {
			break
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 { // if no colon then it's not header
			continue
		}
		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])
		headers[key] = value
	}

	var buff []byte
	if val, ok := headers["Content-Length"]; ok {
		length, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Error parsing Content-Length: ", err.Error())
		}
		buff = make([]byte, length)
	}

	if val, ok := headers["Connection"]; ok {
		if val == "close" {
			closeNext = true
		}
	}

	var body string
	if buff != nil {
		_, err := reader.Read(buff)
		if err != nil {
			fmt.Println("Error reading Body: ", err.Error())
		}
		body = string(buff)
	}

	return &Request{
		Method:   parts[0],
		Path:     parts[1],
		Protocol: parts[2],
		Headers:  headers,
		Body:     body,
	}, closeNext
}

func HandleRequest(req *Request) *Response {
	var response Response
	if req.Path == "/" {
		response = MakeResponse(Response{})
	} else if strings.HasPrefix(req.Path, "/echo/") {
		param := strings.Split(req.Path, "/")[2]
		response = MakeResponse(Response{Body: param})
	} else if req.Path == "/user-agent" {
		response = MakeResponse(Response{Body: req.Headers["User-Agent"]})
	} else if req.Method == "GET" && strings.HasPrefix(req.Path, "/files/") {
		filePath := fmt.Sprintf("%s%s", StaticFileDirectory, strings.TrimPrefix(req.Path, "/files/"))
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			response = MakeResponse(Response{StatusCode: 404, Body: "File not found"})
		} else {
			response = MakeResponse(Response{StatusCode: 200, Headers: map[string]string{
				"Content-Type": "application/octet-stream",
			}, Body: string(fileContent)})
		}
	} else if req.Method == "POST" && strings.HasPrefix(req.Path, "/files/") {
		filePath := fmt.Sprintf("%s%s", StaticFileDirectory, strings.TrimPrefix(req.Path, "/files/"))
		err := os.WriteFile(filePath, []byte(req.Body), 0644)
		if err != nil {
			response = MakeResponse(Response{StatusCode: 500})
		} else {
			response = MakeResponse(Response{StatusCode: 201})
		}
	} else {
		response = MakeResponse(Response{StatusCode: 404})
	}

	if val, ok := req.Headers["Accept-Encoding"]; ok {
		vals := strings.Split(val, ",")
		for i := range vals {
			vals[i] = strings.TrimSpace(vals[i])
		}
		firstSupported := supportedCompression(vals)
		if firstSupported != "" {
			response.Headers["Content-Encoding"] = firstSupported
		}

		if firstSupported == "gzip" {
			response.Body = gzipCompress(response.Body)
			fmt.Printf("%s", response.Body)
			response.Headers["Content-Length"] = fmt.Sprintf("%d", len(response.Body))
		}
	}

	if val, ok := req.Headers["Connection"]; ok {
		response.Headers["Connection"] = val
	}

	return &response
}
