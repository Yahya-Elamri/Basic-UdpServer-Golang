package main

import (
	"UdpServer/module"
	"UdpServer/utils"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	addr       = "127.0.0.1:3000"
	bufferSize = 1024
	tickRate   = 128
)

type GameState struct {
	round   int
	players map[string]*module.Player
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
			players: make(map[string]*module.Player),
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

		go s.HandleClients(remoteAddr, message)
	}
}

func (s *Server) HandleClients(remoteAddr *net.UDPAddr, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if message == "connect" {
		if _, ok := s.Game.players[remoteAddr.String()]; !ok {
			PlayerId := len(s.Game.players) + 1
			s.Game.players[remoteAddr.String()] = &module.Player{ID: PlayerId, Addr: remoteAddr, X: 0, Y: 0}
			fmt.Printf("New client connected: %s (ID: %d)\n", remoteAddr, PlayerId)
		}
	} else if message == "quit" {
		delete(s.Game.players, remoteAddr.String())
		fmt.Printf("Client disconnected: %s\n", remoteAddr)
	} else {
		_, err := s.handleGameLogic(message, remoteAddr)
		if err != nil {
			fmt.Printf("Error processing game logic for %s: %v\n", remoteAddr, err)
		}
	}
}

func (s *Server) handleGameLogic(message string, sender *net.UDPAddr) ([]byte, error) {
	player, exists := s.Game.players[sender.String()]
	if !exists {
		fmt.Printf("Player not found for address %s\n", sender.String())
		return nil, fmt.Errorf("player not found")
	}

	switch message {
	case "z":
		player.Y += 1
	case "s":
		player.Y -= 1
	case "d":
		player.X += 1
	case "q":
		player.X -= 1
	}

	return nil, nil
}

func (s *Server) BroadcastGameState() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, client := range s.Game.players {
		state, err := utils.EncodePlayer(*client)
		if err != nil {
			fmt.Printf("Error encoding player state for %s: %v\n", client.Addr, err)
			continue
		}

		for _, otherClient := range s.Game.players {
			_, err := s.conn.WriteToUDP(state, otherClient.Addr)
			if err != nil {
				fmt.Printf("Error sending to %s: %v\n", otherClient.Addr, err)
			} else {
				fmt.Printf("Sent state update to %s\n", otherClient.Addr)
			}
		}
	}
}

func (s *Server) StartTickLoop() {
	ticker := time.NewTicker(time.Second / tickRate)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.BroadcastGameState()
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

	go server.Read()
	server.StartTickLoop()
}
