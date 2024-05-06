package components

import (
	"main/core/log"
	"main/domain"
	"main/traffic/packet"
)

type inputPort interface {
	connection() Connection
	readIntoBuffer() error
	peakBuffer(priority int) (packet.Flit, bool)
	readOutOfBuffer(priority int) (packet.Flit, bool)
}

type outputPort interface {
	connection() Connection
	allowedToSend(priority int) bool
	sendFlit(flit packet.Flit) error
	updateCredits()
}

type inputPortImpl struct {
	conn Connection
	buff buffer
}

type outputPortImpl struct {
	conn    Connection
	credits map[int]int
}

func newInputPort(conn Connection, buff buffer) (*inputPortImpl, error) {
	if conn == nil || buff == nil {
		log.Log.Error().Msg("nil parameter passed to function")
		return nil, domain.ErrNilParameter
	}

	for priority, credChan := range conn.creditChannels() {
		log.Log.Trace().Int("priority", priority).Msg("publishing input port buffer capacity for priority to connection src")
		credChan <- buff.vChanCapacity()
	}

	log.Log.Trace().Msg("published input port buffer capacity to connection src")

	log.Log.Trace().Msg("new input port")
	return &inputPortImpl{
		conn: conn,
		buff: buff,
	}, nil
}

func newOutputPort(conn Connection, maxPriority int) (*outputPortImpl, error) {
	if conn == nil {
		log.Log.Error().Msg("nil parameter passed to function")
		return nil, domain.ErrNilParameter
	}

	log.Log.Trace().Msg("new output port")
	return &outputPortImpl{
		conn:    conn,
		credits: make(map[int]int, maxPriority),
	}, nil
}

func (i *inputPortImpl) connection() Connection {
	return i.conn
}

func (i *inputPortImpl) readIntoBuffer() (err error) {
	for len(i.conn.flitChannel()) > 0 {
		if err = i.buff.addFlit(<-i.conn.flitChannel()); err != nil {
			return err
		}
	}
	return nil
}

func (i *inputPortImpl) peakBuffer(priority int) (packet.Flit, bool) {
	return i.buff.peakFlit(priority)
}

func (i *inputPortImpl) readOutOfBuffer(priority int) (packet.Flit, bool) {
	flit, exists := i.buff.popFlit(priority)
	if exists {
		i.conn.creditChannel(flit.Priority()) <- 1
	}
	return flit, exists
}

func (o *outputPortImpl) connection() Connection {
	return o.conn
}

func (o *outputPortImpl) allowedToSend(priority int) bool {
	return o.credits[priority] > 0 && len(o.conn.flitChannel()) == 0
}

func (o *outputPortImpl) sendFlit(flit packet.Flit) error {
	if o.allowedToSend(flit.Priority()) {
		o.credits[flit.Priority()]--
		o.conn.flitChannel() <- flit
		return nil
	} else {
		return domain.ErrPortNoCredit
	}
}

func (o *outputPortImpl) updateCredits() {
	for priority := range o.conn.creditChannels() {
		for len(o.conn.creditChannel(priority)) > 0 {
			o.credits[priority] += <-o.conn.creditChannel(priority)
		}
	}
}
