package results

import (
	"strconv"

	"main/src/core/analysis"
	"main/src/domain"

	"github.com/alexeyco/simpletable"
)

type simAnalysisResults struct {
	simResults   localSimResults
	trafficFlows []tfSimAnalysis
}

func NewResultsWithAnalysis(sim domain.FullResults, analyses analysis.AnalysisResults, tfOrder []domain.TrafficFlowConfig) (domain.Results, error) {
	var results simAnalysisResults

	results.simResults = localSimResults(sim.SimResults)

	for i := 0; i < len(tfOrder); i++ {
		tfSimStats, exists := sim.TFStats[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		tfAnalysis, exists := analyses[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		analysisHolds := true
		// if tfAnalysis.AnalysisSchedulable() && (tfSimStats.WorstLatency > (tfOrder[i].Jitter + tfAnalysis.ShiAndBurns)) {
		if tfAnalysis.AnalysisSchedulable() && !tfSimStats.Schedulable() {
			analysisHolds = false
		}

		results.trafficFlows = append(results.trafficFlows, tfSimAnalysis{
			tfSim: tfSim{
				ID:                 tfOrder[i].ID,
				Deadline:           tfOrder[i].Deadline,
				TrafficFlowStatSet: tfSimStats,
			},
			TrafficFlowAnalysisSet: tfAnalysis,
			AnalysisHolds:          analysisHolds,
		})
	}

	return &results, nil
}

func (r *simAnalysisResults) Prettify() (string, error) {
	var str string

	str += r.simResults.prettify()
	str += r.prettifyTfInfo()

	return str, nil
}

func (r *simAnalysisResults) prettifyTfInfo() (str string) {
	str += "Traffic Flow domain.Results\n"
	str += "====================\n"
	str += r.prettifyTfTable()

	return str
}

func (r *simAnalysisResults) prettifyTfTable() (str string) {
	tfTable := simpletable.New()

	tfTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "Traffic Flow"},
			{Align: simpletable.AlignLeft, Text: "Packets Routed"},
			{Align: simpletable.AlignLeft, Text: "Packets Arrived"},
			{Align: simpletable.AlignLeft, Text: "Packets Exceeded Deadline"},
			{Align: simpletable.AlignLeft, Text: "Best Latency"},
			{Align: simpletable.AlignLeft, Text: "Mean Latency"},
			{Align: simpletable.AlignLeft, Text: "Worst Latency"},
			{Align: simpletable.AlignLeft, Text: "Deadline"},
			{Align: simpletable.AlignLeft, Text: "Jitter"},
			{Align: simpletable.AlignLeft, Text: "Basic"},
			{Align: simpletable.AlignLeft, Text: "Shi & Burns"},
			{Align: simpletable.AlignLeft, Text: "Analysis Holds"},
		},
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		tfTable.Body.Cells = append(tfTable.Body.Cells, r.prettifyTfRow(r.trafficFlows[i]))
	}

	str += tfTable.String()
	str += "\n"

	return str
}

func (r *simAnalysisResults) prettifyTfRow(tf tfSimAnalysis) []*simpletable.Cell {
	row := []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: tf.ID},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsRouted)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsArrived)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsExceededDeadline)},
		{Align: simpletable.AlignLeft, Text: cleanBestLatency(tf.BestLatency)},
		{Align: simpletable.AlignLeft, Text: cleanMeanLatency(tf.MeanLatency)},
		{Align: simpletable.AlignLeft, Text: cleanWorstLatency(tf.WorstLatency)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Deadline)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Jitter)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Basic)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.ShiAndBurns)},
		{Align: simpletable.AlignLeft, Text: strconv.FormatBool(tf.AnalysisHolds)},
	}

	return row
}

func (r *simAnalysisResults) OutputCSV(path string) error {
	data := [][]string{
		{
			"Traffic Flow",
			"Packets Routed",
			"Packets Arrived",
			"Packets Exceeded Deadline",
			"Packets Lost",
			"Direct Interference",
			"Indirect Interference",
			"Best Latency",
			"Mean Latency",
			"Worst Latency",
			"Deadline",
			"Jitter",
			"Basic",
			"Shi & Burns",
			"Analysis Schedulable",
			"Schedulable",
			"Analysis Holds",
		},
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		data = append(data, []string{
			r.trafficFlows[i].ID,
			strconv.Itoa(r.trafficFlows[i].PacketsRouted),
			strconv.Itoa(r.trafficFlows[i].PacketsArrived),
			strconv.Itoa(r.trafficFlows[i].PacketsExceededDeadline),
			strconv.Itoa(r.trafficFlows[i].PacketsLost),
			strconv.Itoa(r.trafficFlows[i].DirectInterferenceCount),
			strconv.Itoa(r.trafficFlows[i].IndirectInterferenceCount),
			cleanBestLatency(r.trafficFlows[i].BestLatency),
			cleanMeanLatency(r.trafficFlows[i].MeanLatency),
			cleanWorstLatency(r.trafficFlows[i].WorstLatency),
			strconv.Itoa(r.trafficFlows[i].Deadline),
			strconv.Itoa(r.trafficFlows[i].Jitter),
			strconv.Itoa(r.trafficFlows[i].Basic),
			strconv.Itoa(r.trafficFlows[i].ShiAndBurns),
			strconv.FormatBool(r.trafficFlows[i].AnalysisSchedulable()),
			strconv.FormatBool(r.trafficFlows[i].Schedulable()),
			strconv.FormatBool(r.trafficFlows[i].AnalysisHolds),
		})
	}

	return writeCSV(path, data)
}
