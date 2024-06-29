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

func NewResultsWithAnalysis(sim domain.SimResults, analyses analysis.AnalysisResults, tfOrder []domain.TrafficFlowConfig) (domain.Results, error) {
	var results simAnalysisResults

	results.simResults = localSimResults(sim.SimHeadlineResults)

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
				ID:       tfOrder[i].ID,
				Deadline: tfOrder[i].Deadline,
				StatSet:  tfSimStats,
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
			{Align: simpletable.AlignLeft, Text: "T_i"},
			{Align: simpletable.AlignLeft, Text: "No. pkts"},
			{Align: simpletable.AlignLeft, Text: "No. > D_i"},
			{Align: simpletable.AlignLeft, Text: "max"},
			{Align: simpletable.AlignLeft, Text: "mean"},
			{Align: simpletable.AlignLeft, Text: "min"},
			{Align: simpletable.AlignLeft, Text: "D_i"},
			{Align: simpletable.AlignLeft, Text: "J^R_i + C_i"},
			{Align: simpletable.AlignLeft, Text: "J^R_i + R_i"},
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
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsExceededDeadline)},
		{Align: simpletable.AlignLeft, Text: cleanInt(tf.BestLatency)},
		{Align: simpletable.AlignLeft, Text: cleanFloat(tf.MeanLatency)},
		{Align: simpletable.AlignLeft, Text: cleanInt(tf.WorstLatency)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Deadline)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Jitter + tf.Basic)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Jitter + tf.ShiAndBurns)},
		{Align: simpletable.AlignLeft, Text: strconv.FormatBool(tf.AnalysisHolds)},
	}

	return row
}

func (r *simAnalysisResults) OutputCSV(path string) error {
	data := [][]string{
		{
			"Traffic Flow",
			"Packets Routed",
			"Packets Exceeded Deadline",
			"No. Direct Interference",
			"No. Indirect Interference",
			"Best Latency",
			"Mean Latency",
			"Worst Latency",
			"Deadline",
			"Schedulable",
			"Jitter",
			"Jitter + Basic",
			"Jitter + Shi & Burns",
			"Analysis Schedulable",
			"Analysis Holds",
		},
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		data = append(data, []string{
			r.trafficFlows[i].ID,
			strconv.Itoa(r.trafficFlows[i].PacketsRouted),
			strconv.Itoa(r.trafficFlows[i].PacketsExceededDeadline),
			strconv.Itoa(r.trafficFlows[i].DirectInterferenceCount),
			strconv.Itoa(r.trafficFlows[i].IndirectInterferenceCount),
			cleanInt(r.trafficFlows[i].BestLatency),
			cleanFloat(r.trafficFlows[i].MeanLatency),
			cleanInt(r.trafficFlows[i].WorstLatency),
			strconv.Itoa(r.trafficFlows[i].Deadline),
			strconv.FormatBool(r.trafficFlows[i].Schedulable()),
			strconv.Itoa(r.trafficFlows[i].Jitter),
			strconv.Itoa(r.trafficFlows[i].Jitter + r.trafficFlows[i].Basic),
			strconv.Itoa(r.trafficFlows[i].Jitter + r.trafficFlows[i].ShiAndBurns),
			strconv.FormatBool(r.trafficFlows[i].AnalysisSchedulable()),
			strconv.FormatBool(r.trafficFlows[i].AnalysisHolds),
		})
	}

	return writeCSV(path, data)
}
