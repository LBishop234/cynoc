package analysis

import (
	"context"

	"main/src/domain"
	"main/src/topology"
)

func Analysis(ctx context.Context, conf domain.SimConfig, top *topology.Topology, trafficFlows []domain.TrafficFlowConfig) (domain.AnalysisResults, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		trafficFlows = sortTfsByPriority(trafficFlows)

		analysisTFs, err := constructAnalysisTfs(top, trafficFlows)
		if err != nil {
			return nil, err
		}

		analysisTFsMap := make(map[string]int, len(analysisTFs))
		for i := 0; i < len(analysisTFs); i++ {
			analysisTFsMap[analysisTFs[i].ID] = i
		}

		analysisTFs, err = basicLatency(ctx, conf, analysisTFs)
		if err != nil {
			return nil, err
		}

		analysisTFs, err = shiBurns(ctx, analysisTFs)
		if err != nil {
			return nil, err
		}

		res := make(domain.AnalysisResults, len(analysisTFs))
		for i := 0; i < len(analysisTFs); i++ {
			res[analysisTFs[i].ID] = domain.TrafficFlowAnalysisSet{
				TrafficFlowConfig:         analysisTFs[i].TrafficFlowConfig,
				Basic:                     analysisTFs[i].Basic,
				ShiAndBurns:               analysisTFs[i].ShiAndBurns,
				DirectInterferenceCount:   analysisTFs[i].DirectInterferenceCount,
				IndirectInterferenceCount: analysisTFs[i].IndirectInterferenceCount,
			}
		}

		return res, nil
	}
}
