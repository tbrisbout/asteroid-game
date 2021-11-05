package main

import (
	"errors"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	gopherCount int
	x, y        int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Escape pressed. Bye bye...")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.gopherCount < 10 {
		g.gopherCount++
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.x--
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.x++
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.y--
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.y++
		return nil
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.gopherCount == 0 {
		ebitenutil.DebugPrint(screen, `
Press Enter to create a Gopher
Press Escape to quit
`)
		return
	}

	for i := 0; i < g.gopherCount; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(50*i+g.x*2), float64(g.y*2))
		screen.DrawImage(img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Matrix WoRlD!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
