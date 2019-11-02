package pairplot

import (
	"bufio"
	"bytes"
	"image/png"
	"sort"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type dfPairPlot struct {
	df dataframe.DataFrame

	hue string
}

func NewPairPlotDataFrame(df dataframe.DataFrame) (PairPlot, error) {
	pp := &dfPairPlot{df: df}
	return pp, nil
}

func (pp *dfPairPlot) bars(c draw.Canvas, p *plot.Plot, i1, i2 int) bool {
	var vals []plotter.Values
	first := true
	var min, max float64

	names := pp.df.Names()

	s1 := pp.df.Col(names[i1])
	s3 := pp.df.Col(pp.hue)
	f1 := s1.Float()

	for n := int(s3.Min()); n <= int(s3.Max()); n++ {
		var values plotter.Values
		for i := 0; i < len(f1); i++ {
			if first {
				min = f1[i]
				max = f1[i]
				first = false
			} else {
				smin := f1[i]
				smax := f1[i]
				if min > smin {
					min = smin
				}
				if max < smax {
					max = smax
				}
			}
			values = append(values, f1[i])
		}
		if len(values) == 0 {
			continue
		}
		sort.Float64s(values)
		vals = append(vals, values)
	}
	if len(vals) == 0 {
		return false
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

	if i2 == 0 {
		p.Y.Label.Text = pp.df.Names()[i1]
	} else {
		p.Y.Label.Text = " "
	}
	p.Y.Label.Font.Size = fontSize
	p.Legend.Top = true

	var prev *plotter.BarChart

	drawn := false
	for ci, val := range vals {
		fval := stat.Histogram(nil, divider, val, nil)
		bar, err := plotter.NewBarChart(plotter.Values(fval), vg.Length(20))
		if err != nil {
			return false
		}
		if prev != nil {
			bar.StackOn(prev)
		}
		prev = bar
		bar.Color = plotutil.Color(ci)
		p.Add(bar)
		drawn = true
	}

	return drawn
}

func (pp *dfPairPlot) scatter(c draw.Canvas, p *plot.Plot, i1, i2 int) bool {
	if i2 == 0 {
		p.Y.Label.Text = pp.df.Names()[i1]
	} else {
		p.Y.Label.Text = " "
	}
	p.Y.Label.Font.Size = fontSize

	names := pp.df.Names()

	s1 := pp.df.Col(names[i1])
	s2 := pp.df.Col(names[i2])
	s3 := pp.df.Col(pp.hue)
	f1 := s1.Float()
	f2 := s2.Float()
	f3 := s3.Float()

	for n := int(s3.Min()); n <= int(s3.Max()); n++ {
		var xys plotter.XYs
		for i := 0; i < len(f1); i++ {
			if int(f3[i]) != n {
				continue
			}
			xys = append(xys, plotter.XY{
				X: f1[i],
				Y: f2[i],
			})
		}
		plotData, err := plotter.NewScatter(xys)
		if err != nil {
			return false
		}
		plotData.GlyphStyle.Color = plotutil.Color(n)
		plotData.GlyphStyle.Radius = c.Size().X / 64
		plotData.GlyphStyle.Shape = draw.CircleGlyph{}
		p.Add(plotData)
	}

	return true
}

func (pp *dfPairPlot) plot(c draw.Canvas, p *plot.Plot, i1, i2 int) {
	p, err := plot.New()
	if err != nil {
		return
	}

	p.Add(plotter.NewGrid())

	var drawn bool
	if i1 == i2 {
		drawn = pp.bars(c, p, i1, i2)
	} else {
		drawn = pp.scatter(c, p, i1, i2)
	}
	if !drawn {
		return
	}
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

	// width/height for per graph.
	dx := c.Size().X.Points() / float64(pp.df.Ncol()-1)
	dy := c.Size().Y.Points() / float64(pp.df.Ncol()-1)

	rect := vg.Rectangle{
		Min: vg.Point{
			X: vg.Length(dx * float64(i2)),
			Y: vg.Length(dy * float64(pp.df.Ncol()-i1-2)),
		},
		Max: vg.Point{
			X: vg.Length(dx * float64(i2+1)),
			Y: vg.Length(dy * float64(pp.df.Ncol()-i1-1)),
		},
	}
	c.DrawImage(rect, img)
}

func (pp *dfPairPlot) Plot(c draw.Canvas, p *plot.Plot) {
	pp.df = pp.df.Capply(func(s series.Series) series.Series {
		if s.Name != pp.hue {
			return s
		}
		records := s.Records()
		floats := make([]float64, len(records))
		m := map[string]int{}
		for i, r := range records {
			if _, ok := m[r]; !ok {
				m[r] = len(m)
			}
			floats[i] = float64(m[r])
		}
		return series.Floats(floats)
	})

	for i1 := 0; i1 < pp.df.Ncol(); i1++ {
		for i2 := 0; i2 < pp.df.Ncol(); i2++ {
			pp.plot(c, p, i1, i2)
		}
	}
}

func (pp *dfPairPlot) SetHue(s string) {
	pp.hue = s
}
