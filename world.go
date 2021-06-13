package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

var heatMapClr map[int][4]byte = map[int][4]byte{
	0:  [4]byte{0xff, 0xff, 0xff, 0xff},
	5:  [4]byte{0xff, 0xf3, 0x3b, 0xff},
	10: [4]byte{0xfd, 0xc7, 0x0c, 0xff},
	15: [4]byte{0xf3, 0x90, 0x3f, 0xff},
	20: [4]byte{0xed, 0x68, 0x3c, 0xff},
	25: [4]byte{0xe9, 0x3e, 0x3a, 0xff},
}

type Cell struct {
	alive bool
	age   int
}

type World struct {
	width  int
	height int
	mtx    [][]Cell
}

func (w *World) heavyWeightSpaceship(x, y int) {
	// Check out of bounds
	if x+6 >= w.width || y+4 >= w.height || x < 0 || y < 0 {
		log.Print("heavy weight spaceship index out of range")
		return
	}

	hws := []point{
		point{x: 2, y: 0, alive: true},
		point{x: 3, y: 0, alive: true},
		point{x: 0, y: 1, alive: true},
		point{x: 5, y: 1, alive: true},
		point{x: 6, y: 2, alive: true},
		point{x: 0, y: 3, alive: true},
		point{x: 6, y: 3, alive: true},
		point{x: 1, y: 4, alive: true},
		point{x: 2, y: 4, alive: true},
		point{x: 3, y: 4, alive: true},
		point{x: 4, y: 4, alive: true},
		point{x: 5, y: 4, alive: true},
		point{x: 6, y: 4, alive: true},
	}

	for _, v := range hws {
		w.mtx[x+v.x][y+v.y].alive = v.alive
	}
}

func (w *World) randPopulate() {
	for i := 0; i < 1337; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.mtx[x][y].alive = true
		w.mtx[x][y].age = 1
	}
}

func (w *World) reset() {
	for x := range w.mtx {
		for y := range w.mtx[x] {
			w.mtx[x][y].alive = false
			w.mtx[x][y].age = 0
		}
	}
	globalTick = 0
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

			curr := w.mtx[x][y]

			// needs to sort
			keys := []int{}
			for k := range heatMapClr {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			// find interval
			interClr := 0
			for _, k := range keys {
				if curr.age > k {
					interClr = k
					continue
				}
				break
			}

			for idx, pos := range tmp {
				if w.mtx[x][y].alive {
					p[pos] = heatMapClr[interClr][idx]
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
			if w.mtx[x][y].alive {
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

				curr.alive = w.mtx[curr.x][curr.y].alive
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
			if w.mtx[x][y].alive == true && (alive == 2 || alive == 3) {
				w.mtx[x][y].age++
				continue
			}

			// Any dead cell with three live neighbours becomes a live cell.
			if w.mtx[x][y].alive == false && alive == 3 {
				w.mtx[x][y].age = 1
				next = append(next, point{x: x, y: y, alive: true})
				continue
			}

			// All other live cells die in the next generation. Similarly, all other dead cells stay dead.
			if w.mtx[x][y].alive == true {
				w.mtx[x][y].age = 0
				next = append(next, point{x: x, y: y, alive: false})
			}
		}
	}

	currDead := w.width*w.height - currAlive
	currAliveRatio := (float64(currAlive) / float64(currAlive+currDead)) * 100
	ebiten.SetWindowTitle(
		fmt.Sprintf("tps: %d | gen: %d | alive: %d | dead: %d | alive ratio: %.2f%%",
			int(math.Round(ebiten.CurrentTPS())),
			globalTick,
			currAlive,
			currDead,
			currAliveRatio,
		))

	for _, v := range next {
		w.mtx[v.x][v.y].alive = v.alive
	}
}
