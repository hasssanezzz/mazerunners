package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type GameState struct {
	world *config.Map
}

type Server struct {
	addr      *net.UDPAddr
	clients   map[string]*config.UserInfo
	running   bool
	blacklist map[string]struct{}

	state GameState
}

func NewServer(addr *net.UDPAddr, cfg *config.Config) *Server {
	world := config.NewMap(cfg)
	world.FillRandom()
	world.FillBorders()

	s := &Server{
		addr:      addr,
		clients:   make(map[string]*config.UserInfo),
		running:   true,
		blacklist: make(map[string]struct{}),
		state: GameState{
			world: world,
		},
	}

	return s
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}

	for s.running {
		buf := make([]byte, 1024*4)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("error happend while reading:", err)
			continue
		}
		buf = buf[:n]

		if _, ok := s.blacklist[addr.String()]; ok {
			continue // ignore messages coming from a blacklisted address
		}

		var m config.Message
		if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&m); err != nil {
			log.Println("error happend while decoding:", err)
			continue
		}

		go func() {
			if err := s.handleMessage(conn, addr, &m); err != nil {
				log.Printf("error while handling message from %q: %v\n", addr.String(), err)
			}
		}()
	}

	return nil
}

func (s *Server) handleMessage(conn *net.UDPConn, addr *net.UDPAddr, m *config.Message) error {
	if m.From != addr.String() {
		log.Println("client addr message mismatch, backlisting the client")
	}

	_, ok := s.clients[addr.String()]
	if !ok {
		if m.Event == config.EventPlayerInit {
			info, ok := m.Payload.(*config.UserInfo)
			if !ok {
				return fmt.Errorf("message type payload mismatch")
			}

			s.clients[info.Addr.String()] = info

			response := config.InitResponse{
				World:       s.state.world,
				PlayerCount: len(s.clients),
			}
			buf := bytes.NewBuffer(nil)
			if err := gob.NewEncoder(buf).Encode(response); err != nil {
				return err
			}
			n, err := conn.WriteToUDP(buf.Bytes(), addr)
			if err != nil {
				return err
			}

			if n != buf.Len() {
				panic("n != buf.Len()") // temp debug
			}
		} else {
			log.Println("someone is sending an event before joining :/")
		}
	}

	return nil
}
