package module

import "net"

type (
	Player struct {
		ID   int
		X, Y int
		Addr *net.UDPAddr
	}
)
