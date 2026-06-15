package main

import (
	"fmt"
	"iter"
	"math/rand"
)

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

type Cell struct {
	Point
	Kind CellKind
	Row  int
	Col  int
}

type Map struct {
	matrix [][]CellKind
	cfg    *Config
}

func NewMap(size int, cfg *Config) *Map {
	m := &Map{
		matrix: make([][]CellKind, size),
		cfg:    cfg,
	}

	for i := range size {
		m.matrix[i] = make([]CellKind, size)
	}

	return m
}

func (m *Map) FillRandom() {
	for cell := range m.Items() {
		if rand.Intn(10) == 5 {
			m.matrix[cell.Row][cell.Col] = CellWall
		}
	}
}

func (m *Map) FillBorders() {
	n := len(m.matrix)
	for row := range n {
		for col := range n {
			if row == 0 || row == n-1 || col == 0 || col == n-1 {
				m.matrix[row][col] = CellWall
			}
		}
	}
}

func (m *Map) Items() iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for row, line := range m.matrix {
			for col, kind := range line {
				if !yield(Cell{Kind: kind, Row: row, Col: col, Point: Point{X: col * m.cfg.CellSize, Y: row * m.cfg.CellSize}}) {
					return
				}
			}
		}
	}
}

type Camera struct {
	Point
	cfg *Config
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
