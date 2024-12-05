package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Start listening on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected.")

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the HTTP request
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}

	// Log the request line
	fmt.Println("Request received:", strings.TrimSpace(requestLine))

	// Read and discard the headers (until an empty line)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers:", err)
			return
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}

	// Write a HTTP 200 OK response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello, World!"

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}

	fmt.Println("Response sent to client.")
}
