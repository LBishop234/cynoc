package results

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"main/src/domain"
)

func prettifySimHeadlineResults(r domain.SimHeadlineResults) string {
	str := "Simulation domain.Results\n"
	str += "==================\n"
	str += fmt.Sprintf("Cycles: %d\n", r.Cycles)
	str += fmt.Sprintf("Duration (ms): %d\n\n", r.Duration.Milliseconds())
	str += fmt.Sprintf("Packets Routed: %d\n", r.PacketsRouted)
	str += fmt.Sprintf("Packets Exceeded Deadline: %d\n", r.PacketsExceededDeadline)
	str += "\n"
	return str
}

func cleanInt(val int) string {
	if val == math.MaxInt || val == 0 || val == math.MinInt {
		return "-"
	}
	return strconv.Itoa(val)
}

func cleanFloat(val float64) string {
	if math.IsNaN(val) || val == 0 {
		return "-"
	}
	return strconv.FormatFloat(val, 'f', 2, 64)
}

func writeCSV(path string, data [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)

	return writer.WriteAll(data)
}
