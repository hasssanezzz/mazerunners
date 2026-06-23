package client

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Game struct {
	cfg      *config.Config
	world    *config.Map
	camera   *config.Camera
	network  Network
	userInfo *config.UserInfo

	spirits []Spirit
	player  *Player
	debug   bool
}

func NewGame(cfg *config.Config, network Network, userInfo *config.UserInfo) (*Game, error) {
	response, err := network.Init(userInfo)
	if err != nil {
		return nil, fmt.Errorf("new game failed to init network: %v", err)
	}

	if response.World == nil {
		return nil, fmt.Errorf("server sent nil world")
	}
	response.World.Cfg = cfg

	start := config.Point{X: 1, Y: 1}
	camera := config.NewCamera(start.ToMapCell(cfg.CellSize), cfg)
	player := NewPlayer(start, response.World, camera, cfg)

	g := &Game{
		cfg:      cfg,
		world:    response.World,
		camera:   camera,
		player:   player,
		network:  network,
		userInfo: userInfo,
		spirits:  make([]Spirit, 0, 100),
		debug:    true,
	}

	player.Spawn = func(s Spirit) {
		g.spirits = append(g.spirits, s)
	}

	// TODO: react when other players join

	player.OnStateChange = func(p config.Point, d config.Direction) {
		if err := network.PublishEvent(&config.Message{
			Event: config.EventPlayerStateChange,
			Payload: &config.PlayerStateChangePayload{
				Point: p,
				Dir:   d,
			},
		}); err != nil {
			log.Println("failed to public event:", err)
		}
	}

	return g, nil
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
