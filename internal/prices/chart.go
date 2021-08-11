package prices

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"log"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"time"
)

func (pm *PriceManager) GetPriceChart() []byte {
	var prices []models.Price
	result := pm.db.Order("id asc").Find(&prices)
	if result.Error != nil || len(prices) == 0 {
		log.Printf("Error getting Prices from DB for plotting chart: %v", result.Error)
		return nil
	}

	var max float64
	for _, value := range prices {
		if max < value.SignaPrice {
			max = value.SignaPrice
		}
	}

	var signaSign = "¢"
	var signaMultiplier float64 = 100
	if max >= 1 {
		signaSign = "$"
		signaMultiplier = 1
	}

	graph := chart.Chart{
		Title: fmt.Sprintf("SIGNA and BTC prices (last %v days)", config.CMC_API.LISTENER_DAYS_QUANTITY),
		Background: chart.Style{
			Padding: chart.Box{
				Top:  50,
				Left: 20,
			},
		},
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "SIGNA, " + signaSign,
		},
		YAxisSecondary: chart.YAxis{
			Name: "BTC, $",
		},
		Series: []chart.Series{},
	}

	signaChartTimeSeries := chart.TimeSeries{
		Name: "SIGNA",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(1),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(80),
		},
		XValues: []time.Time{},
		YValues: []float64{},
	}

	btcChartTimeSeries := chart.TimeSeries{
		Name: "BTC",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(20),
		},
		YAxis:   chart.YAxisSecondary,
		XValues: []time.Time{},
		YValues: []float64{},
	}

	annotationSeries := chart.AnnotationSeries{
		Annotations: []chart.Value2{},
	}

	for index, values := range prices {
		signaChartTimeSeries.XValues = append(signaChartTimeSeries.XValues, values.CreatedAt)
		signaChartTimeSeries.YValues = append(signaChartTimeSeries.YValues, values.SignaPrice*signaMultiplier)

		btcChartTimeSeries.XValues = append(btcChartTimeSeries.XValues, values.CreatedAt)
		btcChartTimeSeries.YValues = append(btcChartTimeSeries.YValues, values.BtcPrice)

		if index == len(prices)-1 {
			// annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{XValue: chart.TimeToFloat64(values.CreatedAt), YValue: values.SignaPrice * 100, Label: fmt.Sprintf("SIGNA - %.1f", values.SignaPrice*100)})
			// annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{XValue: chart.TimeToFloat64(values.CreatedAt), YValue: values.BtcPrice, Label: fmt.Sprintf("BTC - %.1f", values.BtcPrice)})
		}
	}

	graph.Series = append(graph.Series, signaChartTimeSeries)
	graph.Series = append(graph.Series, btcChartTimeSeries)
	graph.Series = append(graph.Series, annotationSeries)

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Printf("Could not render chart: %v", err)
		return nil
	}

	return buffer.Bytes()
}
