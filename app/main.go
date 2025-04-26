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

	res := "HTTP/1.1 404 Not Found\r\n\r\n"
	if path == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	}

	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

}
