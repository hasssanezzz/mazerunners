package config

import (
	"fmt"
	"iter"
	"math/rand"
	"slices"
)

type Config struct {
	WindowSize int
	CellSize   int
	MapSize    int
}

func (c *Config) GridLength() int {
	return c.WindowSize / c.CellSize
}

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{x, y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Point) ToMapCell(cellSize int) Point {
	return Point{
		X: p.X * cellSize,
		Y: p.Y * cellSize,
	}
}

func (p Point) ToScreenCord(cellSize int) Point {
	return Point{
		X: p.X / cellSize,
		Y: p.Y / cellSize,
	}
}

type Cell struct {
	Point
	Kind CellKind
	Row  int
	Col  int
}

type Map struct {
	Matrix [][]CellKind
	Cfg    *Config
}

func NewMap(cfg *Config) *Map {
	size := cfg.MapSize
	m := &Map{
		Matrix: make([][]CellKind, size),
		Cfg:    cfg,
	}

	for i := range size {
		m.Matrix[i] = make([]CellKind, size)
	}

	return m
}

func (m *Map) Get(p Point) (CellKind, bool) {
	if p.X < 0 || p.X >= len(m.Matrix) || p.Y < 0 || p.Y >= len(m.Matrix) {
		return 0, false
	}
	return m.Matrix[p.Y][p.X], true
}

func (m *Map) Set(p Point, cell CellKind) {
	_, ok := m.Get(p)
	if !ok {
		return
	}
	m.Matrix[p.Y][p.X] = cell
}

func (m *Map) FillRandom() {
	for cell := range m.Items() {
		if rand.Intn(10) == 5 {
			if rand.Float32() > 0.5 {
				m.Matrix[cell.Row][cell.Col] = CellWall
			} else {
				m.Matrix[cell.Row][cell.Col] = CellWood
			}
			continue
		}

		if rand.Intn(100) == 5 {
			m.Matrix[cell.Row][cell.Col] = CellCoin
		}
	}
}

func (m *Map) FillBorders() {
	n := len(m.Matrix)
	for row := range n {
		for col := range n {
			if row == 0 || row == n-1 || col == 0 || col == n-1 {
				m.Matrix[row][col] = CellWall
			}
		}
	}
}

func (m *Map) Items() iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for row, line := range m.Matrix {
			for col, kind := range line {
				if !yield(Cell{Kind: kind, Row: row, Col: col, Point: Point{X: col * m.Cfg.CellSize, Y: row * m.Cfg.CellSize}}) {
					return
				}
			}
		}
	}
}

func (m *Map) CanMoveTo(x, y int) bool {
	cell, ok := m.Get(Point{x, y})
	if !ok {
		return false
	}

	blocks := []CellKind{CellWall, CellWood}
	return !slices.Contains(blocks, cell)
}

type Camera struct {
	Point
	cfg *Config
}

func NewCamera(pos Point, cfg *Config) *Camera {
	return &Camera{Point: pos, cfg: cfg}
}

func (c *Camera) CenterOn(p Point) {
	c.Point = p
}

func (c *Camera) ToScreen(p Point) Point {
	return Point{X: p.X - (c.X - c.cfg.WindowSize/2),
		Y: p.Y - (c.Y - c.cfg.WindowSize/2)}
}

func (c *Camera) IsVisible(p Point) bool {
	half := c.cfg.WindowSize / 2
	return p.X+c.cfg.CellSize > c.X-half && p.X < c.X+half &&
		p.Y+c.cfg.CellSize > c.Y-half && p.Y < c.Y+half
}
