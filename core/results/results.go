package results

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"main/core/analysis"
	"main/core/simulation"
)

type Results interface {
	Prettify() (string, error)
	OutputCSV(path string) error
}

type localSimResults simulation.SimResults

func (r *localSimResults) prettify() (str string) {
	str += "Simulation Results\n"
	str += "==================\n"
	str += fmt.Sprintf("Cycles: %d\n", r.Cycles)
	str += fmt.Sprintf("Duration (ms): %d\n\n", r.Duration.Milliseconds())
	// str += fmt.Sprintf("Duration (s): %.2f\n\n", float64(r.Duration.Milliseconds())/1000)
	str += fmt.Sprintf("Packets Routed: %d\n", r.PacketsRouted)
	str += fmt.Sprintf("Packets Arrived: %d\n", r.PacketsArrived)
	str += fmt.Sprintf("Packets Exceeded Deadline: %d\n", r.PacketsExceededDeadline)
	str += fmt.Sprintf("Packets Lost: %d\n", r.PacketsRouted-r.PacketsArrived)
	str += "\n"
	return str
}

type tfSim struct {
	ID       string
	Deadline int
	simulation.TrafficFlowStatSet
}

type tfSimAnalysis struct {
	tfSim
	analysis.TrafficFlowAnalysisSet
	AnalysisHolds bool
}

func cleanBestLatency(val int) string {
	if val == math.MaxInt {
		return "-"
	}
	return strconv.Itoa(val)
}

func cleanMeanLatency(val float64) string {
	if math.IsNaN(val) || val == 0 {
		return "-"
	}
	return strconv.FormatFloat(val, 'f', 2, 64)
}

func cleanWorstLatency(val int) string {
	if val == math.MinInt {
		return "-"
	}
	return strconv.Itoa(val)
}

func writeCSV(path string, data [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)

	return writer.WriteAll(data)
}
