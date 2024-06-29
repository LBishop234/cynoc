package results

import (
	"fmt"
	"main/src/domain"
	"strconv"

	"github.com/alexeyco/simpletable"
)

type resultParameter struct {
	name                string
	terminalStr         string
	csvStr              string
	terminalAllowedFlag bool
	reqAnalysisFlag     bool
	value               func(tf tfSimAnalysis) string
}

var parameters = []resultParameter{
	{
		name:                "Traffic Flow ID",
		terminalStr:         "T_i",
		csvStr:              "TF_ID",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return tf.ID },
	},
	{
		name:                "Direct Interference Count",
		terminalStr:         "|S^D_i|",
		csvStr:              "Direct_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.DirectInterferenceCount) },
	},
	{
		name:                "Indirect Interference Count",
		terminalStr:         "|S^I_i|",
		csvStr:              "Indirect_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.IndirectInterferenceCount) },
	},
	{
		name:                "Number of Packets Routed",
		terminalStr:         "No. pkts",
		csvStr:              "Num_Packets_Routed",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.PacketsRouted) },
	},
	{
		name:                "Number of Packets Exceeded Deadline",
		terminalStr:         "No. > D_i",
		csvStr:              "Num_Packets_Exceeded_Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.PacketsExceededDeadline) },
	},
	{
		name:                "Best Latency",
		terminalStr:         "min",
		csvStr:              "Min_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanInt(tf.BestLatency) },
	},
	{
		name:                "Mean Latency",
		terminalStr:         "mean",
		csvStr:              "Mean_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanFloat(tf.MeanLatency) },
	},
	{
		name:                "Worst Latency",
		terminalStr:         "max",
		csvStr:              "Max_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanInt(tf.WorstLatency) },
	},
	{
		name:                "Deadline",
		terminalStr:         "D_i",
		csvStr:              "Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Deadline) },
	},
	{
		name:                "Simulation Schedulable",
		terminalStr:         "Schedulable",
		csvStr:              "Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return fmt.Sprint(tf.Schedulable()) },
	},
	{
		name:                "Jitter",
		terminalStr:         "J^R_i",
		csvStr:              "Jitter",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter) },
	},
	{
		name:                "Jitter + Maximum Basic Network Latency",
		terminalStr:         "J^R_i + C_i",
		csvStr:              "Jitter_Plus_Basic",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter + tf.Basic) },
	},
	{
		name:                "Jitter + Shi & Burns Network Latency",
		terminalStr:         "J^R_i + R_i",
		csvStr:              "Jitter_Plus_Shi_And_Burns",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter + tf.ShiAndBurns) },
	},
	{
		name:                "Shi & Burns Schedulable",
		terminalStr:         "S&B Schedulable",
		csvStr:              "Shi_Burns_Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return fmt.Sprint(tf.AnalysisSchedulable()) },
	},
}

type Results interface {
	Prettify() (string, error)
	OutputCSV(path string) error
}

type simResults struct {
	domain.SimResults
	trafficFlows []tfSim
}

type simAnalaysisResults struct {
	domain.SimResults
	trafficFlows []tfSimAnalysis
}

type tfSim struct {
	ID       string
	Deadline int
	domain.StatSet
}

type tfSimAnalysis struct {
	tfSim
	domain.TrafficFlowAnalysisSet
	AnalysisHolds bool
}

func NewResults(simRes domain.SimResults, tfOrder []domain.TrafficFlowConfig) (Results, error) {
	var results simResults

	results.SimResults = simRes

	for i := 0; i < len(tfOrder); i++ {
		tfStats, exists := simRes.TFStats[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		results.trafficFlows = append(results.trafficFlows, tfSim{
			ID:       tfOrder[i].ID,
			Deadline: tfOrder[i].Deadline,
			StatSet:  tfStats,
		})
	}

	return &results, nil
}

func NewResultsWithAnalysis(simRes domain.SimResults, analyses domain.AnalysisResults, tfOrder []domain.TrafficFlowConfig) (Results, error) {
	var results simAnalaysisResults

	results.SimResults = simRes

	for i := 0; i < len(tfOrder); i++ {
		tfSimStats, exists := simRes.TFStats[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		tfAnalysis, exists := analyses[tfOrder[i].ID]
		if !exists {
			return nil, domain.ErrMissingTrafficFlow
		}

		analysisHolds := true
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

func (r *simResults) Prettify() (string, error) {
	str := prettifySimResults(r.SimHeadlineResults)

	str += "Traffic Flow domain.Results\n"
	str += "====================\n"

	table := simpletable.New()

	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{}}
	for i := 0; i < len(parameters); i++ {
		if parameters[i].terminalAllowedFlag && !parameters[i].reqAnalysisFlag {
			table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Align: simpletable.AlignLeft, Text: parameters[i].terminalStr})
		}
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		row := []*simpletable.Cell{}
		for p := 0; p < len(parameters); p++ {
			if parameters[p].terminalAllowedFlag && !parameters[p].reqAnalysisFlag {
				row = append(row, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  parameters[p].value(tfSimAnalysis{tfSim: r.trafficFlows[i]}),
				})
			}
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	str += table.String()

	return str, nil
}

func (r *simAnalaysisResults) Prettify() (string, error) {
	str := prettifySimResults(r.SimHeadlineResults)

	str += "Traffic Flow domain.Results\n"
	str += "====================\n"

	table := simpletable.New()

	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{}}
	for i := 0; i < len(parameters); i++ {
		if parameters[i].terminalAllowedFlag {
			table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Align: simpletable.AlignLeft, Text: parameters[i].terminalStr})
		}
	}

	for i := 0; i < len(r.trafficFlows); i++ {
		row := []*simpletable.Cell{}
		for p := 0; p < len(parameters); p++ {
			if parameters[p].terminalAllowedFlag {
				row = append(row, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  parameters[p].value(r.trafficFlows[i]),
				})
			}
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	str += table.String()

	return str, nil
}

func (r *simResults) OutputCSV(path string) error {
	data := [][]string{}

	header := []string{}
	for i := 0; i < len(parameters); i++ {
		if !parameters[i].reqAnalysisFlag {
			header = append(header, parameters[i].csvStr)
		}
	}
	data = append(data, header)

	for i := 0; i < len(r.trafficFlows); i++ {
		row := []string{}
		for p := 0; p < len(parameters); p++ {
			if !parameters[p].reqAnalysisFlag {
				row = append(row, parameters[p].value(tfSimAnalysis{tfSim: r.trafficFlows[i]}))
			}
		}
		data = append(data, row)
	}

	return writeCSV(path, data)
}

func (r *simAnalaysisResults) OutputCSV(path string) error {
	data := [][]string{}

	header := []string{}
	for i := 0; i < len(parameters); i++ {
		header = append(header, parameters[i].csvStr)
	}
	data = append(data, header)

	for i := 0; i < len(r.trafficFlows); i++ {
		row := []string{}
		for p := 0; p < len(parameters); p++ {
			row = append(row, parameters[p].value(r.trafficFlows[i]))
		}
		data = append(data, row)
	}

	return writeCSV(path, data)
}
