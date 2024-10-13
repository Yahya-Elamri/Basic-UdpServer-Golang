package main

import (
	"UdpServer/module"
	"UdpServer/utils"
	"fmt"
	"net"
	"sync"
)

const (
	addr       = "127.0.0.1:3000"
	bufferSize = 1024
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
		s.Write(message, remoteAddr)
	}
}

func (s *Server) HandleClients(remoteAddr *net.UDPAddr, message string) {
	if message == "connect" {
		s.mutex.Lock()
		if _, ok := s.Game.players[remoteAddr.String()]; !ok {
			PlayerId := len(s.Game.players) + 1
			s.Game.players[remoteAddr.String()] = &module.Player{ID: PlayerId, Addr: remoteAddr, X: 0, Y: 0}
			fmt.Printf("New client connected: %s (ID: %d)\n", remoteAddr, PlayerId)
		}
		s.mutex.Unlock()
	} else if message == "quit" {
		s.mutex.Lock()
		delete(s.Game.players, remoteAddr.String())
		fmt.Printf("client disconnected: %s \n", remoteAddr)
		s.mutex.Unlock()
	} else {
		s.mutex.Lock()
		_, exists := s.Game.players[remoteAddr.String()]
		s.mutex.Unlock()

		if !exists {
			fmt.Printf("Received game data from unregistered client %s\n", remoteAddr)
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
		player.Y = player.Y + 1
	case "s":
		player.Y = player.Y - 1
	case "d":
		player.X = player.X + 1
	case "q":
		player.X = player.X - 1
	}

	return utils.EncodePlayer(*player)
}

func (s *Server) Write(message string, sender *net.UDPAddr) {
	buffer, _ := s.handleGameLogic(message, sender)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, client := range s.Game.players {
		if client.Addr.String() != sender.String() {
			_, err := s.conn.WriteToUDP(buffer, client.Addr)
			if err != nil {
				fmt.Printf("Error sending to %s: %v\n", client.Addr, err)
			} else {
				fmt.Printf("Sent message to %s: %s\n", client.Addr, message)
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
