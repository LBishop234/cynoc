package shiburns

import (
	"math"

	"main/core/analysis/basic"
	"main/core/analysis/util"
	"main/domain"
)

type ShiBurnsResults struct {
	Latency                   int
	DirectInterferenceCount   int
	IndirectInterferenceCount int
}

func ShiBurns(conf domain.SimConfig, tfrs map[string]util.TrafficFlowAndRoute, tfKey string) (ShiBurnsResults, error) {
	dIntSet, iIntSet := findInterferenceSets(tfrs, tfKey)

	prev := basic.BasicLatency(conf, tfrs[tfKey])
	for {
		interference := 0
		for dIntKey := range dIntSet {
			ji, err := interferenceJitter(conf, tfrs, dIntKey)
			if err != nil {
				return ShiBurnsResults{}, err
			}

			a := int(math.Ceil(float64(prev+tfrs[dIntKey].Jitter+ji) / float64(tfrs[dIntKey].Period)))
			b := basic.BasicLatency(conf, tfrs[dIntKey])
			interference += a * b
		}

		current := interference + basic.BasicLatency(conf, tfrs[tfKey])

		if current > tfrs[tfKey].Deadline || current == prev {
			return ShiBurnsResults{
					Latency:                   current,
					DirectInterferenceCount:   len(dIntSet),
					IndirectInterferenceCount: len(iIntSet),
				},
				nil
		}

		prev = current
	}
}

func filterByPriority(trafficFlows map[string]util.TrafficFlowAndRoute, priority int) map[string]util.TrafficFlowAndRoute {
	filteredTFs := make(map[string]util.TrafficFlowAndRoute)
	for key := range trafficFlows {
		if trafficFlows[key].Priority <= priority {
			filteredTFs[key] = trafficFlows[key]
		}
	}

	return filteredTFs
}

func findInterferenceSets(tfrs map[string]util.TrafficFlowAndRoute, tfKey string) (directIntMap, indirectIntSet map[string]util.TrafficFlowAndRoute) {
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

	indirectIntSet = make(map[string]util.TrafficFlowAndRoute)
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

func findDirectInterferenceSet(tfrs map[string]util.TrafficFlowAndRoute, tfKey string) map[string]util.TrafficFlowAndRoute {
	tfrs = filterByPriority(tfrs, tfrs[tfKey].Priority)

	dIntSet := make(map[string]util.TrafficFlowAndRoute)
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

func interferenceJitter(conf domain.SimConfig, tfrs map[string]util.TrafficFlowAndRoute, tfKey string) (int, error) {
	r, err := ShiBurns(conf, tfrs, tfKey)
	if err != nil {
		return 0, err
	}

	return r.Latency - basic.BasicLatency(conf, tfrs[tfKey]), nil
}
