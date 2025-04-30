package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		request, closeConn := ReadRequest(conn)
		if request == nil {
			break
		}
		response := HandleRequest(request)
		WriteResponse(conn, response)
		if closeConn {
			err := conn.Close()
			if err != nil {
				fmt.Println("Error closing connection: ", err.Error())
			}
			break
		}
	}
}
