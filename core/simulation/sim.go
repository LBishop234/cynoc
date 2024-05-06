package simulation

import (
	"context"
	"math"
	"time"

	"main/core/network"
	"main/domain"
	"main/log"
	"main/traffic"
	"main/traffic/packet"
)

const (
	simProgressMultiple = 0.05
	maxProgressInterval = 100000
)

type simulator struct {
	network      network.Network
	trafficFlows []trafficFlowRoute
	routingAlg   domain.RoutingAlgorithm
	cycleLimit   int

	rcrds *records
}

type trafficFlowRoute struct {
	traffic.TrafficFlow
	route domain.Route
}

func newSimulator(network network.Network, trafficFlows []traffic.TrafficFlow, routingAlg domain.RoutingAlgorithm, cycleLimit int) (*simulator, error) {
	simulator := &simulator{
		network:      network,
		routingAlg:   routingAlg,
		cycleLimit:   cycleLimit,
		trafficFlows: make([]trafficFlowRoute, len(trafficFlows)),

		rcrds: newRecords(),
	}

	for i := 0; i < len(trafficFlows); i++ {
		var route domain.Route
		var err error

		switch routingAlg {
		case domain.XYRouting:
			route, err = network.Topology().XYRoute(
				network.NetworkInterfacesIDMap()[trafficFlows[i].Src()].NodeID(),
				network.NetworkInterfacesIDMap()[trafficFlows[i].Dst()].NodeID(),
			)
		default:
			log.Log.Error().Err(domain.ErrUnknownRoutingAlgorithm).Str("routing_algorithm", string(routingAlg)).Msg("routing algorithm not supported")
			return nil, domain.ErrUnknownRoutingAlgorithm
		}
		if err != nil {
			log.Log.Error().Err(err).Msg("error calculating router")
			return nil, err
		}

		simulator.trafficFlows[i] = trafficFlowRoute{
			TrafficFlow: trafficFlows[i],
			route:       route,
		}
	}

	return simulator, nil
}

func Simulate(ctx context.Context, network network.Network, trafficFlows []traffic.TrafficFlow, routingAlg domain.RoutingAlgorithm, cycleLimit int) (Results, error) {
	select {
	case <-ctx.Done():
		return Results{}, ctx.Err()
	default:
		simulator, err := newSimulator(network, trafficFlows, routingAlg, cycleLimit)
		if err != nil {
			log.Log.Error().Err(nil).Msg("error creating simulator")
			return Results{}, err
		}

		simDuration, rcrds, err := simulator.runSimulation(ctx)
		if err != nil {
			log.Log.Error().Err(err).Msg("error running simulation")
			return Results{}, err
		}

		return results(cycleLimit, simDuration, rcrds, trafficFlows), nil
	}
}

func (s *simulator) runSimulation(ctx context.Context) (time.Duration, *records, error) {
	logProgressInterval := int(math.Round(float64(s.cycleLimit) * simProgressMultiple))
	if logProgressInterval > maxProgressInterval {
		logProgressInterval = maxProgressInterval
	}

	log.Log.Info().Msg("starting simulation")

	start := time.Now()

	for c := 0; c < s.cycleLimit; c++ {
		select {
		case <-ctx.Done():
			return 0, nil, ctx.Err()
		default:
			log.Log.Trace().Int("cycle", c).Msg("starting cycle")

			if err := s.releasePackets(c); err != nil {
				log.Log.Error().Err(err).Msg("error releasing packets")
				return 0, nil, err
			}

			if err := s.network.Cycle(); err != nil {
				log.Log.Error().Err(err).Msg("error cycling network")
				return 0, nil, err
			}

			for i := 0; i < len(s.network.NetworkInterfaces()); i++ {
				pkts := s.network.NetworkInterfaces()[i].PopArrivedPackets()
				for i := 0; i < len(pkts); i++ {
					s.rcrds.recordArrivedPacket(c, pkts[i])
				}
			}

			if c > 0 && c%logProgressInterval == 0 {
				log.Log.Info().Int("cycle", c).Int("limit", s.cycleLimit).Msg("simulation progress")
			}

			log.Log.Trace().Int("cycle", c).Msg("cycle completed")
		}
	}

	simDuration := time.Since(start)

	log.Log.Info().Dur("duration_ms", simDuration).Msg("simulation complete")
	return simDuration, s.rcrds, nil
}

func (s *simulator) releasePackets(cycle int) error {
	for i := 0; i < len(s.trafficFlows); i++ {
		released, periodStartCycle := s.trafficFlows[i].ReleasePacket(cycle)

		if released {
			pkt := packet.NewPacket(
				s.trafficFlows[i].ID(),
				s.trafficFlows[i].Priority(),
				s.trafficFlows[i].Deadline(),
				s.trafficFlows[i].route,
				s.trafficFlows[i].PacketSize(),
			)

			if netwrkIntfc, exists := s.network.NetworkInterfaceMap()[pkt.Route()[0]]; exists {
				if err := netwrkIntfc.RoutePacket(pkt); err != nil {
					log.Log.Error().Err(err).Msg("failed to route packet")
					return err
				}

				s.rcrds.recordTransmittedPacket(periodStartCycle, cycle, pkt)
			} else {
				log.Log.Error().Err(domain.ErrMissingNetworkInterface).Str("network_interface", pkt.Route()[0].Pos.Prettify()).Msg("network interface not found")
				return domain.ErrMissingNetworkInterface
			}
		}
	}

	return nil
}
