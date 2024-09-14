package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Zesty Grids")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
