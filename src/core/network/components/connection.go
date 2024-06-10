package components

import (
	"main/src/domain"
	"main/src/traffic/packet"

	"github.com/rs/zerolog"
)

type Connection interface {
	flitChannel() chan packet.Flit
	creditChannels() map[int]chan int
	creditChannel(priority int) chan int
	flitBandwidth() int

	GetDstRouter() domain.NodeID
	SetDstRouter(nodeID domain.NodeID)

	GetSrcRouter() domain.NodeID
	SetSrcRouter(nodeID domain.NodeID)
}

type connectionImpl struct {
	flitChan   chan packet.Flit
	creditChan map[int]chan int
	destRouter domain.NodeID
	srcRouter  domain.NodeID
	bandwidth  int
	logger     zerolog.Logger
}

func NewConnection(maxPriority, bandwidth int, logger zerolog.Logger) (*connectionImpl, error) {
	creditChan := make(map[int]chan int, maxPriority)
	for i := 1; i <= maxPriority; i++ {
		creditChan[i] = make(chan int, bandwidth)
	}

	logger.Trace().Int("credit_channels", maxPriority).Int("flit_bandwidth", bandwidth).Msg("new connection")
	return &connectionImpl{
		flitChan:   make(chan packet.Flit, bandwidth),
		creditChan: creditChan,
		bandwidth:  bandwidth,
		logger:     logger.With().Logger(),
	}, nil
}

func (c *connectionImpl) flitChannel() chan packet.Flit {
	return c.flitChan
}

func (c *connectionImpl) creditChannels() map[int]chan int {
	return c.creditChan
}

func (c *connectionImpl) creditChannel(priority int) chan int {
	if _, exists := c.creditChan[priority]; !exists {
		c.creditChan[priority] = make(chan int, 1)
	}

	return c.creditChan[priority]
}

func (c *connectionImpl) flitBandwidth() int {
	return c.bandwidth
}

func (c *connectionImpl) GetDstRouter() domain.NodeID {
	return c.destRouter
}

func (c *connectionImpl) SetDstRouter(nodeID domain.NodeID) {
	c.destRouter = nodeID
}

func (c *connectionImpl) GetSrcRouter() domain.NodeID {
	return c.srcRouter
}

func (c *connectionImpl) SetSrcRouter(nodeID domain.NodeID) {
	c.srcRouter = nodeID
}
