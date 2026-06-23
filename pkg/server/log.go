package server

import (
	"net"

	"go.uber.org/zap"
)

func fieldPlayerName(name string) zap.Field {
	return zap.String("playerName", name)
}

func fieldAddr(addr *net.UDPAddr) zap.Field {
	return zap.String("addr", addr.String())
}

func fieldRoom(room int) zap.Field {
	return zap.Int("room", room)
}
