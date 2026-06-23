package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/hasssanezzz/mazerunners/pkg/config"
	"go.uber.org/zap"
)

var logger *zap.Logger = func() *zap.Logger {
	l, _ := zap.NewDevelopment()
	return l
}()

type GameState struct {
	World *config.Map
}

type Room struct {
	ID      int
	State   *GameState
	Clients map[string]*config.UserInfo
}

func (r *Room) addClient(client *config.UserInfo) {
	r.Clients[client.Addr.String()] = client
}

// TODO: we have zero thread safty
type Server struct {
	addr      *net.UDPAddr
	running   bool
	blacklist map[string]struct{}
	rooms     map[int]*Room
	conn      *net.UDPConn
	cfg       *config.Config
}

func NewServer(addr *net.UDPAddr, cfg *config.Config) *Server {
	s := &Server{
		addr:      addr,
		running:   true,
		blacklist: make(map[string]struct{}),
		rooms:     map[int]*Room{},
		cfg:       cfg,
	}

	return s
}

func (s *Server) resolveRoom(roomID int) (*Room, bool) {
	if room, ok := s.rooms[roomID]; ok {
		return room, false
	}

	world := config.NewMap(s.cfg)
	world.FillRandom()
	world.FillBorders()

	r := &Room{
		ID: roomID,
		State: &GameState{
			World: world,
		},
		Clients: map[string]*config.UserInfo{},
	}
	s.rooms[roomID] = r

	return r, true
}

func (s *Server) findClientInRooms(addr string) (*config.UserInfo, bool) {
	for _, room := range s.rooms {
		for _, client := range room.Clients {
			if client.Addr.String() == addr {
				return client, true
			}
		}
	}

	return nil, false
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}
	s.conn = conn

	for s.running {
		buf := make([]byte, 1024*4)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("error happend while reading", zap.Error(err))
			continue
		}
		buf = buf[:n]

		if _, ok := s.blacklist[addr.String()]; ok {
			continue // ignore messages coming from a blacklisted address
		}

		var m config.Message
		if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&m); err != nil {
			logger.Error("error happend while decoding", zap.Error(err))
			continue
		}

		go func() {
			if err := s.handleMessage(addr, &m); err != nil {
				logger.Error("error happend while handling message", zap.Error(err))
			}
		}()
	}

	return nil
}

func (s *Server) broadcast(roomID int, m *config.Message, senderAddr *net.UDPAddr) error {
	room, ok := s.rooms[roomID]
	if !ok {
		return ErrRoomNotFound{roomID}
	}

	for _, client := range room.Clients {
		if client.Addr.String() == senderAddr.String() {
			continue
		}

		if err := s.sendMessage(client.Addr, m); err != nil {
			logger.Warn("failed to send broadcast message", fieldPlayerName(client.Name), fieldAddr(client.Addr), zap.Error(err))
		}
	}

	return nil
}

func (s *Server) sendMessage(addr *net.UDPAddr, m *config.Message) error {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return err
	}

	_, err := s.conn.WriteToUDP(buf.Bytes(), addr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendBytes(addr *net.UDPAddr, data []byte) error {
	_, err := s.conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleMessage(addr *net.UDPAddr, m *config.Message) error {
	if m.From != addr.String() {
		logger.Warn("client addr message mismatch, backlisting the client", fieldAddr(addr), zap.String("msgAddr", m.From))
		s.blacklist[addr.String()] = struct{}{}
		return nil
	}

	client, ok := s.findClientInRooms(addr.String())
	if !ok && m.Event != config.EventPlayerInit {
		logger.Warn("someone is sending an event before joining :/", fieldAddr(addr))
		return nil
	}

	switch m.Event {
	case config.EventPlayerInit:
		return s.handlePlayerInit(addr, m)
	case config.EventPlayerStateChange:
		return s.handlePlayStateChange(client, addr, m)
	}

	return nil
}

func (s *Server) handlePlayerInit(addr *net.UDPAddr, m *config.Message) error {
	info, ok := m.Payload.(*config.UserInfo)
	if !ok {
		return fmt.Errorf("message payload type mismatch")
	}

	room, created := s.resolveRoom(info.Room)
	if created {
		logger.Info("new room created", fieldRoom(info.Room))
	}
	room.addClient(info)

	response := config.InitResponse{
		World:       room.State.World,
		PlayerCount: len(room.Clients),
	}
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(response); err != nil {
		return err
	}

	if err := s.sendBytes(addr, buf.Bytes()); err != nil {
		return err
	}

	logger.Info("player joined", fieldPlayerName(info.Name), fieldAddr(addr), fieldRoom(info.Room))

	return nil
}

func (s *Server) handlePlayStateChange(client *config.UserInfo, addr *net.UDPAddr, m *config.Message) error {
	state, ok := m.Payload.(*config.PlayerStateChangePayload)
	if !ok {
		return fmt.Errorf("message type payload mismatch")
	}

	fmt.Printf("-- player %s is moving, Point: %s Direction: %d\n", client.Name, state.Point, state.Dir)
	if err := s.broadcast(client.Room, m, addr); err != nil {
		return err
	}

	return nil
}
