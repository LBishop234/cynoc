package analysis

import (
	"math"

	"main/src/domain"
)

type shiBurnsResults struct {
	Latency                   int
	DirectInterferenceCount   int
	IndirectInterferenceCount int
}

func shiBurns(conf domain.SimConfig, tfrs map[string]trafficFlowAndRoute, tfKey string) (shiBurnsResults, error) {
	dIntSet, iIntSet := findInterferenceSets(tfrs, tfKey)

	prev := basicLatency(conf, tfrs[tfKey])
	for {
		interference := 0
		for dIntKey := range dIntSet {
			ji, err := interferenceJitter(conf, tfrs, dIntKey)
			if err != nil {
				return shiBurnsResults{}, err
			}

			a := int(math.Ceil(float64(prev+tfrs[dIntKey].Jitter+ji) / float64(tfrs[dIntKey].Period)))
			b := basicLatency(conf, tfrs[dIntKey])
			interference += a * b
		}

		current := interference + basicLatency(conf, tfrs[tfKey])

		if current > tfrs[tfKey].Deadline || current == prev {
			return shiBurnsResults{
					Latency:                   current,
					DirectInterferenceCount:   len(dIntSet),
					IndirectInterferenceCount: len(iIntSet),
				},
				nil
		}

		prev = current
	}
}

func filterByPriority(trafficFlows map[string]trafficFlowAndRoute, priority int) map[string]trafficFlowAndRoute {
	filteredTFs := make(map[string]trafficFlowAndRoute)
	for key := range trafficFlows {
		if trafficFlows[key].Priority <= priority {
			filteredTFs[key] = trafficFlows[key]
		}
	}

	return filteredTFs
}

func findInterferenceSets(tfrs map[string]trafficFlowAndRoute, tfKey string) (directIntMap, indirectIntSet map[string]trafficFlowAndRoute) {
	tfrs = filterByPriority(tfrs, tfrs[tfKey].Priority)

	directIntMap = findDirectInterferenceSet(tfrs, tfKey)

	possibleIIntSet := make([]string, 0, len(tfrs))
	for dIntKey := range directIntMap {
		if dIntKey != tfKey {
			keySet := findDirectInterferenceSet(tfrs, dIntKey)
			for subKey := range keySet {
				possibleIIntSet = append(possibleIIntSet, subKey)
			}
		}
	}

	indirectIntSet = make(map[string]trafficFlowAndRoute)
	for i := 0; i < len(possibleIIntSet); i++ {
		key := possibleIIntSet[i]

		flag := false
		for dIntKey := range directIntMap {
			if key == dIntKey {
				flag = true
			}
		}

		if !flag {
			indirectIntSet[key] = tfrs[key]
		}
	}

	return directIntMap, indirectIntSet
}

func findDirectInterferenceSet(tfrs map[string]trafficFlowAndRoute, tfKey string) map[string]trafficFlowAndRoute {
	tfrs = filterByPriority(tfrs, tfrs[tfKey].Priority)

	dIntSet := make(map[string]trafficFlowAndRoute)
	for key := range tfrs {
		if key != tfKey && intersectingRoutes(tfrs[key].Route, tfrs[tfKey].Route) {
			dIntSet[key] = tfrs[key]
		}
	}

	return dIntSet
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

func interferenceJitter(conf domain.SimConfig, tfrs map[string]trafficFlowAndRoute, tfKey string) (int, error) {
	r, err := shiBurns(conf, tfrs, tfKey)
	if err != nil {
		return 0, err
	}

	return r.Latency - basicLatency(conf, tfrs[tfKey]), nil
}
