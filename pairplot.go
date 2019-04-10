package pairplot

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"image/png"
	"os"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type PairPlot struct {
	header []string
	data   [][]string

	Hue string
}

func NewPairPlotFromRows(header []string, data [][]string) (*PairPlot, error) {
	pp := &PairPlot{}
	pp.header = header
	pp.data = data
	return pp, nil
}

func NewPairPlot(fname string) (*PairPlot, error) {
	pp := &PairPlot{}

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	pp.header, err = r.Read()
	if err != nil {
		return nil, err
	}
	pp.data, err = r.ReadAll()
	if err != nil {
		return nil, err
	}

	return pp, nil
}

func (pp *PairPlot) Plot(c draw.Canvas, p *plot.Plot) {
	names := []string{}
	attr := -1
	for i, s := range pp.header {
		if s == pp.Hue {
			attr = i
			for _, row := range pp.data {
				name := row[i]
				found := false
				for _, n := range names {
					if n == name {
						found = true
						break
					}
				}
				if !found {
					names = append(names, name)
				}
			}
			break
		}
	}
	if len(names) == 0 {
		names = []string{""}
	}

	dx := c.Size().X.Points() / float64(len(pp.header)-1)
	dy := c.Size().Y.Points() / float64(len(pp.header)-1)

	for i1 := 0; i1 < len(pp.header); i1++ {
		for i2 := 0; i2 < len(pp.header); i2++ {
			p, err := plot.New()
			if err != nil {
				return
			}

			p.Add(plotter.NewGrid())

			empty := true
			if i1 == i2 {
				var vals []plotter.Values
				first := true
				var min, max float64
				for _, name := range names {
					var values plotter.Values
					for _, row := range pp.data {
						if name != "" && row[attr] != name {
							continue
						}
						f, err := strconv.ParseFloat(row[i1], 64)
						if err != nil {
							continue
						}
						if first {
							first = false
							min, max = f, f
						} else {
							if f < min {
								min = f
							}
							if f > max {
								max = f
							}
						}
						values = append(values, f)
					}
					if len(values) == 0 {
						continue
					}
					sort.Float64s(values)
					vals = append(vals, values)
				}
				if len(vals) == 0 {
					continue
				}
				divider := []float64{}
				i := min
				for {
					divider = append(divider, float64(i))
					if i > max {
						break
					}
					i += (max - min) / 10.0
				}
				var prev *plotter.BarChart
				for ci, val := range vals {
					fval := stat.Histogram(nil, divider, val, nil)
					bar, err := plotter.NewBarChart(plotter.Values(fval), vg.Length(20))
					if err != nil {
						panic(err)
					}
					if prev != nil {
						bar.StackOn(prev)
					}
					prev = bar
					bar.Color = plotutil.Color(ci)
					p.Add(bar)
					empty = false
				}
				if i2 == 0 {
					p.Y.Label.Text = pp.header[i1]
				} else {
					p.Y.Label.Text = " "
				}
				p.Y.Label.Font.Size = 40.5
				p.Legend.Top = true
			} else {
				if i2 == 0 {
					p.Y.Label.Text = pp.header[i1]
				} else {
					p.Y.Label.Text = " "
				}
				p.Y.Label.Font.Size = 40.5
				for ci, name := range names {
					var xys plotter.XYs
					for _, row := range pp.data {
						if name != "" && row[attr] != name {
							continue
						}
						f1, err := strconv.ParseFloat(row[i1], 64)
						if err != nil {
							continue
						}
						f2, err := strconv.ParseFloat(row[i2], 64)
						if err != nil {
							continue
						}
						xys = append(xys, plotter.XY{
							X: f1,
							Y: f2,
						})
					}
					if len(xys) == 0 {
						continue
					}

					plotData, err := plotter.NewScatter(xys)
					if err != nil {
						return
					}
					plotData.GlyphStyle.Color = plotutil.Color(ci)
					plotData.GlyphStyle.Radius = c.Size().X / 64
					plotData.GlyphStyle.Shape = draw.CircleGlyph{}
					p.Add(plotData)
					empty = false
				}
			}

			if !empty {
				w, err := p.WriterTo(c.Size().X, c.Size().Y, "png")
				if err != nil {
					return
				}
				var buf bytes.Buffer
				writer := bufio.NewWriter(&buf)
				w.WriteTo(writer)
				writer.Flush()

				img, err := png.Decode(bytes.NewReader(buf.Bytes()))
				if err != nil {
					return
				}

				rect := vg.Rectangle{
					Min: vg.Point{
						X: vg.Length(dx * float64(i2)),
						Y: vg.Length(dy * float64(len(pp.header)-i1-2)),
					},
					Max: vg.Point{
						X: vg.Length(dx * float64(i2+1)),
						Y: vg.Length(dy * float64(len(pp.header)-i1-1)),
					},
				}
				c.DrawImage(rect, img)
			}
		}
	}
}
