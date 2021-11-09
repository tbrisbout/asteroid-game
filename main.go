package main

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	img       *ebiten.Image
	playerImg *ebiten.Image

	redColor   = color.RGBA{0xff, 0, 0, 0xff}
	blueColor  = color.RGBA{0, 0, 0xdd, 0xff}
	greenColor = color.RGBA{0, 0xff, 0, 0xaa}
)

const (
	// welcomeMessage is intended to be used in DebugPrint function
	welcomeMessage = `
Press Enter to create a Gopher
Press Escape to quit
Use arrows (or hjkl) to move

Press x to show hitboxes
`

	// startX is the initial x position of the falling gopher
	startX = 200

	// bottomThreshold is the line where image should disappear
	bottomThreshold = 300

	// imgWidth is the full width (including transparent) of the gopher image
	imgWidth = 240
)

type Game struct {
	fallingX, fallingY int
	failedCount        int

	gopherCount  int
	x, y         int
	mousePressed bool

	showHitboxes bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Escape pressed. Bye bye...")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.gopherCount < 10 {
		g.gopherCount++
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyH) {
		g.x -= 3
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL) {
		g.x += 3
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyK) {
		g.y -= 3
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.y += 3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.showHitboxes = !g.showHitboxes
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mousePressed = true
	} else {
		g.mousePressed = false
	}

	if g.fallingY < bottomThreshold {
		g.fallingY++
	} else {
		g.fallingY = 0
		g.failedCount++
	}

	return nil
}

func overlap(a, b, wa, wb int) bool {
	return a < b+wb && b < a+wa
}

func (g *Game) hit() bool {
	return overlap(g.x, g.fallingX, imgWidth, imgWidth) && overlap(g.y, g.fallingY, imgWidth, imgWidth)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.gopherCount == 0 {
		ebitenutil.DebugPrint(screen, welcomeMessage)
		return
	}

	// draw player gopher
	for i := 0; i < g.gopherCount; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(50*i+g.x), float64(g.y))

		if g.mousePressed {
			op.GeoM.Scale(1.5, 1)
		}

		if g.showHitboxes {
			playerImgHB := ebiten.NewImageFromImage(img)
			playerImgHB.Fill(redColor)
			screen.DrawImage(playerImgHB, op)
		}

		screen.DrawImage(playerImg, op)
	}

	// draw falling gopher
	if g.fallingY < bottomThreshold {
		fallingOp := &ebiten.DrawImageOptions{}
		fallingOp.GeoM.Translate(float64(g.fallingX), float64(g.fallingY))

		if g.hit() {
			if g.showHitboxes {
				imgHB := ebiten.NewImageFromImage(img)
				imgHB.Fill(greenColor)
				screen.DrawImage(imgHB, fallingOp)
			}
		} else {
			if g.showHitboxes {
				imgHB := ebiten.NewImageFromImage(img)
				imgHB.Fill(blueColor)
				screen.DrawImage(imgHB, fallingOp)
			}
		}
		screen.DrawImage(img, fallingOp)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("failed: %d", g.failedCount))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	playerImg, _, err = ebitenutil.NewImageFromFile("gopher.png")

	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Hitbox WoRlD!")
	if err := ebiten.RunGame(&Game{fallingX: startX}); err != nil {
		log.Fatal(err)
	}
}
