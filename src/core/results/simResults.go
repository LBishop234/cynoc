package results

import (
	"strconv"

	"main/src/domain"

	"github.com/alexeyco/simpletable"
)

type simResults struct {
	simResults   localSimResults
	trafficFlows []tfSim
}

func NewResults(sim domain.FullResults, tfOrder []domain.TrafficFlowConfig) (domain.Results, error) {
	var results simResults

	results.simResults = localSimResults(sim.SimResults)

	for i := 0; i < len(tfOrder); i++ {
		tfStats, exists := sim.TFStats[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		results.trafficFlows = append(results.trafficFlows, tfSim{
			ID:                 tfOrder[i].ID,
			Deadline:           tfOrder[i].Deadline,
			TrafficFlowStatSet: tfStats,
		})
	}

	return &results, nil
}

func (r *simResults) Prettify() (string, error) {
	var str string

	str += r.simResults.prettify()
	str += r.prettifyTfInfo()

	return str, nil
}

func (r *simResults) prettifyTfInfo() (str string) {
	str += "Traffic Flow domain.Results\n"
	str += "====================\n"
	str += r.prettifyTfTable()

	return str
}

func (r *simResults) prettifyTfTable() (str string) {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "T_i"},
			{Align: simpletable.AlignLeft, Text: "No. pkts"},
			{Align: simpletable.AlignLeft, Text: "No. > D_i"},
			{Align: simpletable.AlignLeft, Text: "max"},
			{Align: simpletable.AlignLeft, Text: "mean"},
			{Align: simpletable.AlignLeft, Text: "min"},
			{Align: simpletable.AlignLeft, Text: "D_i"},
		},
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		row := r.prettifyTfRow(r.trafficFlows[i])
		table.Body.Cells = append(table.Body.Cells, row)
	}

	str += table.String()

	return str
}

func (r *simResults) prettifyTfRow(tf tfSim) []*simpletable.Cell {
	row := []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: tf.ID},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsRouted)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.PacketsExceededDeadline)},
		{Align: simpletable.AlignLeft, Text: cleanBestLatency(tf.BestLatency)},
		{Align: simpletable.AlignLeft, Text: cleanMeanLatency(tf.MeanLatency)},
		{Align: simpletable.AlignLeft, Text: cleanWorstLatency(tf.WorstLatency)},
		{Align: simpletable.AlignLeft, Text: strconv.Itoa(tf.Deadline)},
	}

	return row
}

func (r *simResults) OutputCSV(path string) error {
	data := [][]string{
		{
			"Traffic Flow",
			"Packets Routed",
			"Packets Exceeded Deadline",
			"Packets Lost",
			"Best Latency",
			"Mean Latency",
			"Worst Latency",
			"Deadline",
			"Schedulable",
		},
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		data = append(data, []string{
			r.trafficFlows[i].ID,
			strconv.Itoa(r.trafficFlows[i].PacketsRouted),
			strconv.Itoa(r.trafficFlows[i].PacketsExceededDeadline),
			strconv.Itoa(r.trafficFlows[i].PacketsLost),
			cleanBestLatency(r.trafficFlows[i].BestLatency),
			cleanMeanLatency(r.trafficFlows[i].MeanLatency),
			cleanWorstLatency(r.trafficFlows[i].WorstLatency),
			strconv.Itoa(r.trafficFlows[i].Deadline),
			strconv.FormatBool(r.trafficFlows[i].Schedulable()),
		})
	}

	return writeCSV(path, data)
}
