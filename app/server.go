package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Printf("received %d bytes\n", n)
	fmt.Printf("received the following data: %s", string(buf[:n]))

	requestLine := strings.Split(string(buf), "\r\n")
	fmt.Println(requestLine[0])

	request := strings.Fields(requestLine[0])
	method := request[0]
	path := request[1]

	if method == "GET" {
		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else if strings.HasPrefix(path, "/echo") {
			param := strings.SplitN(path, "/", 3)[2]

			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(param), param)
			conn.Write([]byte(response))
		} else if path == "/user-agent" {
			var userAgent string
			for _, header := range requestLine {
				if strings.HasPrefix(header, "User-Agent") {
					userAgent = header
					break
				}
			}

			data := strings.Split(strings.Fields(userAgent)[1], "\r\n")[0]
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			conn.Write([]byte(response))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
}
