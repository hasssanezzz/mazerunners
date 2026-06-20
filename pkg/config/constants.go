package config

import "image/color"

type CellKind uint8

const (
	CellEmpty CellKind = iota
	CellWall
	CellTrap
	CellCoin
)

var (
	ColorBlack   = color.RGBA{39, 51, 56, 255}   // dark slate
	ColorPrimary = color.RGBA{43, 87, 72, 255}   // forest green (player)
	ColorBG      = color.RGBA{97, 135, 100, 255} // muted green (background)
	ColorRed     = color.RGBA{245, 73, 39, 255}  // trap / pointer
	ColorGrey    = color.RGBA{70, 80, 88, 255}   // wall
	ColorGold    = color.RGBA{255, 191, 00, 255}

	ColorWhite       = color.RGBA{255, 255, 255, 255}
	ColorTransparent = color.RGBA{0, 0, 0, 0}
	ColorPureBlack   = color.RGBA{0, 0, 0, 255}
)

type Direction uint8

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Event uint8

const (
	EventNoOp Event = iota
)
