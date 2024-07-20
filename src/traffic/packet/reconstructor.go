package packet

import (
	"io"

	"main/src/domain"

	"github.com/rs/zerolog"
)

type Reconstructor interface {
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

func NewReconstructor(headerFlit HeaderFlit, logger zerolog.Logger) (*reconstructor, error) {
	if headerFlit == nil {
		return nil, domain.ErrNilParameter
	}

	r := &reconstructor{
		bodyFlits: make([]BodyFlit, 0),
		logger:    logger.With().Str("packet", headerFlit.PacketID()).Logger(),
	}
	r.logger.Trace().Msg("new packet reconstructor")

	r.headerFlit = headerFlit
	r.logger.Trace().Str("flit", headerFlit.ID()).Str("type", headerFlit.Type().String()).Msg("set header flit")

	return r, nil
}

func (r *reconstructor) AddBody(bodyFlit BodyFlit) error {
	if bodyFlit == nil {
		return domain.ErrNilParameter
	}

	r.bodyFlits = append(r.bodyFlits, bodyFlit)
	r.logger.Trace().Str("flit", bodyFlit.ID()).Str("type", bodyFlit.Type().String()).Msg("added body flit")

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
	r.logger.Trace().Str("flit", tailFlit.ID()).Str("type", tailFlit.Type().String()).Msg("set tail flit")

	return nil
}

func (r *reconstructor) Reconstruct() (Packet, error) {
	if r.headerFlit == nil || r.tailFlit == nil {
		return nil, domain.ErrFlitUnset
	}

	pkt := NewPacket(
		r.headerFlit.TrafficFlowID(),
		r.headerFlit.PacketIndex(),
		r.headerFlit.Priority(),
		r.headerFlit.Deadline(),
		r.headerFlit.Route(),
		// The packet size is the sum of the header, body, and tail flits.
		len(r.bodyFlits)+2,
		// Reconstructed packets should only be used for records purposes, so the logger can be discarded.
		zerolog.New(io.Discard),
	)

	r.logger.Trace().Str("packet", pkt.ID()).Msg("reconstructed packet")
	return pkt, nil
}
