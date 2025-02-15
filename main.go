package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8

	jumpVelocity = -10
	gravity      = 0.5
)

var (
	RunnerImage *ebiten.Image
)

type Game struct {
	count       int
	x           float64
	y           float64
	vy          float64
	facingRight bool
}

func (g *Game) Update() error {
	// Handle jumping
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.y == 0 {
		g.vy = jumpVelocity
	}

	// Handle walking
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		fmt.Println("Right arrow key pressed")
		g.x += 2  // Move right
		g.count++ // Increment count to update animation frame
		g.facingRight = true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		fmt.Println("Left arrow key pressed")
		g.x -= 2  // Move left
		g.count++ // Increment count to update animation frame
		g.facingRight = false
	}

	// Update vertical position and velocity
	g.vy += gravity
	g.y += g.vy

	// Prevent the character from falling below the ground
	if g.y > 0 {
		g.y = 0
		g.vy = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Scale(2, 2)
	if !g.facingRight {
		op.GeoM.Scale(-1, 1)                         // Mirror the image horizontally
		op.GeoM.Translate(-float64(frameWidth)/2, 0) // Adjust position after mirroring
	}
	op.GeoM.Translate(g.x, screenHeight/2+g.y)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(RunnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Load the image from the "img" folder.
	imgPath := filepath.Join("img", "runner.png")
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	RunnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
