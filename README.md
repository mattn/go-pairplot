# go-pairplot

The gonum/plot/plotter package behave like seaborn.pairplot.

![](https://raw.githubusercontent.com/mattn/go-pairplot/master/screenshot.png)

## Usage

```go
p, err := plot.New()
if err != nil {
	log.Fatal(err)
}
pp, err := pairplot.NewPairPlot("iris.csv")
if err != nil {
	log.Fatal(err)
}
pp.Hue = "Name"
p.HideAxes()
p.Add(pp)
p.Save(8*vg.Inch, 8*vg.Inch, "example.png")
```

## Installation

```
$ go build
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
