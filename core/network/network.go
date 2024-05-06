package network

import (
	"main/core/network/components"
	"main/domain"
	"main/log"
	"main/topology"
)

type Network interface {
	NetworkInterfaces() []components.NetworkInterface
	Routers() []components.Router

	NetworkInterfaceMap() map[domain.NodeID]components.NetworkInterface
	RouterMap() map[domain.NodeID]components.Router

	NetworkInterfacesIDMap() map[string]components.NetworkInterface
	RoutersIDMap() map[string]components.Router

	Topology() *topology.Topology

	Cycle() error
}

type networkImpl struct {
	netwrkIntfcs []components.NetworkInterface
	routers      []components.Router

	netwrkIntfcMap map[domain.NodeID]components.NetworkInterface
	routerMap      map[domain.NodeID]components.Router

	netwrkIntfcIDMap map[string]components.NetworkInterface
	routerIDMap      map[string]components.Router

	top *topology.Topology
}

func NewNetwork(top *topology.Topology, conf domain.SimConfig) (Network, error) {
	routerNodes, err := buildNetwork(top, conf)
	if err != nil {
		log.Log.Error().Err(err).Msg("error building network")
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

func buildNetwork(top *topology.Topology, conf domain.SimConfig) (map[string]components.RouterNode, error) {
	log.Log.Debug().Msg("constructing network from topology")

	routerNodes := make(map[string]components.RouterNode)

	log.Log.Debug().Msg("creating routers")
	for id := range top.Nodes() {
		node, exists := top.Node(id)
		if !exists {
			log.Log.Error().Err(domain.ErrInvalidTopology).Str("id", id).Msg("node does not exist")
			return nil, domain.ErrInvalidTopology
		}

		rNode, err := components.NewRouterNode(components.RouterConfig{
			NodeID:    node.NodeID(),
			SimConfig: conf,
		})
		if err != nil {
			log.Log.Error().Err(err).Str("id", node.NodeID().ID).Msg("error creating router")
			return nil, err
		}

		routerNodes[rNode.NodeID().ID] = rNode
	}

	log.Log.Debug().Msg("connecting routers")
	for id := range top.Edges() {
		edge, exists := top.Edge(id)
		if !exists {
			log.Log.Error().Err(domain.ErrInvalidTopology).Str("id", id).Msg("edge does not exist")
			return nil, domain.ErrInvalidTopology
		}

		aRouterNode, exists := routerNodes[edge.A()]
		if !exists {
			log.Log.Error().Err(domain.ErrInvalidTopology).Str("edge_id", edge.ID()).Str("router_id", edge.A()).Msg("edge A router does not exist")
			return nil, domain.ErrInvalidTopology
		}

		bRouterNode, exists := routerNodes[edge.B()]
		if !exists {
			log.Log.Error().Err(domain.ErrInvalidTopology).Str("edge_id", edge.ID()).Str("router_id", edge.B()).Msg("edge B router does not exist")
			return nil, domain.ErrInvalidTopology
		}

		aToB, err := components.NewConnection(conf.MaxPriority)
		if err != nil {
			log.Log.Error().Err(err).Str("id", edge.ID()).Msg("error creating connection")
			return nil, err
		}
		if err := aRouterNode.Router.RegisterOutputPort(aToB); err != nil {
			log.Log.Error().Err(err).Str("edge_id", edge.ID()).Str("router_id", aRouterNode.NodeID().ID).Msg("error registering output port")
			return nil, err
		}
		if err := bRouterNode.Router.RegisterInputPort(aToB); err != nil {
			log.Log.Error().Err(err).Str("edge_id", edge.ID()).Str("router_id", bRouterNode.NodeID().ID).Msg("error registering input port")
			return nil, err
		}

		bToA, err := components.NewConnection(conf.MaxPriority)
		if err != nil {
			log.Log.Error().Err(err).Str("id", edge.ID()).Msg("error creating connection")
			return nil, err
		}
		if err := bRouterNode.Router.RegisterOutputPort(bToA); err != nil {
			log.Log.Error().Err(err).Str("edge_id", edge.ID()).Str("router_id", bRouterNode.NodeID().ID).Msg("error registering output port")
			return nil, err
		}
		if err := aRouterNode.Router.RegisterInputPort(bToA); err != nil {
			log.Log.Error().Err(err).Str("edge_id", edge.ID()).Str("router_id", aRouterNode.NodeID().ID).Msg("error registering input port")
			return nil, err
		}
	}

	log.Log.Info().Msg("build network from topology")
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

func (n *networkImpl) Cycle() error {
	for i := 0; i < len(n.netwrkIntfcs); i++ {
		if err := n.netwrkIntfcs[i].TransmitPendingPackets(); err != nil {
			log.Log.Error().Err(err).Str("id", n.netwrkIntfcs[i].NodeID().ID).Msg("error transmitting network interface's pending packets")
			return err
		}
	}

	if err := n.cycleRouters(); err != nil {
		log.Log.Error().Err(err).Msg("error cycling routers")
		return err
	}

	for i := 0; i < len(n.netwrkIntfcs); i++ {
		if err := n.netwrkIntfcs[i].HandleArrivingFlits(); err != nil {
			log.Log.Error().Err(err).Str("id", n.netwrkIntfcs[i].NodeID().ID).Msg("error handling network interface's arriving flits")
			return err
		}
	}

	return nil
}

func (n *networkImpl) Topology() *topology.Topology {
	return n.top
}

func (n *networkImpl) cycleRouters() error {
	for i := 0; i < len(n.routers); i++ {
		n.routers[i].UpdateOutputMap()
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].UpdateOutputPortsCredit(); err != nil {
			log.Log.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error updating router output ports credit")
			return err
		}
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].RouteBufferedFlits(); err != nil {
			log.Log.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error routing buffered flits")
			return err
		}
	}

	for i := 0; i < len(n.routers); i++ {
		if err := n.routers[i].ReadFromInputPorts(); err != nil {
			log.Log.Error().Err(err).Str("id", n.routers[i].NodeID().ID).Msg("error reading from router input ports")
			return err
		}
	}

	log.Log.Trace().Msg("performed network cycle for network routers")
	return nil
}
