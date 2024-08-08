package analysis

import (
	"main/src/domain"
	"main/src/topology"
)

type analysisTF struct {
	domain.TrafficFlowAnalysisSet
	domain.Route
	directIntSet   map[string]int
	indirectIntSet map[string]int
}

func constructAnalysisTfs(top *topology.Topology, trafficFlows []domain.TrafficFlowConfig) ([]analysisTF, error) {
	analysisTFs := make([]analysisTF, len(trafficFlows))

	var err error
	for i := 0; i < len(trafficFlows); i++ {
		if analysisTFs[i], err = newAnalysisTF(top, trafficFlows[i]); err != nil {
			return nil, err
		}
	}

	return analysisTFs, nil
}

func newAnalysisTF(top *topology.Topology, trafficFlow domain.TrafficFlowConfig) (analysisTF, error) {
	strRoute, err := trafficFlow.RouteArray()
	if err != nil {
		return analysisTF{}, err
	}

	route, err := top.Route(strRoute)
	if err != nil {
		return analysisTF{}, err
	}

	return analysisTF{
		TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
			TrafficFlowConfig:         trafficFlow,
			Basic:                     -1,
			ShiAndBurns:               -1,
			DirectInterferenceCount:   -1,
			IndirectInterferenceCount: -1,
		},
		Route: route,
	}, nil
}

func sortTfsByPriority(trafficFlows []domain.TrafficFlowConfig) []domain.TrafficFlowConfig {
	tfs := make([]domain.TrafficFlowConfig, len(trafficFlows))
	copy(tfs, trafficFlows)

	for i := 0; i < len(tfs); i++ {
		for j := 0; j < len(tfs)-i-1; j++ {
			if tfs[j].Priority > tfs[j+1].Priority {
				tfs[j], tfs[j+1] = tfs[j+1], tfs[j]
			}
		}
	}

	return tfs
}
