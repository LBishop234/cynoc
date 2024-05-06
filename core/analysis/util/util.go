package util

import (
	"main/domain"
	"main/topology"
)

const (
	NoAdditionalFlits       = 2
	LinkBandwidthFlitFactor = 1
)

type TrafficFlowAndRoute struct {
	domain.TrafficFlowConfig
	domain.Route
}

func ConstructTrafficFlowAndRoutes(top *topology.Topology, trafficFlows map[string]domain.TrafficFlowConfig) (map[string]TrafficFlowAndRoute, error) {
	tfrs := make(map[string]TrafficFlowAndRoute, len(trafficFlows))

	var err error
	for key := range trafficFlows {
		if tfrs[key], err = NewTrafficFlowAndRoute(top, trafficFlows[key]); err != nil {
			return nil, err
		}
	}

	return tfrs, nil
}

func NewTrafficFlowAndRoute(top *topology.Topology, trafficFlow domain.TrafficFlowConfig) (TrafficFlowAndRoute, error) {
	src, exists := top.Node(trafficFlow.Src)
	if !exists {
		return TrafficFlowAndRoute{}, domain.ErrInvalidTopology
	}

	dst, exists := top.Node(trafficFlow.Dst)
	if !exists {
		return TrafficFlowAndRoute{}, domain.ErrInvalidTopology
	}

	route, err := top.XYRoute(src.NodeID(), dst.NodeID())
	if err != nil {
		return TrafficFlowAndRoute{}, err
	}

	return TrafficFlowAndRoute{trafficFlow, route}, nil
}
