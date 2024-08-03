package components

import (
	"main/src/traffic/packet"

	"github.com/rs/zerolog"
)

type Connection interface {
	flitChannel() chan packet.Flit
	creditChannels() map[int]chan int
	creditChannel(priority int) chan int

	GetDstRouter() string
	SetDstRouter(nodeID string)

	GetSrcRouter() string
	SetSrcRouter(nodeID string)
}

type connectionImpl struct {
	flitChan   chan packet.Flit
	creditChan map[int]chan int
	destRouter string
	srcRouter  string
	logger     zerolog.Logger
}

func NewConnection(maxPriority int, logger zerolog.Logger) (*connectionImpl, error) {
	creditChan := make(map[int]chan int, maxPriority)
	for i := 1; i <= maxPriority; i++ {
		creditChan[i] = make(chan int, 1)
	}

	logger.Trace().Int("credit_channels", maxPriority).Msg("new connection")
	return &connectionImpl{
		flitChan:   make(chan packet.Flit, 1),
		creditChan: creditChan,
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

func (c *connectionImpl) GetDstRouter() string {
	return c.destRouter
}

func (c *connectionImpl) SetDstRouter(nodeID string) {
	c.destRouter = nodeID
}

func (c *connectionImpl) GetSrcRouter() string {
	return c.srcRouter
}

func (c *connectionImpl) SetSrcRouter(nodeID string) {
	c.srcRouter = nodeID
}
