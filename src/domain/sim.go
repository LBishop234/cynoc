package domain

import "time"

type SimConfig struct {
	CycleLimit       int              `yaml:"cycle_limit" json:"cycle_limit"`
	RoutingAlgorithm RoutingAlgorithm `yaml:"routing_algorithm" json:"routing_algorithm"`
	MaxPriority      int              `yaml:"max_priority" json:"max_priority"`
	FlitSize         int              `yaml:"flit_size" json:"flit_size"`
	BufferSize       int              `yaml:"buffer_size" json:"buffer_size"`
	LinkBandwidth    int              `yaml:"link_bandwidth" json:"link_bandwidth"`
	ProcessingDelay  int              `yaml:"processing_delay" json:"processing_delay"`
}

type Results struct {
	SimResults SimResults
	TFStats    map[string]TrafficFlowStatSet
}

type SimResults struct {
	Cycles   int
	Duration time.Duration
	StatSet
}

type TrafficFlowStatSet struct {
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

func (s *TrafficFlowStatSet) Schedulable() bool {
	return s.PacketsExceededDeadline == 0
}
