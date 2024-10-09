package main

import (
	"fmt"
	"net"

	"github.com/eiannone/keyboard"
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
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyEsc {
			_, err = conn.Write([]byte("quit"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
			break
		} else {
			_, err = conn.Write([]byte(string(char)))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}
