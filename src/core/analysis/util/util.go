package util

import (
	"main/src/domain"
	"main/src/topology"
)

const (
	NoAdditionalFlits = 2
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
	strRoute, err := trafficFlow.RouteArray()
	if err != nil {
		return TrafficFlowAndRoute{}, err
	}

	route, err := top.Route(strRoute)
	if err != nil {
		return TrafficFlowAndRoute{}, err
	}

	return TrafficFlowAndRoute{trafficFlow, route}, nil
}
