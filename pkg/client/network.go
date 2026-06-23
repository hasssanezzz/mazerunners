package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Handler func(*config.Message)

type Network interface {
	Init(*config.UserInfo) (*config.InitResponse, error)
	PublishEvent(*config.Message) error
	Subscribe(Handler) error
}

type UDPNetwork struct {
	serverAddress *net.UDPAddr
	conn          *net.UDPConn
	localAddr     *net.UDPAddr

	once sync.Once
}

var _ Network = (*UDPNetwork)(nil)

func NewUDPNetwork(addr *net.UDPAddr) *UDPNetwork {
	n := &UDPNetwork{
		serverAddress: addr,
	}

	return n
}

func (n *UDPNetwork) Init(info *config.UserInfo) (*config.InitResponse, error) {
	conn, err := net.DialUDP("udp", nil, n.serverAddress)
	if err != nil {
		return nil, err
	}
	n.conn = conn

	n.localAddr = conn.LocalAddr().(*net.UDPAddr)
	info.Addr = n.localAddr

	m := &config.Message{
		From:    n.localAddr.String(),
		Event:   config.EventPlayerInit,
		Payload: info,
	}

	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	if _, err := n.conn.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	rbuf := make([]byte, 1024*64)
	nr, err := n.conn.Read(rbuf)
	if err != nil {
		return nil, err
	}
	var response config.InitResponse
	if err := gob.NewDecoder(bytes.NewBuffer(rbuf[:nr])).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (n *UDPNetwork) PublishEvent(m *config.Message) error {
	m.From = n.localAddr.String()

	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return err
	}
	_, err := n.conn.Write(buf.Bytes())
	return err
}

func (n *UDPNetwork) Subscribe(handler Handler) error {
	if n.conn == nil {
		return fmt.Errorf("failed to subscribe, consumer must init first")
	}
	n.once.Do(func() {
		go n.subscribe(handler)
	})
	return nil
}

func (n *UDPNetwork) subscribe(handler Handler) {
	buf := make([]byte, 1024*64)
	for {
		nr, err := n.conn.Read(buf)
		if err != nil {
			log.Println("can't read message:", err)
			continue
		}
		var msg config.Message
		if err := gob.NewDecoder(bytes.NewBuffer(buf[:nr])).Decode(&msg); err != nil {
			log.Println("can't decode message:", err)
			continue
		}
		handler(&msg)
	}
}
