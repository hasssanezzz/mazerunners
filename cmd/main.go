package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasssanezzz/mazerunners/pkg/client"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

func runClient() {
	os.Getenv("MZRN")

	// playerName := os.Getenv("MZRN_PLAYER_NAME")
	// playerPassword := os.Getenv("MZRN_PLAYER_PASSWORD")
	// serverAddress := os.Getenv("MZRN_SERVER_ADDRESS")

	cfg := &config.Config{
		WindowSize: 800,
		CellSize:   25,
		MapSize:    50,
	}

	ebiten.SetWindowSize(cfg.WindowSize, cfg.WindowSize)
	ebiten.SetWindowTitle("Maze Runners")

	if err := ebiten.RunGame(client.NewGame(cfg)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		return
	}

	runClient()
}
