package analysis

import (
	"main/src/domain"
	"main/src/topology"
)

type trafficFlowAndRoute struct {
	domain.TrafficFlowConfig
	domain.Route
}

func constructTrafficFlowAndRoutes(top *topology.Topology, trafficFlows map[string]domain.TrafficFlowConfig) (map[string]trafficFlowAndRoute, error) {
	tfrs := make(map[string]trafficFlowAndRoute, len(trafficFlows))

	var err error
	for key := range trafficFlows {
		if tfrs[key], err = newTrafficFlowAndRoute(top, trafficFlows[key]); err != nil {
			return nil, err
		}
	}

	return tfrs, nil
}

func newTrafficFlowAndRoute(top *topology.Topology, trafficFlow domain.TrafficFlowConfig) (trafficFlowAndRoute, error) {
	strRoute, err := trafficFlow.RouteArray()
	if err != nil {
		return trafficFlowAndRoute{}, err
	}

	route, err := top.Route(strRoute)
	if err != nil {
		return trafficFlowAndRoute{}, err
	}

	return trafficFlowAndRoute{trafficFlow, route}, nil
}
