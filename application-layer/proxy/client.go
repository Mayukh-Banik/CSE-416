package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	// Start listening on port 9900
	listener, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatalf("Error starting listener on port 9900: %v", err)
	}
	defer listener.Close()
	fmt.Println("Listening for incoming data on port 9900...")

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection established.")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buffer bytes.Buffer
	reader := bufio.NewReader(conn)

	// Keep reading from the connection
	for {
		// Read data from the connection
		temp := make([]byte, 1024) // Temporary buffer
		n, err := reader.Read(temp)
		if n > 0 {
			buffer.Write(temp[:n])
			fmt.Printf("Received data: %s\n", buffer.String())
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading data:", err)
			return
		}

		// Try to parse the HTTP request
		req, err := http.ReadRequest(bufio.NewReader(&buffer))
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				// Incomplete request, keep reading
				continue
			}
			log.Println("Error parsing HTTP request:", err)
			return
		}

		// Forward the request to localhost:8080
		forwardToLocalhost(req, conn)
		break
	}
}

func forwardToLocalhost(req *http.Request, conn net.Conn) {
	// Dial localhost:8080
	dialConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Printf("Error connecting to localhost:8080: %v", err)
		return
	}
	defer dialConn.Close()

	// Write the request to localhost:8080
	err = req.Write(dialConn)
	if err != nil {
		log.Printf("Error forwarding request: %v", err)
		return
	}

	// Read the response from localhost:8080
	respReader := bufio.NewReader(dialConn)
	resp, err := http.ReadResponse(respReader, req)
	if err != nil {
		log.Printf("Error reading response from localhost:8080: %v", err)
		return
	}

	// Send the response back to the original connection
	err = resp.Write(conn)
	if err != nil {
		log.Printf("Error writing response back to client: %v", err)
		return
	}
}
