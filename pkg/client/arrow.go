package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hasssanezzz/mazerunners/pkg/config"
)

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
	if !a.alive {
		return
	}
	a.ticksAlive++
	if a.ticksAlive%arrowSpeed != 0 {
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
