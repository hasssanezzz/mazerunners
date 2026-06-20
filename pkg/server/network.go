package server

import "github.com/hasssanezzz/mazerunners/pkg/config"

type Message struct {
	From      string
	EventType config.Event
	Payload   any
}

type Network interface {
	Listen(addr string) error
	Subscribe(handler func(Message)) error
	Broadcast(Message) error
}
