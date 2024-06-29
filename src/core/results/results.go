package results

import (
	"main/src/domain"

	"github.com/alexeyco/simpletable"
)

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
	str := prettifySimHeadlineResults(r.SimHeadlineResults)

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
	str := prettifySimHeadlineResults(r.SimHeadlineResults)

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
