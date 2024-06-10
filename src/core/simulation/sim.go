package simulation

import (
	"context"
	"math"
	"time"

	"main/src/core/network"
	"main/src/domain"
	"main/src/traffic"

	"github.com/rs/zerolog"
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

	rcrds *Records

	logger zerolog.Logger
}

type trafficFlowRoute struct {
	traffic.TrafficFlow
	route domain.Route
}

func newSimulator(network network.Network, trafficFlows []traffic.TrafficFlow, routingAlg domain.RoutingAlgorithm, cycleLimit int, logger zerolog.Logger) (*simulator, error) {
	simulator := &simulator{
		network:      network,
		routingAlg:   routingAlg,
		cycleLimit:   cycleLimit,
		trafficFlows: make([]trafficFlowRoute, len(trafficFlows)),

		rcrds: newRecords(logger),

		logger: logger,
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
			logger.Error().Err(domain.ErrUnknownRoutingAlgorithm).Str("routing_algorithm", string(routingAlg)).Msg("routing algorithm not supported")
			return nil, domain.ErrUnknownRoutingAlgorithm
		}
		if err != nil {
			logger.Error().Err(err).Msg("error calculating router")
			return nil, err
		}

		simulator.trafficFlows[i] = trafficFlowRoute{
			TrafficFlow: trafficFlows[i],
			route:       route,
		}
	}

	return simulator, nil
}

func Simulate(ctx context.Context, network network.Network, trafficFlows []traffic.TrafficFlow, routingAlg domain.RoutingAlgorithm, cycleLimit int, logger zerolog.Logger) (domain.FullResults, error) {
	select {
	case <-ctx.Done():
		return domain.FullResults{}, ctx.Err()
	default:
		simulator, err := newSimulator(network, trafficFlows, routingAlg, cycleLimit, logger)
		if err != nil {
			logger.Error().Err(nil).Msg("error creating simulator")
			return domain.FullResults{}, err
		}

		simDuration, rcrds, err := simulator.runSimulation(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("error running simulation")
			return domain.FullResults{}, err
		}

		return results(cycleLimit, simDuration, rcrds, trafficFlows), nil
	}
}

func (s *simulator) runSimulation(ctx context.Context) (time.Duration, *Records, error) {
	logProgressInterval := int(math.Round(float64(s.cycleLimit) * simProgressMultiple))
	if logProgressInterval > maxProgressInterval {
		logProgressInterval = maxProgressInterval
	}

	s.logger.Info().Msg("starting simulation")

	start := time.Now()

	for c := 0; c < s.cycleLimit; c++ {
		select {
		case <-ctx.Done():
			return 0, nil, ctx.Err()
		default:
			s.logger.Trace().Int("cycle", c).Msg("starting cycle")

			if err := s.releasePackets(c); err != nil {
				s.logger.Error().Err(err).Msg("error releasing packets")
				return 0, nil, err
			}

			if err := s.network.Cycle(c); err != nil {
				s.logger.Error().Err(err).Msg("error cycling network")
				return 0, nil, err
			}

			for i := 0; i < len(s.network.NetworkInterfaces()); i++ {
				pkts := s.network.NetworkInterfaces()[i].PopArrivedPackets(c)
				for i := 0; i < len(pkts); i++ {
					s.rcrds.recordArrivedPacket(c, pkts[i])
				}
			}

			if c > 0 && c%logProgressInterval == 0 {
				s.logger.Info().Int("cycle", c).Int("limit", s.cycleLimit).Msg("simulation progress")
			}

			s.logger.Trace().Int("cycle", c).Msg("cycle completed")
		}
	}

	simDuration := time.Since(start)

	s.logger.Info().Dur("duration_ms", simDuration).Msg("simulation complete")
	return simDuration, s.rcrds, nil
}

func (s *simulator) releasePackets(cycle int) error {
	for i := 0; i < len(s.trafficFlows); i++ {
		released, pkt, periodStartCycle := s.trafficFlows[i].ReleasePacket(cycle, s.trafficFlows[i].TrafficFlow, s.trafficFlows[i].route, s.logger)

		if released {
			if netwrkIntfc, exists := s.network.NetworkInterfaceMap()[pkt.Route()[0]]; exists {
				if err := netwrkIntfc.RoutePacket(cycle, pkt); err != nil {
					s.logger.Error().Err(err).Msg("failed to route packet")
					return err
				}

				s.rcrds.recordTransmittedPacket(periodStartCycle, cycle, pkt)
			} else {
				s.logger.Error().Err(domain.ErrMissingNetworkInterface).Str("network_interface", pkt.Route()[0].Pos.Prettify()).Msg("network interface not found")
				return domain.ErrMissingNetworkInterface
			}
		}
	}

	return nil
}
