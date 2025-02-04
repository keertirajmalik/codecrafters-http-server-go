package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

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
	defer conn.Close()

	request, _ := http.ReadRequest(bufio.NewReader(conn))

	if request.URL.Path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(request.URL.Path, "/echo") {
		param := strings.Split(request.URL.Path, "/")[2]
		encodings := strings.Split(request.Header.Get("Accept-Encoding"), ",")
		for _, encoding := range encodings {
			if strings.TrimSpace(encoding) == "gzip" {
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s", len(param), param)
				conn.Write([]byte(response))
			}
		}
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(param), param)
		conn.Write([]byte(response))
	} else if request.URL.Path == "/user-agent" {
		data := request.UserAgent()

		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
		conn.Write([]byte(response))
	} else if strings.HasPrefix(request.URL.Path, "/files") {
		args := os.Args
		directoryPath := args[2]
		filePath := directoryPath + strings.Split(request.URL.Path, "/")[2]

		if request.Method == "GET" {
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("file read failed with error:", err)
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return
			}

			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)
			conn.Write([]byte(response))
		}

		if request.Method == "POST" {
			content, err := io.ReadAll(request.Body)
			if err != nil {
				fmt.Println("reading from body failed with error:", err)
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return
			}
			os.WriteFile(filePath, []byte(content), os.ModeAppend)
			conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		}
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
