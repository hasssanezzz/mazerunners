package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

type Spirit interface {
	HandleEvent(config.Event)
	Update()
	Draw(*ebiten.Image)
}

type Player struct {
	Position  config.Point
	Direction config.Direction
	Camera    *config.Camera
	World     *config.Map

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

func (p *Player) HandleEvent(e config.Event) {}

const (
	keyRepeatDelay = 18 // ticks before repeat starts (~0.3s at 60fps)
	keyRepeatRate  = 6  // ticks between repeated moves while held
)

func isKeyActive(key ebiten.Key) bool {
	d := inpututil.KeyPressDuration(key)
	return d == 1 || (d > keyRepeatDelay && (d-keyRepeatDelay)%keyRepeatRate == 0)
}

func (p *Player) Update() {
	if isKeyActive(ebiten.KeyRight) {
		p.Direction = config.DirectionRight
		if p.canMoveTo(p.Position.X+1, p.Position.Y) {
			p.Position.X += 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
	}

	if isKeyActive(ebiten.KeyLeft) {
		p.Direction = config.DirectionLeft
		if p.canMoveTo(p.Position.X-1, p.Position.Y) {
			p.Position.X -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
	}

	if isKeyActive(ebiten.KeyUp) {
		p.Direction = config.DirectionUp
		if p.canMoveTo(p.Position.X, p.Position.Y-1) {
			p.Position.Y -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
	}

	if isKeyActive(ebiten.KeyDown) {
		p.Direction = config.DirectionDown
		if p.canMoveTo(p.Position.X, p.Position.Y+1) {
			p.Position.Y += 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
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
		8,
		config.ColorRed,
		false,
	)
}

func (p *Player) canMoveTo(x, y int) bool {
	return p.World.Matrix[y][x] != config.CellWall
}
