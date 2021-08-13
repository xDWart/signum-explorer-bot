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
		Title: fmt.Sprintf("SIGNA and BTC prices (last %v days)", config.CMC_API.SAVING_DAYS_QUANTITY),
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

	var annotationColor = chart.ColorGreen
	actualPrices := pm.cmcClient.GetPrices()
	if actualPrices["SIGNA"].PercentChange24h < 0 {
		annotationColor = chart.ColorRed
	}

	for _, values := range prices {
		signaChartTimeSeries.XValues = append(signaChartTimeSeries.XValues, values.CreatedAt)
		signaChartTimeSeries.YValues = append(signaChartTimeSeries.YValues, values.SignaPrice*signaMultiplier)

		btcChartTimeSeries.XValues = append(btcChartTimeSeries.XValues, values.CreatedAt)
		btcChartTimeSeries.YValues = append(btcChartTimeSeries.YValues, values.BtcPrice)
	}

	lastPrices := pm.cmcClient.GetPrices()
	signaChartTimeSeries.XValues = append(signaChartTimeSeries.XValues, time.Now())
	signaChartTimeSeries.YValues = append(signaChartTimeSeries.YValues, lastPrices["SIGNA"].Price*signaMultiplier)
	btcChartTimeSeries.XValues = append(btcChartTimeSeries.XValues, time.Now())
	btcChartTimeSeries.YValues = append(btcChartTimeSeries.YValues, lastPrices["BTC"].Price)
	annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{
		XValue: chart.TimeToFloat64(time.Now()),
		YValue: lastPrices["SIGNA"].Price * signaMultiplier,
		Label:  fmt.Sprintf("%.2f", lastPrices["SIGNA"].Price*signaMultiplier),
		Style:  chart.Style{StrokeColor: annotationColor}})
	//annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{
	//	XValue: chart.TimeToFloat64(time.Now()),
	//	YValue: lastPrices["BTC"].Price,
	//	Label:  fmt.Sprintf("%.f", lastPrices["BTC"].Price)})

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
