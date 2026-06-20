package config

import "image/color"

type CellKind uint8

const (
	CellEmpty CellKind = iota
	CellWall
	CellTrap
	CellCoin
	CellWood
)

var (
	ColorBlack   = color.RGBA{39, 51, 56, 255} // dark slate
	ColorPrimary = color.RGBA{57, 177, 209, 255}

	ColorBG   = color.RGBA{255, 242, 219, 255}
	ColorWood = color.RGBA{157, 102, 56, 255}

	ColorRed  = color.RGBA{245, 73, 39, 255}
	ColorGrey = color.RGBA{70, 80, 88, 255} // wall
	ColorGold = color.RGBA{255, 191, 00, 255}

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
	EventPlayerShoot
	EventPlayerInit
)
