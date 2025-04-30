package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	method   string
	path     string
	protocol string
	headers  map[string]string
	body     string
}

type Response struct {
	protocol   string
	statusCode int
	statusText string
	headers    map[string]string
	body       string
}

var (
	StaticFileDirectory string
)

func main() {

	flag.StringVar(&StaticFileDirectory, "directory", "", "Static Directory")
	flag.Parse()

	if StaticFileDirectory != "" {
		fmt.Println("Static directory set to ", StaticFileDirectory)
	} else if StaticFileDirectory != "" && StaticFileDirectory[len(StaticFileDirectory)-1:] == "/" { // if no trailing slash, add it
		StaticFileDirectory = StaticFileDirectory + "/"
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Error closing listener: ", err.Error())
		}
	}(l)

	fmt.Println("Listening on :4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go readRequest(conn)
	}
}

func readRequest(conn net.Conn) {

	reader := bufio.NewReader(conn)

	// first line = GET /path HTTP/1.1
	firstLine, err := reader.ReadString('\n')
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
	// headers
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

	var body string
	if buff != nil {
		_, err := reader.Read(buff)
		if err != nil {
			fmt.Println("Error reading body: ", err.Error())
		}
		body = string(buff)
	}

	handleRequest(conn, &Request{
		method:   parts[0],
		path:     parts[1],
		protocol: parts[2],
		headers:  headers,
		body:     body,
	})
}

func handleRequest(conn net.Conn, req *Request) {
	var response Response
	if req.path == "/" {
		response = makeResponse(Response{})
	} else if strings.HasPrefix(req.path, "/echo/") {
		parts := strings.Split(req.path, "/")
		response = makeResponse(Response{body: parts[2]})
	} else if req.path == "/user-agent" {
		response = makeResponse(Response{body: req.headers["User-Agent"]})
	} else if req.method == "GET" && strings.HasPrefix(req.path, "/files/") {
		filePath := fmt.Sprintf("%s%s", StaticFileDirectory, strings.TrimPrefix(req.path, "/files/"))
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			response = makeResponse(Response{statusCode: 404, body: "File not found"})
		} else {
			response = makeResponse(Response{statusCode: 200, headers: map[string]string{
				"Content-Type": "application/octet-stream",
			}, body: string(fileContent)})
		}
	} else if req.method == "POST" && strings.HasPrefix(req.path, "/files/") {
		filePath := fmt.Sprintf("%s%s", StaticFileDirectory, strings.TrimPrefix(req.path, "/files/"))
		err := os.WriteFile(filePath, []byte(req.body), 0644)
		if err != nil {
			response = makeResponse(Response{statusCode: 500})
		} else {
			response = makeResponse(Response{statusCode: 201})
		}
	} else {
		response = makeResponse(Response{statusCode: 404})
	}

	writeResponse(conn, response)
}

func writeResponse(conn net.Conn, res Response) {
	_, err := fmt.Fprintf(conn, "%s %d %s\r\n", res.protocol, res.statusCode, res.statusText)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
	for key, value := range res.headers {
		_, err := fmt.Fprintf(conn, "%s: %s\r\n", key, value)
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
		}
	}
	_, err = fmt.Fprintf(conn, "\r\n%s", res.body) // line break from as the end of headers then body
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}

func makeResponse(res Response) Response {
	if res.protocol == "" {
		res.protocol = "HTTP/1.1" // default
	}

	switch res.statusCode {
	case 0:
		res.statusCode = 200
		res.statusText = "OK"
	case 200:
		res.statusText = "OK"
	case 201:
		res.statusText = "Created"
	case 404:
		res.statusText = "Not Found"
	case 500:
		res.statusText = "Internal Server Error"
	}

	if res.headers == nil {
		res.headers = make(map[string]string)
	}

	if res.body != "" {
		if _, ok := res.headers["Content-Type"]; !ok {
			res.headers["Content-Type"] = "text/plain"
		}
		res.headers["Content-Length"] = fmt.Sprintf("%d", len(res.body))
	}

	return res
}
