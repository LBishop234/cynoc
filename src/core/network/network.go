package network

import (
	"main/src/core/network/components"
	"main/src/domain"
	"main/src/topology"

	"github.com/rs/zerolog"
)

type Network interface {
	NetworkInterfaces() []components.NetworkInterface
	Routers() []components.Router

	NetworkInterfaceMap() map[domain.NodeID]components.NetworkInterface
	RouterMap() map[domain.NodeID]components.Router

	NetworkInterfacesIDMap() map[string]components.NetworkInterface
	RoutersIDMap() map[string]components.Router

	Topology() *topology.Topology

	Cycle(cycle int) error
}

type networkImpl struct {
	netwrkIntfcs []components.NetworkInterface
	routers      []components.Router

	netwrkIntfcMap map[domain.NodeID]components.NetworkInterface
	routerMap      map[domain.NodeID]components.Router

	netwrkIntfcIDMap map[string]components.NetworkInterface
	routerIDMap      map[string]components.Router

	top *topology.Topology

	logger zerolog.Logger
}

func NewNetwork(top *topology.Topology, conf domain.SimConfig, logger zerolog.Logger) (Network, error) {
	routerNodes, err := buildNetwork(top, conf, logger)
	if err != nil {
		logger.Error().Err(err).Msg("error building network")
		return nil, err
	}

	netwrkIntfcs := make([]components.NetworkInterface, len(routerNodes))
	index := 0
	for id := range routerNodes {
		netwrkIntfcs[index] = routerNodes[id].NetworkInterface
		index++
	}

	routers := make([]components.Router, len(routerNodes))
	index = 0
	for id := range routerNodes {
		routers[index] = routerNodes[id].Router
		index++
	}

	netwrkIntfcMap := make(map[domain.NodeID]components.NetworkInterface)
	for id := range routerNodes {
		netwrkIntfcMap[routerNodes[id].NetworkInterface.NodeID()] = routerNodes[id].NetworkInterface
	}

	routerMap := make(map[domain.NodeID]components.Router)
	for id := range routerNodes {
		routerMap[routerNodes[id].Router.NodeID()] = routerNodes[id].Router
	}

	netwrkIntfcIDMap := make(map[string]components.NetworkInterface)
	for id := range routerNodes {
		netwrkIntfcIDMap[routerNodes[id].NetworkInterface.NodeID().ID] = routerNodes[id].NetworkInterface
	}

	routerIDMap := make(map[string]components.Router)
	for id := range routerNodes {
		routerIDMap[routerNodes[id].Router.NodeID().ID] = routerNodes[id].Router
	}

	return &networkImpl{
		netwrkIntfcs: netwrkIntfcs,
		routers:      routers,

		netwrkIntfcMap: netwrkIntfcMap,
		routerMap:      routerMap,

		netwrkIntfcIDMap: netwrkIntfcIDMap,
		routerIDMap:      routerIDMap,

		top: top,
	}, nil
}

