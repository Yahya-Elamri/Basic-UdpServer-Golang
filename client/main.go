package main

import (
	"UdpServer/utils"
	"fmt"
	"net"

	"github.com/eiannone/keyboard"
)

type Player struct {
	ID   int
	X, Y int
}

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

	_, err = conn.Write([]byte("connect"))
	if err != nil {
		fmt.Println("Error sending connect message:", err)
		return
	}

	fmt.Println("Connected to server at", serverAddr)
	buffer := make([]byte, 1024)

	go writer(conn)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}

		player, err := utils.DecodePlayer(buffer[:n])
		if err != nil {
			fmt.Println("Error decoding:", err)
		}

		fmt.Printf("Received player struct: %+v\n", player)
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
