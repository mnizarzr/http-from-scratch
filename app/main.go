package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

func main() {

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

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
	}

	req := string(buff)
	lines := strings.Split(req, CRLF)
	path := strings.Split(lines[0], " ")[1] // 0 = method, 1 = path, 2 = protocol and or version
	parameter := strings.Split(path, "/")

	resStatus := "HTTP/1.1 404 Not Found" + CRLF
	resHeaders := CRLF
	var resBody string

	if path == "/" {
		resStatus = "HTTP/1.1 200 OK" + CRLF
	} else if strings.HasPrefix(path, "/echo/") {
		resStatus = "HTTP/1.1 200 OK" + CRLF
		resHeaders = "Content-Type: text/plain" + CRLF
		resHeaders += "Content-Length: " + fmt.Sprintf("%d", len(parameter[2])) + CRLF
		resHeaders += CRLF
		resBody = parameter[2]
	}

	_, err = conn.Write([]byte(resStatus + resHeaders + resBody))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

}
