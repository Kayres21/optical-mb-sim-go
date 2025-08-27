package plotter

import (
	"fmt"
	"image/color"
	"path/filepath"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func GenerateScatterPlot(xData, yData []float64, title, xLabel, yLabel string) error {
	if len(xData) != len(yData) {
		return fmt.Errorf("los datos de X y Y deben tener la misma longitud")
	}

	points := make(plotter.XYs, len(xData))
	for i := range xData {
		points[i].X = xData[i]
		points[i].Y = yData[i]
	}

	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	s, err := plotter.NewScatter(points)
	if err != nil {
		return err
	}

	s.Color = color.RGBA{R: 255, A: 255}
	s.Radius = vg.Points(3)

	p.Add(s)

	filename := fmt.Sprintf("%s.png", title)
	filePath := filepath.Join("result", filename)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filePath); err != nil {
		return err
	}

	fmt.Printf("Gráfico generado con éxito en '%s'\n", filePath)
	return nil
}
