package utils

import (
	"UdpServer/module"
	"bytes"
	"encoding/gob"
)

func DecodePlayer(buffer []byte) (module.Player, error) {
	var player module.Player
	buf := bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&player)
	if err != nil {
		return player, err
	}
	return player, nil
}
