package main

import (
	"log"

	"github.com/mattn/go-pairplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	pp, err := pairplot.NewPairPlotCSV("iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	pp.SetHue("Name")
	p.HideAxes()
	p.Add(pp)
	p.Save(8*vg.Inch, 8*vg.Inch, "example.png")
}
