package analysis

import (
	"context"

	"main/src/core/analysis/basic"
	"main/src/core/analysis/shiburns"
	"main/src/core/analysis/util"
	"main/src/domain"
	"main/src/topology"
)

func Analysis(ctx context.Context, conf domain.SimConfig, top *topology.Topology, trafficFlows []domain.TrafficFlowConfig) (domain.AnalysisResults, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		trafficFlowMap := make(map[string]domain.TrafficFlowConfig, len(trafficFlows))
		for i := 0; i < len(trafficFlows); i++ {
			trafficFlowMap[trafficFlows[i].ID] = trafficFlows[i]
		}

		tfrs, err := util.ConstructTrafficFlowAndRoutes(top, trafficFlowMap)
		if err != nil {
			return nil, err
		}

		analyses := make(domain.AnalysisResults, len(trafficFlows))
		for key := range trafficFlowMap {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				analyses[key], err = analyseTrafficFlow(conf, tfrs, key)
				if err != nil {
					return nil, err
				}
			}
		}

		return analyses, nil
	}
}

func analyseTrafficFlow(conf domain.SimConfig, tfrs map[string]util.TrafficFlowAndRoute, key string) (domain.TrafficFlowAnalysisSet, error) {
	shiAndBurns, err := shiburns.ShiBurns(conf, tfrs, key)
	if err != nil {
		return domain.TrafficFlowAnalysisSet{}, err
	}

	return domain.TrafficFlowAnalysisSet{
		TrafficFlowConfig:         tfrs[key].TrafficFlowConfig,
		Basic:                     basic.BasicLatency(conf, tfrs[key]),
		ShiAndBurns:               shiAndBurns.Latency,
		DirectInterferenceCount:   shiAndBurns.DirectInterferenceCount,
		IndirectInterferenceCount: shiAndBurns.IndirectInterferenceCount,
	}, nil
}
