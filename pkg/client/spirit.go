package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
