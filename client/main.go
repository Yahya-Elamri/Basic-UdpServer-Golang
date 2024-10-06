package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at", serverAddr)
	buffer := make([]byte, 1024)

	// Start the writer goroutine to handle user input
	go writer(conn)

	// Reader loop on the main thread
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}
		fmt.Printf("\nReceived from %s: %s\n", addr, string(buffer[:n]))
	}
}

func writer(conn *net.UDPConn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message to send (or 'quit' to exit): ")
		if scanner.Scan() {
			message := scanner.Text()
			if message == "quit" {
				_, err := conn.Write([]byte(message))
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
				fmt.Println("Exiting...")
				os.Exit(0)
			}
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		}
	}
}
