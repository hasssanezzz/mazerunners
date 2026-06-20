package main

import (
	"log"
	"net"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasssanezzz/mazerunners/pkg/client"
	"github.com/hasssanezzz/mazerunners/pkg/config"
	"github.com/hasssanezzz/mazerunners/pkg/server"
)

func requireEnvVar(name, value, msg string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		if len(value) > 0 {
			return value
		}
		println(msg)
		os.Exit(0)
	}
	return v
}

func runClient() {
	playerName := requireEnvVar("MZRN_PLAYER_NAME", "", "please provide MZRN_PLAYER_NAME")
	playerPassword := requireEnvVar("MZRN_PLAYER_PASSWORD", "", "please provide MZRN_PLAYER_PASSWORD")
	serverAddress := requireEnvVar("MZRN_SERVER_ADDRESS", "0.0.0.0:8000", "please provide MZRN_SERVER_ADDRESS")

	serverUdpAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		panic(err)
	}

	userInfo := &config.UserInfo{
		Name:   playerName,
		Secret: playerPassword,
	}

	network := client.NewUDPNetwork(serverUdpAddr)

	cfg := &config.Config{
		WindowSize: 800,
		CellSize:   25,
		MapSize:    50,
	}

	ebiten.SetWindowSize(cfg.WindowSize, cfg.WindowSize)
	ebiten.SetWindowTitle("Maze Runners")

	game, err := client.NewGame(cfg, network, userInfo)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func runServer() {
	serverAddress := requireEnvVar("MZRN_SERVER_ADDRESS", "0.0.0.0:8000", "")

	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.Config{
		CellSize: 25,
		MapSize:  50,
	}

	s := server.NewServer(addr, cfg)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		println("running server")
		runServer()
		return
	}

	runClient()
}
