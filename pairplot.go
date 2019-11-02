package pairplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
)

const fontSize = 40.5

type PairPlot interface {
	Plot(c draw.Canvas, p *plot.Plot)
	SetHue(string)
}
