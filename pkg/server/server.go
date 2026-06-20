package server

import "net"

type Client struct {
	address net.Addr
	name    string
}

type Server struct {
	addr    *net.UDPAddr
	clients []Client
}

func NewServer(addr *net.UDPAddr) *Server {
	s := &Server{
		addr: addr,
	}

	return s
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}

	_ = conn

	return nil
}
