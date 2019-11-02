package main

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/mattn/go-pairplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	df := dataframe.ReadCSV(f)
	pp, err := pairplot.NewPairPlotDataFrame(df)
	if err != nil {
		log.Fatal(err)
	}
	pp.SetHue("Name")
	p.HideAxes()
	p.Add(pp)
	p.Save(8*vg.Inch, 8*vg.Inch, "example.png")
}
