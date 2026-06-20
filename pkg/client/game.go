package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Game struct {
	cfg    *config.Config
	world  *config.Map
	camera *config.Camera
	player *Player
	debug  bool
}

func NewGame(cfg *config.Config) *Game {
	world := config.NewMap(cfg.MapSize, cfg)
	world.FillRandom()
	world.FillBorders()

	start := config.Point{X: 1, Y: 1}
	camera := config.NewCamera(start.ToMapCell(cfg.CellSize), cfg)
	player := NewPlayer(start, world, camera, cfg)

	return &Game{
		cfg:    cfg,
		world:  world,
		camera: camera,
		player: player,
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	for cell := range g.world.Items() {
		if !g.camera.IsVisible(cell.Point) {
			continue
		}
		if cell.Kind == config.CellWall {
			pos := g.camera.ToScreen(cell.Point)
			vector.FillRect(screen, float32(pos.X), float32(pos.Y), float32(g.cfg.CellSize), float32(g.cfg.CellSize), config.ColorGrey, false)
		}
	}
}

func (g *Game) Update() error {
	g.player.Update()

	if g.debug {
		if ebiten.IsKeyPressed(ebiten.KeyLeftBracket) {
			if g.cfg.CellSize > 10 {
				g.cfg.CellSize--
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRightBracket) {
			g.cfg.CellSize++
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(config.ColorBG)
	g.drawWorld(screen)
	g.player.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.cfg.WindowSize, g.cfg.WindowSize
}
