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
	// Start a listener on port 9900
	listen, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	fmt.Println("Listening on port 9900...")

	for {
		// Accept a new connection
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		handleConnection(conn)
	}
}

// Handle incoming connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buffer bytes.Buffer // Create a dynamic buffer using bytes.Buffer
	var req *http.Request
	temp := make([]byte, 1024) // Temporary buffer for reading data

	for {
		// Read data from the connection
		n, err := conn.Read(temp)
		if n > 0 {
			buffer.Write(temp[:n]) // Append the read data to the buffer
			fmt.Printf("Received data: %s\n", buffer.String())
		}

		if err != nil {
			if err != io.EOF {
				log.Println("Error reading data:", err)
			}
			break
		}

		req, err = http.ReadRequest(bufio.NewReader(&buffer))
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		} else {
			break
		}
	}

	// Print the parsed HTTP request details
	fmt.Println("Method:", req.Method)
	fmt.Println("URL:", req.URL)
	fmt.Println("Headers:", req.Header)
	fmt.Println("Body:", req.Body)

	responseBody := "This is a custom 200 OK response"
	resp := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
	}

	// Set headers for the response
	resp.Header.Set("Content-Type", "text/plain")
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(responseBody)))

	// Write the HTTP response to the connection
	err := resp.Write(conn)
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}

	fmt.Println("200 OK response sent.")
}