func buildNetwork(top *topology.Topology, conf domain.SimConfig, logger zerolog.Logger) (map[string]components.RouterNode, error) {
	logger.Debug().Msg("constructing network from topology")

	routerNodes := make(map[string]components.RouterNode)

	logger.Debug().Msg("creating routers")
	for id := range top.Nodes() {
		node, exists := top.Node(id)
		if !exists {
			logger.Error().Err(domain.ErrInvalidTopology).Str("node_id", id).Msg("node does not exist")
			return nil, domain.ErrInvalidTopology
		}

		rNode, err := components.NewRouterNode(
			components.RouterConfig{
				NodeID:    node.NodeID(),
				SimConfig: conf,
			},
			logger,
		)
		if err != nil {
			logger.Error().Err(err).Str("node_id", node.NodeID().ID).Msg("error creating router")
			return nil, err
		}

		routerNodes[rNode.NodeID().ID] = rNode
	}

	logger.Debug().Msg("connecting routers")
	for id := range top.Edges() {
		edge, exists := top.Edge(id)
		if !exists {
			logger.Error().Err(domain.ErrInvalidTopology).Str("node_id", id).Msg("edge does not exist")
			return nil, domain.ErrInvalidTopology
		}

		aRouterNode, exists := routerNodes[edge.A()]
		if !exists {
			logger.Error().Err(domain.ErrInvalidTopology).Str("edge_id", edge.ID()).Str("node_id", edge.A()).Msg("edge A router does not exist")
			return nil, domain.ErrInvalidTopology
		}

		bRouterNode, exists := routerNodes[edge.B()]
		if !exists {
			logger.Error().Err(domain.ErrInvalidTopology).Str("edge_id", edge.ID()).Str("node_id", edge.B()).Msg("edge B router does not exist")
			return nil, domain.ErrInvalidTopology
		}

		aToB, err := components.NewConnection(conf.MaxPriority, logger)
		if err != nil {
			logger.Error().Err(err).Str("id", edge.ID()).Msg("error creating connection")
			return nil, err
		}
		if err := aRouterNode.Router.RegisterOutputPort(aToB); err != nil {
			logger.Error().Err(err).Str("edge_id", edge.ID()).Str("node_id", aRouterNode.NodeID().ID).Msg("error registering output port")
			return nil, err
		}
		if err := bRouterNode.Router.RegisterInputPort(aToB); err != nil {
			logger.Error().Err(err).Str("edge_id", edge.ID()).Str("node_id", bRouterNode.NodeID().ID).Msg("error registering input port")
			return nil, err
		}

		bToA, err := components.NewConnection(conf.MaxPriority, logger)
		if err != nil {
			logger.Error().Err(err).Str("id", edge.ID()).Msg("error creating connection")
			return nil, err
		}
		if err := bRouterNode.Router.RegisterOutputPort(bToA); err != nil {
			logger.Error().Err(err).Str("edge_id", edge.ID()).Str("node_id", bRouterNode.NodeID().ID).Msg("error registering output port")
			return nil, err
		}
		if err := aRouterNode.Router.RegisterInputPort(bToA); err != nil {
			logger.Error().Err(err).Str("edge_id", edge.ID()).Str("node_id", aRouterNode.NodeID().ID).Msg("error registering input port")
			return nil, err
		}
	}

	logger.Info().Msg("build network from topology")
	return routerNodes, nil
}

func (n *networkImpl) NetworkInterfaces() []components.NetworkInterface {
	return n.netwrkIntfcs
}

func (n *networkImpl) Routers() []components.Router {
	return n.routers
}

func (n *networkImpl) NetworkInterfaceMap() map[domain.NodeID]components.NetworkInterface {
	return n.netwrkIntfcMap
}

func (n *networkImpl) RouterMap() map[domain.NodeID]components.Router {
	return n.routerMap
}

func (n *networkImpl) NetworkInterfacesIDMap() map[string]components.NetworkInterface {
	return n.netwrkIntfcIDMap
}

func (n *networkImpl) RoutersIDMap() map[string]components.Router {
	return n.routerIDMap
}

func (n *networkImpl) Cycle(cycle int) error {
	for i := 0; i < len(n.netwrkIntfcs); i++ {
		if err := n.netwrkIntfcs[i].TransmitPendingPackets(cycle); err != nil {
			n.logger.Error().Err(err).Str("id", n.netwrkIntfcs[i].NodeID().ID).Msg("error transmitting network interface's pending packets")
			return err
		}
	}

	if err := n.cycleRouters(cycle); err != nil {
		n.logger.Error().Err(err).Msg("error cycling routers")
		return err
	}

	for i := 0; i < len(n.netwrkIntfcs); i++ {
		if err := n.netwrkIntfcs[i].HandleArrivingFlits(cycle); err != nil {
			n.logger.Error().Err(err).Str("node_id", n.netwrkIntfcs[i].NodeID().ID).Msg("error handling network interface's arriving flits")
			return err
		}
	}

	return nil
}

func (n *networkImpl) Topology() *topology.Topology {
	return n.top
}

func (n *networkImpl) cycleRouters(cycle int) error {
	for i := 0; i < len(n.routers); i++ {
		n.routers[i].UpdateOutputMap()
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].UpdateOutputPortsCredit(); err != nil {
			n.logger.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error updating router output ports credit")
			return err
		}
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].RouteBufferedFlits(cycle); err != nil {
			n.logger.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error routing buffered flits")
			return err
		}
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].ReadFromInputPorts(cycle); err != nil {
			n.logger.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error reading from router input ports")
			return err
		}
	}

	return nil
}
