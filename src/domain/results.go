package domain

import "time"

type SimResults struct {
	SimHeadlineResults SimHeadlineResults
	TFStats            map[string]StatSet
}

type SimHeadlineResults struct {
	Cycles   int
	Duration time.Duration
	StatSet
}

type StatSet struct {
	PacketsRouted           int     `csv:"PacketsRouted"`
	PacketsArrived          int     `csv:"PacketsArrived"`
	PacketsLost             int     `csv:"PacketLost"`
	PacketsExceededDeadline int     `csv:"PacketsExceededDeadline"`
	BestLatency             int     `csv:"BestLatency"`
	MeanLatency             float64 `csv:"MeanLatency"`
	WorstLatency            int     `csv:"WorstLatency"`
}

func (s *StatSet) Schedulable() bool {
	return s.PacketsExceededDeadline == 0
}

type AnalysisResults map[string]TrafficFlowAnalysisSet

func (r AnalysisResults) AnalysesSchedulable() bool {
	for _, tf := range r {
		if !tf.AnalysisSchedulable() {
			return false
		}
	}

	return true
}

type TrafficFlowAnalysisSet struct {
	TrafficFlowConfig
	Basic                     int
	ShiAndBurns               int
	DirectInterferenceCount   int
	IndirectInterferenceCount int
}

func (a TrafficFlowAnalysisSet) AnalysisSchedulable() bool {
	return (a.Jitter + a.ShiAndBurns) < a.Deadline
}
