package domain

import "time"

type Results interface {
	Prettify() (string, error)
	OutputCSV(path string) error
}

type FullResults struct {
	SimResults SimResults
	TFStats    map[string]StatSet
}

type SimResults struct {
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
