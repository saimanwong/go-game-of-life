package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	width  int
	height int
	mtx    [][]bool
}

func (w *World) heavyWeightSpaceship(x, y int) {
	// Check out of bounds
	if x+6 >= w.width || y+4 >= w.height {
		log.Print("heavy weight spaceship index out of range")
		return
	}
	w.mtx[x+2][y] = true
	w.mtx[x+3][y] = true

	w.mtx[x][y+1] = true
	w.mtx[x+5][y+1] = true

	w.mtx[x+6][y+2] = true

	w.mtx[x][y+3] = true
	w.mtx[x+6][y+3] = true

	w.mtx[x+1][y+4] = true
	w.mtx[x+2][y+4] = true
	w.mtx[x+3][y+4] = true
	w.mtx[x+4][y+4] = true
	w.mtx[x+5][y+4] = true
	w.mtx[x+6][y+4] = true
}

func (w *World) randPopulate() {
	for i := 0; i < 1337; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.mtx[x][y] = true
	}
	globalTick = 0
}

func (w *World) reset() {
	for x := range w.mtx {
		for y := range w.mtx[x] {
			w.mtx[x][y] = false
		}
	}
}

func (w *World) toPixel(p []byte) []byte {
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			tmp := [4]int{}
			i := 4 * x
			j := 4 * y * w.width
			for n := range tmp {
				tmp[n] = i + j + n
			}
			for _, pos := range tmp {
				if w.mtx[x][y] {
					p[pos] = 0xff
					continue
				}
				p[pos] = 0
			}
		}
	}
	return p
}

type point struct {
	x     int
	y     int
	alive bool
}

var direction map[string]point = map[string]point{
	"UP": point{
		x: 0,
		y: -1,
	},
	"UP_RIGHT": point{
		x: 1,
		y: -1,
	},
	"RIGHT": point{
		x: 1,
		y: 0,
	},
	"DOWN_RIGHT": point{
		x: 1,
		y: 1,
	},
	"DOWN": point{
		x: 0,
		y: 1,
	},
	"DOWN_LEFT": point{
		x: -1,
		y: 1,
	},
	"LEFT": point{
		x: -1,
		y: 0,
	},
	"UP_LEFT": point{
		x: -1,
		y: -1,
	},
}

func (w *World) update() {
	next := []point{}
	currAlive := 0
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if w.mtx[x][y] {
				currAlive++
			}
			// Check all direction
			neigh := []point{}
			for _, p := range direction {
				curr := point{
					x: x + p.x,
					y: y + p.y,
				}

				// Check out of bounds
				if curr.x < 0 || curr.x >= w.width {
					continue
				}
				if curr.y < 0 || curr.y >= w.height {
					continue
				}

				curr.alive = w.mtx[curr.x][curr.y]
				neigh = append(neigh, curr)

			}
			dead := 0
			alive := 0
			for _, v := range neigh {
				if v.alive {
					alive++
					continue
				}
				dead++
			}

			// Any live cell with two or three live neighbours survives.
			if w.mtx[x][y] == true && (alive == 2 || alive == 3) {
				continue
			}

			// Any dead cell with three live neighbours becomes a live cell.
			if w.mtx[x][y] == false && alive == 3 {
				next = append(next, point{x: x, y: y, alive: true})
				continue
			}

			// All other live cells die in the next generation. Similarly, all other dead cells stay dead.
			if w.mtx[x][y] == true {
				next = append(next, point{x: x, y: y, alive: false})
			}
		}
	}

	currDead := w.width*w.height - currAlive
	currAliveRatio := (float64(currAlive) / float64(currAlive+currDead)) * 100
	ebiten.SetWindowTitle(
		fmt.Sprintf("gen: %d | alive: %d | dead: %d | alive ratio: %.2f%%",
			globalTick,
			currAlive,
			currDead,
			currAliveRatio,
		))

	for _, v := range next {
		w.mtx[v.x][v.y] = v.alive
	}
}
