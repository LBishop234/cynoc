package simulation

import (
	"time"

	"main/src/domain"
	"main/src/traffic"
)

func results(cycles int, dur time.Duration, rcrds *Records, trafficFlows []traffic.TrafficFlow) domain.FullResults {
	results := domain.FullResults{
		SimResults: domain.SimResults{
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
		results.TFStats[trafficFlows[i].ID()] = newTFStatSet(rcrds, trafficFlows[i].ID())
	}

	return results
}

func newTFStatSet(rcrds *Records, tfID string) domain.StatSet {
	return domain.StatSet{
		PacketsRouted:           rcrds.noTransmittedByTF(tfID),
		PacketsArrived:          rcrds.noArrivedByTF(tfID),
		PacketsLost:             rcrds.noLostByTF(tfID),
		PacketsExceededDeadline: rcrds.noExceededDeadlineByTF(tfID),
		BestLatency:             rcrds.bestLatencyByTF(tfID),
		MeanLatency:             rcrds.meanLatencyByTF(tfID),
		WorstLatency:            rcrds.worstLatencyByTF(tfID),
	}
}
