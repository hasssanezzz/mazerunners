package config

import (
	"encoding/gob"
	"net"
)

func init() {
	gob.Register(&UserInfo{})
	gob.Register(&InitResponse{})
}

type Message struct {
	From    string
	Event   Event
	Payload any
}

type UserInfo struct {
	Name   string
	Secret string
	Addr   *net.UDPAddr
}

type InitResponse struct {
	World       *Map
	PlayerCount int
}
