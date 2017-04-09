package chart

import (
	"github.com/wcharczuk/go-chart"
	"fmt"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func DrawTcpappChart(resultFile string, outFile string) error {

	b, err := ioutil.ReadFile(resultFile)
	if err != nil {
		return err
	}

	var measure migrate.Measure
	err = json.Unmarshal(b, &measure)
	if err != nil {
		return err
	}

	dataLen := len(measure.ClientResults)

	times := make([]time.Time, dataLen)
	values := make([]float64, dataLen)

	for i, res := range measure.ClientResults {
		times[i] = res.Time
		values[i] = float64(res.Rtt)
	}

	annotations := make([]chart.Value2, 0)
	for i, mileStone := range measure.Milestones {
		annotations = append(annotations, chart.Value2{
			XValue: float64(mileStone.At.UnixNano()),
			Style: chart.Style{
				Show:            true,
				StrokeColor:     chart.ColorAlternateGray,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			Label: fmt.Sprintf("%d.", i+1),
		})
	}

	mainSeries := chart.TimeSeries{
		XValues: times,
		YValues: values,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlue,
			FillColor:   chart.ColorBlue.WithAlpha(100),
		},
	}

	minSeries := &chart.MinSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	maxSeries := &chart.MaxSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Width: 1600,
		Height: 500,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 0,
			},
		},
		XAxis: chart.XAxis{
			//TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.Time.FromFloat64(typed)
				return fmt.Sprintf("%02d:%03d", typedDate.Second(), typedDate.Nanosecond() / 1000000)
			},
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 0.5,
			},
			Name:      "Time (sec:ms)",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 0.0,
			},
		},
		YAxis: chart.YAxis{
			Name:      "Response time",
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%.0f ms", typed / 1000000)
			},
		},
		Series: []chart.Series{
			mainSeries,
			chart.AnnotationSeries{
				Annotations: annotations,
			},
			minSeries,
			maxSeries,
			chart.LastValueAnnotation(minSeries),
			chart.LastValueAnnotation(maxSeries),
		},
	}
	//graph.Elements = []chart.Renderable{ chart.Legend(&graph) }


	w, err := os.Create(outFile)
	if err != nil {
		return err
	}
	graph.Render(chart.SVG, w)

	return w.Close()
}