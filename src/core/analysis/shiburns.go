package analysis

import (
	"context"
	"math"

	"main/src/domain"
)

// Assumes analysisTFs are sorted by priority & Basic latency has been calculated.
func shiBurns(ctx context.Context, analysisTFs []analysisTF) ([]analysisTF, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		analysisTFs = findIntereferenceSets(analysisTFs)

		for i := 0; i < len(analysisTFs); i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				current := analysisTFs[i].Basic
				prev := 0

				for current != prev && current < analysisTFs[i].Deadline {
					prev = current

					interference := 0
					for _, dIntIndex := range analysisTFs[i].directIntSet {
						JIi := analysisTFs[dIntIndex].ShiAndBurns - analysisTFs[dIntIndex].Basic

						x := int(math.Ceil((float64(prev + analysisTFs[dIntIndex].Jitter + JIi)) / float64(analysisTFs[dIntIndex].Period)))
						interference += x * analysisTFs[dIntIndex].Basic
					}

					current = interference + analysisTFs[i].Basic
				}

				analysisTFs[i].ShiAndBurns = current
			}
		}

		return analysisTFs, nil
	}
}

// Assumes that the analysisTFs are sorted by priority.
func findIntereferenceSets(analysisTFs []analysisTF) []analysisTF {
	for i := 0; i < len(analysisTFs); i++ {
		analysisTFs[i].directIntSet = make(map[string]int)
		for j := i - 1; j >= 0; j-- {
			if intersectingRoutes(analysisTFs[i].Route, analysisTFs[j].Route) {
				analysisTFs[i].directIntSet[analysisTFs[j].ID] = j
			}
		}
		analysisTFs[i].DirectInterferenceCount = len(analysisTFs[i].directIntSet)
	}

	for i := 0; i < len(analysisTFs); i++ {
		optionsMap := make(map[string]int)
		for _, dIndex := range analysisTFs[i].directIntSet {
			for iKey, iIndex := range analysisTFs[dIndex].directIntSet {
				optionsMap[iKey] = iIndex
			}
			for iKey, iIndex := range analysisTFs[dIndex].indirectIntSet {
				optionsMap[iKey] = iIndex
			}
		}

		analysisTFs[i].indirectIntSet = make(map[string]int)
		for oKey, oIndex := range optionsMap {
			if _, exists := analysisTFs[i].directIntSet[oKey]; !exists {
				analysisTFs[i].indirectIntSet[oKey] = oIndex
			}
		}
		analysisTFs[i].IndirectInterferenceCount = len(analysisTFs[i].indirectIntSet)
	}

	return analysisTFs
}

func intersectingRoutes(rA, rB domain.Route) bool {
	if len(rA) == 0 || len(rB) == 0 {
		return false
	}

	if rA[0] == rB[0] {
		return true
	}

	if rA[len(rA)-1] == rB[len(rB)-1] {
		return true
	}

	for i := 0; i < len(rA)-1; i++ {
		for j := 0; j < len(rB)-1; j++ {
			if rA[i] == rB[j] && rA[i+1] == rB[j+1] {
				return true
			}
		}
	}

	return false
}
