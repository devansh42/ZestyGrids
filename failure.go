package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *game) checkRestart() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		if inBox(float32(x), float32(y), screenW*0.4, screenH/2+fontSize, screenW*0.4+150, screenH/2+fontSize*2.5) {
			g.reset()
		}
	}
}
func (g *game) drawFailureText(mainscreen *ebiten.Image) {
	var failureTextOp, restartBtnOp text.DrawOptions
	failureTextOp.ColorScale.ScaleWithColor(color.White)
	failureTextOp.GeoM.Translate(screenW*0.2, screenH/2-fontSize)
	if g.fellIntoHole {
		text.Draw(mainscreen, "Boom! You messed up with the wrong guy", g.defaultFont, &failureTextOp)
	} else {
		text.Draw(mainscreen, "Oops! seems like you missed a lot of shots!", g.defaultFont, &failureTextOp)
	}

	restartBtnOp.GeoM.Translate(screenW*0.4, screenH/2+fontSize)
	restartBtnOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(mainscreen, "Restart", g.defaultFont, &restartBtnOp)
	vector.DrawFilledRect(mainscreen, screenW*0.4, screenH/2+fontSize, 150, fontSize*1.5, redColor, false)

}

func (g *game) reset() {
	*g = game{
		blinkDuration: time.Second * 1,
		blinkInterval: time.Second * 2,
		level:         1,
	}
	g.initFont()
}
