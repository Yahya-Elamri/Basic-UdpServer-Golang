package main

import (
	"fmt"
	"net"
	"sync"
)

const (
	addr       = "127.0.0.1:3000"
	bufferSize = 1024
)

type Player struct {
	id   int
	x, y int
	hp   int
	addr *net.UDPAddr
}

type GameState struct {
	round   int
	players map[string]*Player
}

type Server struct {
	Game  GameState
	conn  *net.UDPConn
	mutex sync.Mutex
}

func NewServer() (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn: conn,
		Game: GameState{
			round:   0,
			players: make(map[string]*Player),
		},
	}, nil

}

func (s *Server) Read() {
	buffer := make([]byte, bufferSize)

	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		message := string(buffer[:n])
		fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr, message)

		// Register client if it's not already in the list
		s.mutex.Lock()
		if _, ok := s.Game.players[remoteAddr.String()]; !ok {
			PlayerId := len(s.Game.players) + 1
			s.Game.players[remoteAddr.String()] = &Player{id: PlayerId, addr: remoteAddr, x: 5, y: 5, hp: 100}
			fmt.Printf("New client connected: %s (ID: %d)\n", remoteAddr, PlayerId)
		}
		s.mutex.Unlock()

		// Broadcast the message to all clients
		s.Write(message, remoteAddr)
	}
}

func (s *Server) Write(message string, sender *net.UDPAddr) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, client := range s.Game.players {
		if client.addr.String() != sender.String() { // Don't send back to the sender
			_, err := s.conn.WriteToUDP([]byte(message), client.addr)
			if err != nil {
				fmt.Printf("Error sending to %s: %v\n", client.addr, err)
			} else {
				fmt.Printf("Sent message to %s: %s\n", client.addr, message)
			}
		}
	}
}

func main() {
	server, err := NewServer()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer server.conn.Close()

	fmt.Println("Server started on", server.conn.LocalAddr())

	// Start reading and broadcasting messages
	server.Read()
}
