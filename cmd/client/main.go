package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hasssanezzz/mazerunners/pkg/client"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

func main() {
	cfg := &config.Config{
		WindowSize: 500,
		CellSize:   50,
		MapSize:    100,
	}

	ebiten.SetWindowSize(cfg.WindowSize, cfg.WindowSize)
	ebiten.SetWindowTitle("Maze Runners")

	if err := ebiten.RunGame(client.NewGame(cfg)); err != nil {
		log.Fatal(err)
	}
}
