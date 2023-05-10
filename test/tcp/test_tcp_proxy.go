package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	addr := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v\n", addr, err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Failed to read data: %v\n", err)
			}
			return
		}

		fmt.Printf("Received request: %s\n", string(buf[:n]))
		conn.Write([]byte("pong\n"))
	}
}
