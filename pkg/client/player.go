package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Player struct {
	Position      config.Point
	Direction     config.Direction
	Camera        *config.Camera
	World         *config.Map
	Coins         int
	Spawn         func(Spirit)
	OnStateChange func(config.Point, config.Direction)

	cfg *config.Config
}

var _ Spirit = (*Player)(nil)

func NewPlayer(pos config.Point, m *config.Map, c *config.Camera, cfg *config.Config) *Player {
	p := &Player{
		Position: pos,
		World:    m,
		Camera:   c,
		cfg:      cfg,
	}

	return p
}

func (p *Player) Alive() bool { return true }

func (p *Player) HandleEvent(e config.Event) {
	switch e {
	case config.EventPlayerShoot:
		if p.Spawn == nil {
			return
		}
		p.Spawn(&Arrow{
			position:  p.Position,
			direction: p.Direction,
			world:     p.World,
			camera:    p.Camera,
			cfg:       p.cfg,
			alive:     true,
		})
	}
}

func (p *Player) Update() {
	if isKeyActive(ebiten.KeyRight) {
		p.Direction = config.DirectionRight
		if p.World.CanMoveTo(p.Position.X+1, p.Position.Y) {
			p.Position.X += 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
		p.OnStateChange(p.Position, p.Direction)
	}

	if isKeyActive(ebiten.KeyLeft) {
		p.Direction = config.DirectionLeft
		if p.World.CanMoveTo(p.Position.X-1, p.Position.Y) {
			p.Position.X -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
		p.OnStateChange(p.Position, p.Direction)
	}

	if isKeyActive(ebiten.KeyUp) {
		p.Direction = config.DirectionUp
		if p.World.CanMoveTo(p.Position.X, p.Position.Y-1) {
			p.Position.Y -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
		p.OnStateChange(p.Position, p.Direction)
	}

	if isKeyActive(ebiten.KeyDown) {
		p.Direction = config.DirectionDown
		if p.World.CanMoveTo(p.Position.X, p.Position.Y+1) {
			p.Position.Y += 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
		p.OnStateChange(p.Position, p.Direction)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	cord := p.Camera.ToScreen(p.Position.ToMapCell(p.cfg.CellSize))
	vector.FillRect(screen, float32(cord.X), float32(cord.Y), float32(p.cfg.CellSize), float32(p.cfg.CellSize), config.ColorPrimary, false)

	cellSize := p.cfg.CellSize

	var start, end config.Point
	switch p.Direction {
	case config.DirectionUp:
		start = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize, p.Position.Y*cellSize))
		end = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize+cellSize, p.Position.Y*cellSize))
	case config.DirectionRight:
		start = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize+cellSize, p.Position.Y*cellSize))
		end = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize+cellSize, p.Position.Y*cellSize+cellSize))
	case config.DirectionDown:
		start = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize, p.Position.Y*cellSize+cellSize))
		end = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize+cellSize, p.Position.Y*cellSize+cellSize))
	case config.DirectionLeft:
		start = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize, p.Position.Y*cellSize))
		end = p.Camera.ToScreen(config.NewPoint(p.Position.X*cellSize, p.Position.Y*cellSize+cellSize))
	}

	vector.StrokeLine(
		screen,
		float32(start.X),
		float32(start.Y),
		float32(end.X),
		float32(end.Y),
		4,
		config.ColorRed,
		false,
	)
}
