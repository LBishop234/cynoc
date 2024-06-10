package packet

import (
	"io"
	"main/src/domain"

	"github.com/rs/zerolog"
)

type Reconstructor interface {
	SetHeader(headerFlit HeaderFlit) error
	AddBody(bodyFlit BodyFlit) error
	SetTail(tailFlit TailFlit) error

	Reconstruct() (Packet, error)
}

type reconstructor struct {
	headerFlit HeaderFlit
	bodyFlits  []BodyFlit
	tailFlit   TailFlit

	logger zerolog.Logger
}

func NewReconstructor(logger zerolog.Logger) *reconstructor {
	logger.Trace().Msg("new packet reconstructor")

	return &reconstructor{
		bodyFlits: make([]BodyFlit, 0),
		logger:    logger,
	}
}

func (r *reconstructor) SetHeader(headerFlit HeaderFlit) error {
	if headerFlit == nil {
		return domain.ErrNilParameter
	}

	if r.headerFlit != nil {
		return domain.ErrFlitAlreadySet
	}

	r.headerFlit = headerFlit
	return nil
}

func (r *reconstructor) AddBody(bodyFlit BodyFlit) error {
	if bodyFlit == nil {
		return domain.ErrNilParameter
	}

	r.bodyFlits = append(r.bodyFlits, bodyFlit)
	return nil
}

func (r *reconstructor) SetTail(tailFlit TailFlit) error {
	if tailFlit == nil {
		return domain.ErrNilParameter
	}

	if r.tailFlit != nil {
		return domain.ErrFlitAlreadySet
	}

	r.tailFlit = tailFlit
	return nil
}

func (r *reconstructor) Reconstruct() (Packet, error) {
	if r.headerFlit == nil || r.tailFlit == nil {
		return nil, domain.ErrFlitUnset
	}

	bodySize := 0
	for i := 0; i < len(r.bodyFlits); i++ {
		bodySize += r.bodyFlits[i].DataSize()
	}

	return NewPacket(
		r.headerFlit.TrafficFlowID(),
		r.headerFlit.PacketID(),
		r.headerFlit.Priority(),
		r.headerFlit.Deadline(),
		r.headerFlit.Route(),
		bodySize,
		zerolog.New(io.Discard),
	), nil
}
