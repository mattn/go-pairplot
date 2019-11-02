package pairplot

import (
	"os"

	"github.com/go-gota/gota/dataframe"
)

func NewPairPlotFromRows(header []string, data [][]string) (PairPlot, error) {
	df := dataframe.LoadRecords(append([][]string{header}, data...), dataframe.HasHeader(true))
	pp := &dfPairPlot{df: df}
	return pp, nil
}

func NewPairPlotCSV(fname string) (PairPlot, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)
	pp := &dfPairPlot{df: df}
	return pp, nil
}
