package utils

import (
	"UdpServer/module"
	"bytes"
	"encoding/gob"
)

func EncodePlayer(player module.Player) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(player)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
