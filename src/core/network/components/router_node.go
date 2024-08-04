package components

import (
	"github.com/rs/zerolog"
)

type RouterNode struct {
	nodeID           string
	Router           Router
	NetworkInterface NetworkInterface
}

func (r *RouterNode) NodeID() string {
	return r.nodeID
}

func NewRouterNode(conf RouterConfig, logger zerolog.Logger) (RouterNode, error) {
	router, err := newRouter(conf, logger)
	if err != nil {
		logger.Error().Err(err).Msg("error creating new router")
		return RouterNode{}, err
	}

	netIntfc, err := newNetworkInterface(conf.NodeID, conf.BufferSize, conf.MaxPriority, logger)
	if err != nil {
		logger.Error().Err(err).Msg("error creating new network interface")
		return RouterNode{}, err
	}

	if err := router.SetNetworkInterface(netIntfc); err != nil {
		logger.Error().Err(err).Str("id", router.NodeID()).Msg("error setting router network interface")
		return RouterNode{}, err
	}

	logger.Trace().Str("id", router.NodeID()).Msg("new router and network interface")
	return RouterNode{
		nodeID:           conf.NodeID,
		Router:           router,
		NetworkInterface: netIntfc,
	}, nil
}
