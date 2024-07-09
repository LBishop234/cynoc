package components

import (
	"main/src/domain"
	"main/src/traffic/packet"

	"github.com/rs/zerolog"
)

type inputPort interface {
	connection() Connection
	readIntoBuffer(cycle int) error
	peakBuffer(priority int) (packet.Flit, bool)
	readOutOfBuffer(cycle, priority int) (packet.Flit, bool)
}

type outputPort interface {
	connection() Connection
	allowedToSend(priority int) bool
	sendFlit(cycle int, flit packet.Flit) error
	updateCredits()
}

type inputPortImpl struct {
	conn   Connection
	buff   buffer
	logger zerolog.Logger
}

type outputPortImpl struct {
	conn    Connection
	credits map[int]int
	logger  zerolog.Logger
}

func newInputPort(conn Connection, buff buffer, logger zerolog.Logger) (*inputPortImpl, error) {
	localLogger := logger.With().Str("port", "input_port").Logger()

	if conn == nil || buff == nil {
		localLogger.Error().Err(domain.ErrNilParameter).Msg("invalid input port parameters")
		return nil, domain.ErrNilParameter
	}

	for priority, credChan := range conn.creditChannels() {
		localLogger.Trace().
			Int("priority", priority).Int("capacity", cap(credChan)).Int("Credit", buff.vChanCapacity()).
			Msg("publishing input port virtual channel credits to connection source object")

		credChan <- buff.vChanCapacity()
	}
	localLogger.Trace().Msg("published input port buffer capacity to connection source")

	localLogger.Trace().Msg("new input port")
	return &inputPortImpl{
		conn:   conn,
		buff:   buff,
		logger: localLogger,
	}, nil
}

func newOutputPort(conn Connection, maxPriority int, logger zerolog.Logger) (*outputPortImpl, error) {
	localLogger := logger.With().Str("port", "output_port").Logger()

	if conn == nil {
		localLogger.Error().Err(domain.ErrNilParameter).Msg("invalid output port parameters")
		return nil, domain.ErrNilParameter
	}

	localLogger.Trace().Msg("new output port")
	return &outputPortImpl{
		conn:    conn,
		credits: make(map[int]int, maxPriority),
		logger:  localLogger,
	}, nil
}

func (i *inputPortImpl) connection() Connection {
	return i.conn
}

func (i *inputPortImpl) readIntoBuffer(cycle int) (err error) {
	for len(i.conn.flitChannel()) > 0 {
		flit := <-i.conn.flitChannel()

		if err = i.buff.addFlit(flit); err != nil {
			return err
		}

		i.logger.Debug().
			Int("cycle", cycle).Str("flit", flit.ID()).Str("type", flit.Type().String()).
			Msg("flit arrived at component")
	}
	return nil
}

func (i *inputPortImpl) peakBuffer(priority int) (packet.Flit, bool) {
	return i.buff.peakFlit(priority)
}

func (i *inputPortImpl) readOutOfBuffer(cycle, priority int) (packet.Flit, bool) {
	flit, exists := i.buff.popFlit(priority)
	if exists {
		i.logger.Trace().
			Int("cycle", cycle).Str("flit", flit.ID()).Str("type", flit.Type().String()).
			Msg("flit read out of buffer")

		i.conn.creditChannel(flit.Priority()) <- 1
	}
	return flit, exists
}

func (o *outputPortImpl) connection() Connection {
	return o.conn
}

func (o *outputPortImpl) allowedToSend(priority int) bool {
	return o.credits[priority] > 0 && len(o.conn.flitChannel()) < o.conn.flitBandwidth()
}

func (o *outputPortImpl) sendFlit(cycle int, flit packet.Flit) error {
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
