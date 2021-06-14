package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var WIDTH int = 160
var HEIGHT int = 120
var globalTick int = 0

type Game struct {
	world  *World
	pixels []byte
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Draw(screen *ebiten.Image) {
	// h + left click: heavy- weight spaceship (HWSS)
	if ebiten.IsKeyPressed(ebiten.KeyH) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.world.heavyWeightSpaceship(x, y)
	}

	// r: reset world
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		log.Print("reseting world")
		g.world.reset()
	}

	// p: populate world
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		log.Print("populating world")
		g.world.randPopulate()
	}

	// update tps, up or down
	currTPS := int(math.Round(ebiten.CurrentTPS()))
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ebiten.SetMaxTPS(currTPS + 1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) && currTPS > 0 {
		ebiten.SetMaxTPS(currTPS - 1)
	}

	// draw
	g.pixels = g.world.toPixel(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func (g *Game) Update() error {
	globalTick++
	alive, dead := g.world.update()
	aliveRatio := (alive / (alive + dead)) * 100
	ebiten.SetWindowTitle(
		fmt.Sprintf("tps: %d | gen: %d | alive: %d | dead: %d | alive ratio: %.2f%%",
			int(math.Round(ebiten.CurrentTPS())),
			globalTick,
			int(alive),
			int(dead),
			aliveRatio,
		))

	return nil
}

func initMtx(width, height int) [][]Cell {
	m := make([][]Cell, width)
	for i := range m {
		m[i] = make([]Cell, height)
	}
	return m
}

func main() {
	g := &Game{
		world: &World{
			width:  WIDTH,
			height: HEIGHT,
			mtx:    initMtx(WIDTH, HEIGHT),
		},
		pixels: make([]byte, WIDTH*HEIGHT*4),
	}
	ebiten.SetMaxTPS(1)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
