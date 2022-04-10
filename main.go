package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// startX is the initial x position of the falling asteroid
	startX = 200

	// bottomThreshold is the line where image should disappear
	bottomThreshold = 300

	asteroidWidth           = 120
	canonWidth, canonHeight = 48, 32
	shotWidth, shotHeight   = 10, 100

	canonY = screenHeight - canonHeight

	nbFramesBeforeRespawn = 30

	// welcomeMessage is intended to be used in DebugPrint function
	welcomeMessage = `
Press Enter to create an asteroid
Press Escape to quit
Use arrows (or hjkl) to move

Press x to show hitboxes
`
	gameOverMessage     = "Sorry, You lose! Press Enter to restart or Escape to quit"
	gameFinishedMessage = "Wow! You rock at this! Press Enter to restart or Escape to quit"
)

type gameMode int

const (
	notStarted gameMode = iota
	started
	gameOver
	gameFinished
)

var (
	spaceImg    *ebiten.Image
	asteroidImg *ebiten.Image
	canonImg    *ebiten.Image
	shotImg     *ebiten.Image

	redColor    = color.RGBA{0xff, 0, 0, 0xff}
	blueColor   = color.RGBA{0, 0, 0xdd, 0xff}
	greenColor  = color.RGBA{0, 0xff, 0, 0xaa}
	yellowColor = color.RGBA{0xff, 0xff, 0, 0xff}
)

type Game struct {
	mode gameMode

	fallingX, fallingY int
	failedCount        int
	shotCount          int

	canonX  int
	shootX  int
	shootY  int
	isShoot bool

	showHitboxes bool
	isFullScreen bool

	countBeforeRespawn int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Escape pressed. Bye bye...")
	}

	if g.mode != started && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.mode = started
		g.failedCount = 0
		g.shotCount = 0
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyH) {
		g.canonX -= 3
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL) {
		g.canonX += 3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.showHitboxes = !g.showHitboxes
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.isFullScreen = !g.isFullScreen
	}

	if !g.isShoot && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.isShoot = true
		g.shootX = g.canonX - shotWidth/2 + canonWidth/2
	}

	if g.isShoot {
		g.shootY -= 6
		if g.shootY <= 0 {
			g.isShoot = false
			g.shootY = canonY
		}
	}

	if g.fallingY < bottomThreshold {
		g.fallingY += int(math.Ceil(float64(g.shotCount)/3)) + 1
	} else {
		g.fallingY = 0
		g.fallingX = rand.Intn(screenWidth - asteroidWidth/2)

		g.failedCount++
	}

	// maybe use a goroutine for this
	if g.countBeforeRespawn > 0 {
		g.countBeforeRespawn++
		g.fallingY = 0
	}

	if g.countBeforeRespawn >= nbFramesBeforeRespawn {
		g.countBeforeRespawn = 0
	}

	if g.hit() {
		g.fallingY = 0
		g.fallingX = rand.Intn(screenWidth - asteroidWidth/2)

		g.shotCount++

		g.shootY = screenHeight
		g.isShoot = false
		g.countBeforeRespawn++
	}

	if g.failedCount >= 10 {
		g.mode = gameOver
	} else if g.shotCount >= 10 {
		g.mode = gameFinished
	}

	return nil
}

func overlap(a, b, wa, wb int) bool {
	return a < b+wb && b < a+wa
}

func (g *Game) hit() bool {
	// hack to have a hit that feels more precise than the hitbox
	gapToReachRealAsteroid := 10

	shootRect := image.Rect(g.shootX, g.shootY, g.shootX+shotWidth, g.shootY+shotHeight)
	fallingRect := image.Rect(g.fallingX, g.fallingY, g.fallingX+asteroidWidth, g.fallingY+asteroidWidth-gapToReachRealAsteroid)
	return fallingRect.Overlaps(shootRect)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case notStarted:
		ebitenutil.DebugPrint(screen, welcomeMessage)
		return
	case gameOver:
		ebitenutil.DebugPrint(screen, gameOverMessage)
		return
	case gameFinished:
		ebitenutil.DebugPrint(screen, gameFinishedMessage)
		return
	}

	screen.DrawImage(spaceImg, nil)
	if g.isShoot {
		g.drawShoot(screen)
	}
	g.drawCanon(screen)

	if g.countBeforeRespawn == 0 {
		g.drawAsteroid(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("failed: %d\nshot: %d", g.failedCount, g.shotCount))
}

func (g *Game) drawCanon(screen *ebiten.Image) {
	// img := ebiten.NewImage(canonWidth, canonHeight)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(g.canonX), canonY)

	// img.Fill(greenColor)
	screen.DrawImage(canonImg, op)
}

func (g *Game) drawShoot(screen *ebiten.Image) {
	// img := ebiten.NewImage(shotWidth, shotHeight)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(g.shootX), float64(g.shootY))

	// img.Fill(yellowColor)
	screen.DrawImage(shotImg, op)
}

func (g *Game) drawAsteroid(screen *ebiten.Image) {
	if g.fallingY < bottomThreshold {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.fallingX), float64(g.fallingY))
		if g.showHitboxes {
			imgHB := ebiten.NewImageFromImage(asteroidImg)

			if g.hit() {
				imgHB.Fill(greenColor)
			} else {
				imgHB.Fill(blueColor)
			}

			screen.DrawImage(imgHB, op)
		}

		screen.DrawImage(asteroidImg, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebiten.SetFullscreen(g.isFullScreen)
	return screenWidth, screenHeight
}

func main() {
	var err error

	spaceImg, _, err = ebitenutil.NewImageFromFile("space.png")
	asteroidImg, _, err = ebitenutil.NewImageFromFile("asteroid.png")
	canonImg, _, err = ebitenutil.NewImageFromFile("canon.png")
	shotImg, _, err = ebitenutil.NewImageFromFile("shot.png")

	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, Shoot WoRlD!")

	game := &Game{
		fallingX: startX,
		canonX:   (screenWidth - canonWidth) / 2,
		shootY:   canonY,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
