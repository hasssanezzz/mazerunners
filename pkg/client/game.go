package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Game struct {
	cfg    *config.Config
	world  *config.Map
	camera *config.Camera

	spirits []Spirit
	player  *Player
	debug   bool
}

func NewGame(cfg *config.Config) *Game {
	world := config.NewMap(cfg.MapSize, cfg)
	world.FillRandom()
	world.FillBorders()

	start := config.Point{X: 1, Y: 1}
	camera := config.NewCamera(start.ToMapCell(cfg.CellSize), cfg)
	player := NewPlayer(start, world, camera, cfg)

	g := &Game{
		cfg:     cfg,
		world:   world,
		camera:  camera,
		player:  player,
		spirits: make([]Spirit, 0, 100),
		debug:   true,
	}

	player.Spawn = func(s Spirit) {
		g.spirits = append(g.spirits, s)
	}

	return g
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	for cell := range g.world.Items() {
		if !g.camera.IsVisible(cell.Point) {
			continue
		}

		if cell.Kind == config.CellEmpty {
			continue
		}

		pos := g.camera.ToScreen(cell.Point)
		color := config.ColorGrey

		switch cell.Kind {
		case config.CellWall:
			color = config.ColorGrey
		case config.CellCoin:
			color = config.ColorGold
		case config.CellWood:
			color = config.ColorWood
		}

		vector.FillRect(screen, float32(pos.X), float32(pos.Y), float32(g.cfg.CellSize), float32(g.cfg.CellSize), color, false)
	}
}

func (g *Game) Update() error {
	g.player.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.HandleEvent(config.EventPlayerShoot)
	}

	alive := g.spirits[:0]
	for _, s := range g.spirits {
		s.Update()
		if s.Alive() {
			alive = append(alive, s)
		}
	}
	for i := len(alive); i < len(g.spirits); i++ {
		g.spirits[i] = nil
	}
	g.spirits = alive

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
	for _, s := range g.spirits {
		s.Draw(screen)
	}
	g.player.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.cfg.WindowSize, g.cfg.WindowSize
}
