package networkinfo

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"log"
	"signum-explorer-bot/internal/config"
	"signum-explorer-bot/internal/database/models"
	"time"
)

func (ni *NetworkInfoListener) GetNetworkChart(duration time.Duration) []byte {
	var networkInfos []models.NetworkInfo
	result := ni.db.Where("created_at > ?", time.Now().Add(-duration)).Order("id asc").Find(&networkInfos)
	if result.Error != nil || len(networkInfos) == 0 {
		log.Printf("Error getting Network Infos from DB for plotting chart: %v", result.Error)
		return nil
	}

	var lastText = "since rebranding"
	switch duration {
	case config.DAY:
		lastText = "last 24 hours"
	case config.WEEK:
		lastText = "last week"
	case config.MONTH:
		lastText = "last month"
	}

	graph := chart.Chart{
		Title: fmt.Sprintf("Network Statistic (%v)", lastText),
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
			Name: "Commitment, SIGNA / TiB",
		},
		YAxisSecondary: chart.YAxis{
			Name: "Difficulty, PiB",
		},
		Series: []chart.Series{},
	}

	difficultyChartTimeSeries := chart.TimeSeries{
		Name: "Difficulty",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
		},
		YAxis:   chart.YAxisSecondary,
		XValues: []time.Time{},
		YValues: []float64{},
	}

	commitmentChartTimeSeries := chart.TimeSeries{
		Name: "Commitment",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(1),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
		},
		XValues: []time.Time{},
		YValues: []float64{},
	}

	annotationSeries := chart.AnnotationSeries{
		Annotations: []chart.Value2{},
	}

	for _, values := range networkInfos {
		difficultyChartTimeSeries.XValues = append(difficultyChartTimeSeries.XValues, values.CreatedAt)
		difficultyChartTimeSeries.YValues = append(difficultyChartTimeSeries.YValues, values.NetworkDifficulty/1024)

		commitmentChartTimeSeries.XValues = append(commitmentChartTimeSeries.XValues, values.CreatedAt)
		commitmentChartTimeSeries.YValues = append(commitmentChartTimeSeries.YValues, values.AverageCommitment)
	}

	actualMiningInfo := ni.GetLastMiningInfo()
	difficultyChartTimeSeries.XValues = append(difficultyChartTimeSeries.XValues, time.Now())
	difficultyChartTimeSeries.YValues = append(difficultyChartTimeSeries.YValues, actualMiningInfo.ActualNetworkDifficulty/1024)
	commitmentChartTimeSeries.XValues = append(commitmentChartTimeSeries.XValues, time.Now())
	commitmentChartTimeSeries.YValues = append(commitmentChartTimeSeries.YValues, actualMiningInfo.ActualCommitment)
	annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{
		XValue: chart.TimeToFloat64(time.Now()),
		YValue: actualMiningInfo.ActualCommitment,
		Label:  fmt.Sprintf("%.f.00", actualMiningInfo.ActualCommitment),
		Style:  chart.Style{StrokeColor: chart.ColorGreen}})
	// annotationSeries.Annotations = append(annotationSeries.Annotations, chart.Value2{XValue: chart.TimeToFloat64(values.CreatedAt), YValue: values.NetworkDifficulty / 1024, Label: fmt.Sprintf("%.1f", values.NetworkDifficulty/1024)})

	graph.Series = append(graph.Series, commitmentChartTimeSeries)
	graph.Series = append(graph.Series, difficultyChartTimeSeries)
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
