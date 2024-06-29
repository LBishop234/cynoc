package simulation

import (
	"time"

	"main/src/domain"
	"main/src/traffic"
)

func simResults(cycles int, dur time.Duration, rcrds *Records, trafficFlows []traffic.TrafficFlow) domain.SimResults {
	results := domain.SimResults{
		SimHeadlineResults: domain.SimHeadlineResults{
			Cycles:   cycles,
			Duration: dur,
			StatSet: domain.StatSet{
				PacketsRouted:           rcrds.noTransmitted(),
				PacketsArrived:          rcrds.noArrived(),
				PacketsLost:             rcrds.noLost(),
				PacketsExceededDeadline: rcrds.noExceededDeadline(),
				BestLatency:             rcrds.bestLatency(),
				MeanLatency:             rcrds.meanLatency(),
				WorstLatency:            rcrds.worstLatency(),
			},
		},
		TFStats: make(map[string]domain.StatSet, len(trafficFlows)),
	}

	for i := 0; i < len(trafficFlows); i++ {
		results.TFStats[trafficFlows[i].ID()] = domain.StatSet{
			PacketsRouted:           rcrds.noTransmittedByTF(trafficFlows[i].ID()),
			PacketsArrived:          rcrds.noArrivedByTF(trafficFlows[i].ID()),
			PacketsLost:             rcrds.noLostByTF(trafficFlows[i].ID()),
			PacketsExceededDeadline: rcrds.noExceededDeadlineByTF(trafficFlows[i].ID()),
			BestLatency:             rcrds.bestLatencyByTF(trafficFlows[i].ID()),
			MeanLatency:             rcrds.meanLatencyByTF(trafficFlows[i].ID()),
			WorstLatency:            rcrds.worstLatencyByTF(trafficFlows[i].ID()),
		}
	}

	return results
}
