package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

const (
	keyRepeatDelay = 18 // ticks before repeat starts (~0.3s at 60fps)
	keyRepeatRate  = 6  // ticks between repeated moves while held
)

func isKeyActive(key ebiten.Key) bool {
	d := inpututil.KeyPressDuration(key)
	return d == 1 || (d > keyRepeatDelay && (d-keyRepeatDelay)%keyRepeatRate == 0)
}

type Spirit interface {
	HandleEvent(config.Event)
	Update()
	Draw(*ebiten.Image)
	Alive() bool
}

const arrowSpeed = 2 // ticks between each tile move

type Arrow struct {
	position   config.Point
	direction  config.Direction
	alive      bool
	ticksAlive int

	world  *config.Map
	camera *config.Camera
	cfg    *config.Config
}

var _ Spirit = (*Arrow)(nil)

func (a *Arrow) HandleEvent(e config.Event) {}
func (a *Arrow) Alive() bool                { return a.alive }

func (a *Arrow) Update() {
	a.ticksAlive++
	if !a.alive || a.ticksAlive%arrowSpeed != 0 {
		return
	}

	nx, ny := a.position.X, a.position.Y
	switch a.direction {
	case config.DirectionUp:
		ny--
	case config.DirectionDown:
		ny++
	case config.DirectionLeft:
		nx--
	case config.DirectionRight:
		nx++
	}

	if !a.world.CanMoveTo(nx, ny) {
		a.alive = false
		return
	}

	a.position.X = nx
	a.position.Y = ny
}

func (a *Arrow) Draw(screen *ebiten.Image) {
	if !a.alive {
		return
	}
	pos := a.camera.ToScreen(a.position.ToMapCell(a.cfg.CellSize))
	size := float32(a.cfg.CellSize)
	vector.FillRect(screen, float32(pos.X), float32(pos.Y), size, size, config.ColorRed, false)
}

type Player struct {
	Position  config.Point
	Direction config.Direction
	Camera    *config.Camera
	World     *config.Map
	Coins     int
	Spawn     func(Spirit)

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
	}

	if isKeyActive(ebiten.KeyLeft) {
		p.Direction = config.DirectionLeft
		if p.World.CanMoveTo(p.Position.X-1, p.Position.Y) {
			p.Position.X -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
	}

	if isKeyActive(ebiten.KeyUp) {
		p.Direction = config.DirectionUp
		if p.World.CanMoveTo(p.Position.X, p.Position.Y-1) {
			p.Position.Y -= 1
			p.Camera.CenterOn(p.Position.ToMapCell(p.cfg.CellSize))
		}
	}

	if isKeyActive(ebiten.KeyDown) {
		p.Direction = config.DirectionDown
		if p.World.CanMoveTo(p.Position.X, p.Position.Y+1) {
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
		4,
		config.ColorRed,
		false,
	)
}
