package network_info

import (
	"bytes"
	"github.com/wcharczuk/go-chart/v2"
	"log"
	"signum-explorer-bot/internal/database/models"
	"time"
)

func (ni *NetworkInfoListener) GetNetworkChart() []byte {
	var networkInfos []models.NetworkInfo
	result := ni.db.Order("id desc").Find(&networkInfos)
	if result.Error != nil || len(networkInfos) == 0 {
		log.Printf("Error getting Network Infos from DB for plotting chart: %v", result.Error)
		return nil
	}

	graph := chart.Chart{
		Title: "Network Statistic",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 20,
			},
		},
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeMinuteValueFormatter,
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

	for _, values := range networkInfos {
		difficultyChartTimeSeries.XValues = append(difficultyChartTimeSeries.XValues, values.CreatedAt)
		difficultyChartTimeSeries.YValues = append(difficultyChartTimeSeries.YValues, values.NetworkDifficulty/1024)

		commitmentChartTimeSeries.XValues = append(commitmentChartTimeSeries.XValues, values.CreatedAt)
		commitmentChartTimeSeries.YValues = append(commitmentChartTimeSeries.YValues, values.AverageCommitment)
	}

	graph.Series = append(graph.Series, commitmentChartTimeSeries)
	graph.Series = append(graph.Series, difficultyChartTimeSeries)

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
