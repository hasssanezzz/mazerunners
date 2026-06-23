package config

import (
	"encoding/gob"
	"net"
)

func init() {
	gob.Register(&UserInfo{})
	gob.Register(&InitResponse{})
	gob.Register(&PlayerStateChangePayload{})
	gob.Register(&PlayerShootPayload{})
	gob.Register(&Point{})
}

type PlayerShootPayload struct {
	Point Point
	Dir   Direction
}

type PlayerStateChangePayload struct {
	Point Point
	Dir   Direction
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
	Room   int
}

type InitResponse struct {
	World       *Map
	PlayerCount int
}
