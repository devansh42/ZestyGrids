package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW                      = 1080
	screenH                      = 1080
	OriginX              float32 = 0
	OriginY              float32 = 0
	boardW                       = 640
	boardH                       = 640
	gridW                float32 = boardW / gridCount
	gridH                float32 = boardH / gridCount
	gridCount                    = 4
	leftMargin                   = 220
	topMargin                    = 180
	strokeBorder                 = 3
	textHorizontalMargin         = 50
	textVerticalMargin           = 50
	fontSize                     = 36
)

var (
	bgColor     = color.RGBA{242, 224, 213, 1}
	greyColor   = color.RGBA{132, 129, 129, 0}
	redColor    = color.RGBA{255, 0, 0, 0}
	happyColors = []color.RGBA{{255, 215, 0, 1}, // Sunshine Yellow
		{135, 206, 235, 1}, // Sky Blue
		{152, 255, 152, 1}, // Mint Green
		{255, 111, 97, 1},  // Coral
		{255, 105, 180, 1}, // Bright Pink
		{230, 230, 250, 1}, // Lavender
		{255, 165, 0, 1},   // Golden Orange
		{50, 205, 50, 1},   // Lime Green
		{64, 224, 208, 1},  // Turquoise
	}
)

type game struct {
	failed                          bool
	level                           int
	hits                            int
	levelHits                       int
	levelTrials                     int
	defaultFont                     text.Face
	targetGrids                     [3][2]int
	gridColors                      [3]int
	evilInTargetGrids               int
	levelFontDrawOpts, hitsDrawOpts text.DrawOptions
	blinkInterval                   time.Duration
	lastUpdated                     time.Time
	blinkDuration                   time.Duration
	blinkVisible                    bool
	activeGrid                      int
	fellIntoHole                    bool
}

func newGame() *game {
	var g = game{
		level:         1,
		blinkInterval: time.Second * 2,
		blinkDuration: time.Second * 1,
		activeGrid:    1,
	}

	g.initFont()
	return &g
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.failed {
		g.drawFailureText(screen)
		return
	}
	g.drawGameRectangle(screen)
	g.drawText(screen)
}
func (g *game) drawGameRectangle(mainscreen *ebiten.Image) {
	var drawnActiveIndex int
	for i := 0; i < gridCount; i++ {
		for j := 0; j < gridCount; j++ {
			posX, posY := g.gridDimm(i, j)
			if g.isActiveGrid(i, j) && !g.isEvilGrid(i, j) {
				vector.DrawFilledRect(mainscreen, posX, posY,
					gridW, gridH, happyColors[g.gridColors[g.activeIndex(i, j)]], false)
				drawnActiveIndex++
			} else if g.isEvilGrid(i, j) {
				vector.DrawFilledRect(mainscreen, posX, posY,
					gridW, gridH, color.Black, false)
			} else {
				vector.DrawFilledRect(mainscreen, posX, posY, gridW, gridH, bgColor, false)

			}
			vector.StrokeRect(mainscreen,
				posX,
				posY,
				gridW, gridH, strokeBorder, greyColor, false)
		}
	}
}

func (g *game) activeIndex(i, j int) int {
	for x := 0; x < g.activeGrid; x++ {
		if g.targetGrids[x] == [2]int{i, j} {
			return x
		}
	}
	return -1
}

func (g *game) isActiveGrid(row, col int) bool {
	for i := 0; i < g.activeGrid; i++ {
		if g.targetGrids[i][0] == row && g.targetGrids[i][1] == col {
			return true
		}
	}
	return false
}

func (g *game) isEvilGrid(row, col int) bool {
	if g.evilInTargetGrids > -1 {
		return g.targetGrids[g.evilInTargetGrids][0] == row && g.targetGrids[g.evilInTargetGrids][1] == col
	}
	return false
}

func (g *game) gridDimm(row, col int) (float32, float32) {
	return leftMargin + OriginX + float32(row)*gridW, topMargin + OriginY + float32(col)*gridH

}

func (g *game) Update() error {
	if g.failed {
		g.checkRestart()
		return nil
	}

	now := time.Now()
	dur := now.Sub(g.lastUpdated)
	if !g.blinkVisible && dur > g.blinkInterval {
		if g.increaseLevel() {
			g.failed = true
			return nil
		}
		g.activateGrids()
		g.lastUpdated = now
		g.blinkVisible = true
		g.levelTrials++
	} else if g.blinkVisible && dur > g.blinkDuration {
		g.blinkVisible = false
		g.resetPaintedGrids()
	} else if g.blinkVisible {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			scored, scoredEvil := g.scored(ebiten.CursorPosition())
			if scored {
				if scoredEvil {
					g.failed = true
					g.fellIntoHole = true
					return nil
				}
				g.hits++
				g.blinkVisible = false
				g.resetPaintedGrids()
				g.levelHits++
			}
		}
	}

	return nil
}

func (g *game) resetPaintedGrids() {
	for i := 0; i < g.activeGrid; i++ {
		for j := 0; j < 2; j++ {
			g.targetGrids[i][j] = -1
		}
	}
}

func (g *game) activateGrids() {
	g.activeGrid = 1 + rand.Intn(3)                    // Atleast one grid
	g.evilInTargetGrids = -1 + rand.Intn(g.activeGrid) // Index for evil grid, if -1 then we would ignore evil grid
	var rows, cols byte                                // Assuming matrix is less than 8x8
	var colors uint16                                  // Assuming we only have at most 16 colors
	for i := 0; i < g.activeGrid; i++ {
		var grid [2]int
		for grid = [2]int{rand.Intn(gridCount), rand.Intn(gridCount)}; rows&(1<<grid[0]) == (1<<grid[0]) && cols&(1<<grid[1]) == (1<<grid[1]); grid = [2]int{rand.Intn(gridCount), rand.Intn(gridCount)} {
		}

		rows |= 1 << (grid[0])
		cols |= 1 << (grid[1])

		var color int
		for color = rand.Intn(len(happyColors)); colors&(1<<(color)) == 1<<color; color = rand.Intn(len(happyColors)) {
		}
		colors |= (1 << color)
		g.gridColors[i] = color
		g.targetGrids[i] = grid
	}
}

// increaseLevel increases game level
// this affects next tick
func (g *game) increaseLevel() (levelFailed bool) {
	if g.levelTrials >= 8 {
		if g.levelHits < 4 {
			return true
		}
		g.levelTrials = 0
		g.levelHits = 0
		g.level++
		// 2% drop per level
		g.blinkInterval = time.Millisecond * time.Duration(2000*math.Pow(0.98, float64(g.level)-1))
		g.blinkDuration = time.Millisecond * time.Duration(1000*math.Pow(0.98, float64(g.level)-1))
	}
	return
}

func (g *game) scored(posX, posY int) (scoredGrid bool, wasEvilGrid bool) {
	fx, fy := float32(posX), float32(posY)

	for i := 0; i < g.activeGrid; i++ {
		gridX, gridY := g.gridDimm(g.targetGrids[i][0], g.targetGrids[i][1])
		cornerX, cornerY := g.gridDimm(g.targetGrids[i][0]+1, g.targetGrids[i][1]+1)
		if inBox(fx, fy, gridX, gridY, cornerX, cornerY) {
			return true, g.evilInTargetGrids == i
		}
	}
	return
}

func inBox(fx, fy, boxX, boxY, cornerX, cornerY float32) bool {
	return fx >= boxX && fx <= cornerX && fy >= boxY && fy <= cornerY

}

func (g *game) Layout(w, h int) (int, int) {
	return screenW, screenH
}
