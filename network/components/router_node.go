package components

import (
	"main/core/log"
	"main/domain"
)

type RouterNode struct {
	nodeID           domain.NodeID
	Router           Router
	NetworkInterface NetworkInterface
}

func (r *RouterNode) NodeID() domain.NodeID {
	return r.nodeID
}

func NewRouterNode(conf RouterConfig) (RouterNode, error) {
	router, err := newRouter(conf)
	if err != nil {
		log.Log.Error().Err(err).Msg("error creating new router")
		return RouterNode{}, err
	}

	netIntfc, err := newNetworkInterface(conf.NodeID, conf.BufferSize, conf.FlitSize, conf.MaxPriority)
	if err != nil {
		log.Log.Error().Err(err).Msg("error creating new network interface")
		return RouterNode{}, err
	}

	if err := router.SetNetworkInterface(netIntfc); err != nil {
		log.Log.Error().Err(err).Str("id", router.NodeID().ID).Msg("error setting router network interface")
		return RouterNode{}, err
	}

	log.Log.Trace().Str("id", router.NodeID().ID).Msg("new router and network interface")
	return RouterNode{
		nodeID:           conf.NodeID,
		Router:           router,
		NetworkInterface: netIntfc,
	}, nil
}
