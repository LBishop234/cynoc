package components

import (
	"main/log"
	"main/src/domain"
	"main/src/traffic/packet"
)

type Connection interface {
	flitChannel() chan packet.Flit
	creditChannels() map[int]chan int
	creditChannel(priority int) chan int

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
}

func NewConnection(maxPriority int) (*connectionImpl, error) {
	creditChan := make(map[int]chan int, maxPriority)
	for i := 0; i <= maxPriority; i++ {
		creditChan[i] = make(chan int, 1)
	}

	log.Log.Trace().Msg("new connection")
	return &connectionImpl{
		flitChan:   make(chan packet.Flit, 1),
		creditChan: creditChan,
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
