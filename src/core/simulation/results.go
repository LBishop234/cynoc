package simulation

import (
	"time"

	"main/src/traffic"
)

type StatSet struct {
	PacketsRouted           int     `csv:"PacketsRouted"`
	PacketsArrived          int     `csv:"PacketsArrived"`
	PacketsLost             int     `csv:"PacketLost"`
	PacketsExceededDeadline int     `csv:"PacketsExceededDeadline"`
	BestLatency             int     `csv:"BestLatency"`
	MeanLatency             float64 `csv:"MeanLatency"`
	WorstLatency            int     `csv:"WorstLatency"`
}

type TrafficFlowStatSet struct {
	StatSet
}

type SimResults struct {
	Cycles   int
	Duration time.Duration
	StatSet
}

type Results struct {
	SimResults SimResults
	TFStats    map[string]TrafficFlowStatSet
}

func results(cycles int, dur time.Duration, rcrds *records, trafficFlows []traffic.TrafficFlow) Results {
	results := Results{
		SimResults: SimResults{
			Cycles:   cycles,
			Duration: dur,
			StatSet: StatSet{
				PacketsRouted:           rcrds.noTransmitted(),
				PacketsArrived:          rcrds.noArrived(),
				PacketsLost:             rcrds.noLost(),
				PacketsExceededDeadline: rcrds.noExceededDeadline(),
				BestLatency:             rcrds.bestLatency(),
				MeanLatency:             rcrds.meanLatency(),
				WorstLatency:            rcrds.worstLatency(),
			},
		},
		TFStats: make(map[string]TrafficFlowStatSet, len(trafficFlows)),
	}

	for i := 0; i < len(trafficFlows); i++ {
		results.TFStats[trafficFlows[i].ID()] = newTFStatSet(rcrds, trafficFlows[i].ID())
	}

	return results
}

func newTFStatSet(rcrds *records, tfID string) TrafficFlowStatSet {
	return TrafficFlowStatSet{
		StatSet: StatSet{
			PacketsRouted:           rcrds.noTransmittedByTF(tfID),
			PacketsArrived:          rcrds.noArrivedByTF(tfID),
			PacketsLost:             rcrds.noLostByTF(tfID),
			PacketsExceededDeadline: rcrds.noExceededDeadlineByTF(tfID),
			BestLatency:             rcrds.bestLatencyByTF(tfID),
			MeanLatency:             rcrds.meanLatencyByTF(tfID),
			WorstLatency:            rcrds.worstLatencyByTF(tfID),
		},
	}
}

func (s *TrafficFlowStatSet) Schedulable() bool {
	return s.PacketsExceededDeadline == 0
}
