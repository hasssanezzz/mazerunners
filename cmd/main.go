package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Config struct {
	WindowSize int
	CellSize   int
	MapSize    int
}

func (c *Config) GridLength() int {
	return c.WindowSize / c.CellSize
}

type Game struct {
	cfg    *Config
	world  *Map
	camera *Camera
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for i := range g.cfg.GridLength() {
		pxl_cord := i * g.cfg.CellSize
		vector.StrokeLine(screen, float32(pxl_cord), 0, float32(pxl_cord), float32(g.cfg.WindowSize), 1, ColorGrey, false)
		vector.StrokeLine(screen, 0, float32(pxl_cord), float32(g.cfg.WindowSize), float32(pxl_cord), 1, ColorGrey, false)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	for cell := range g.world.Items() {
		if !g.camera.IsVisible(cell.Point) {
			continue
		}

		// face := &text.GoTextFace{
		// 	Source: mplusFaceSource,
		// 	Size:   10,
		// }
		// w, h := text.Measure(cell.Point.String(), face, 0)
		// op := &text.DrawOptions{}
		// pos := g.camera.ToScreen(Point{(cell.X) + ((g.cfg.CellSize)-int(w))/2, (cell.Y) + ((g.cfg.CellSize)-int(h))/2})
		// op.GeoM.Translate(float64(pos.X), float64(pos.Y))
		// op.ColorScale.ScaleWithColor(ColorBlack)
		// text.Draw(screen, cell.Point.String(), face, op)

		if cell.Kind == CellWall {
			pos := g.camera.ToScreen(cell.Point)
			vector.FillRect(
				screen,
				float32(pos.X),
				float32(pos.Y),
				float32(g.cfg.CellSize),
				float32(g.cfg.CellSize),
				ColorGrey,
				false,
			)
		}
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.CenterOn(Point{g.camera.X + 10, g.camera.Y})
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.CenterOn(Point{g.camera.X - 10, g.camera.Y})
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.CenterOn(Point{g.camera.X, g.camera.Y - 10})
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.CenterOn(Point{g.camera.X, g.camera.Y + 10})
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeftBracket) {
		if g.cfg.CellSize-1 > 10 {
			g.cfg.CellSize -= 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRightBracket) {
		g.cfg.CellSize += 1
	}

	return nil
}

var mplusFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBG)

	// g.drawGrid(screen)
	g.drawWorld(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.WindowSize, g.cfg.WindowSize
}

func main() {
	cfg := &Config{
		WindowSize: 500,
		CellSize:   50,
		MapSize:    100,
	}

	world := NewMap(cfg.MapSize, cfg)
	world.FillRandom()
	world.FillBorders()

	ebiten.SetWindowSize(cfg.WindowSize, cfg.WindowSize)
	ebiten.SetWindowTitle("Game dev 101")

	if err := ebiten.RunGame(&Game{
		cfg:   cfg,
		world: world,
		camera: &Camera{
			Point: Point{0, 0},
			cfg:   cfg,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
