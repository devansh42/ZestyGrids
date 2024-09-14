package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *game) drawText(mainscreen *ebiten.Image) {
	text.Draw(mainscreen, fmt.Sprintf("Level: %d", g.level), g.defaultFont, &g.levelFontDrawOpts)
	text.Draw(mainscreen, fmt.Sprintf("Score: %d", g.hits), g.defaultFont, &g.hitsDrawOpts)
}

func (g *game) initFont() {
	x, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	g.defaultFont = &text.GoTextFace{
		Source: x,
		Size:   fontSize,
	}
	var op text.DrawOptions
	op.GeoM.Translate(textHorizontalMargin, textVerticalMargin)
	op.ColorScale.ScaleWithColor(color.White)
	g.levelFontDrawOpts = op

	g.hitsDrawOpts = text.DrawOptions{}
	g.hitsDrawOpts.ColorScale.ScaleWithColor(color.White)
	g.hitsDrawOpts.GeoM.Translate((screenW*3)/4, textVerticalMargin)
}
