package packet

import (
	"errors"
	"fmt"

	"main/src/domain"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
)

// type PktEvent string

// const (
// 	PktCreated           PktEvent = "Packet created"
// 	PktReleasedToNetwork PktEvent = "Packet released to the network"
// 	PktSplitToFlits      PktEvent = "Packet split into flits"
// )

type Packet interface {
	ID() string
	TrafficFlowID() string
	PacketIndex() string
	Priority() int
	Deadline() int
	Route() domain.Route
	BodySize() int
	Flits() []Flit
}

type packet struct {
	id            string
	trafficFlowID string
	packetIndex   string
	priority      int
	deadline      int
	route         domain.Route
	bodySize      int

	logger zerolog.Logger
}

func newPacketID(trafficFlowID, packetIndex string) string {
	return fmt.Sprintf("%s-%s", trafficFlowID, packetIndex)
}

func NewPacket(trafficFlowID string, packetIndex string, priority, deadline int, route domain.Route, packetSize int, logger zerolog.Logger) *packet {
	id := newPacketID(trafficFlowID, packetIndex)
	bodySize := packetSize - 2

	logger.Trace().Str("packet", id).Msg("new packet")

	return &packet{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetIndex:   packetIndex,
		priority:      priority,
		deadline:      deadline,
		route:         route,
		bodySize:      bodySize,
		logger:        logger,
	}
}

func (p *packet) ID() string {
	return p.id
}

func (p *packet) PacketIndex() string {
	return p.packetIndex
}

func (p *packet) TrafficFlowID() string {
	return p.trafficFlowID
}

func (p *packet) Priority() int {
	return p.priority
}

func (p *packet) Deadline() int {
	return p.deadline
}

func (p *packet) Route() domain.Route {
	return p.route
}

func (p *packet) BodySize() int {
	return p.bodySize
}

func (p *packet) Flits() []Flit {
	flits := make([]Flit, 1+p.bodySize+1)

	flits[0] = NewHeaderFlit(p.TrafficFlowID(), p.PacketIndex(), 0, p.priority, p.deadline, p.route, p.logger)

	bodyFlits := p.bodyFlits()
	for i := 0; i < len(bodyFlits); i++ {
		flits[i+1] = bodyFlits[i]
	}

	flits[len(flits)-1] = NewTailFlit(p.TrafficFlowID(), p.PacketIndex(), len(flits)-1, p.priority, p.logger)

	return flits
}

func (p *packet) bodyFlits() []BodyFlit {
	bodyFlits := make([]BodyFlit, p.bodySize)
	for i := 0; i < p.bodySize; i++ {
		bodyFlits[i] = NewBodyFlit(p.TrafficFlowID(), p.PacketIndex(), i+1, p.priority, p.logger)
	}

	return bodyFlits
}

func EqualPackets(pkt1, pkt2 Packet) error {
	if pkt1 == nil || pkt2 == nil {
		return domain.ErrNilParameter
	}

	if pkt1.ID() != pkt2.ID() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("ID: %s != %s", pkt1.ID(), pkt2.ID()))
	}

	if pkt1.TrafficFlowID() != pkt2.TrafficFlowID() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("TrafficFlowID: %s != %s", pkt1.TrafficFlowID(), pkt2.TrafficFlowID()))
	}

	if pkt1.PacketIndex() != pkt2.PacketIndex() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("ID: %s != %s", pkt1.PacketIndex(), pkt2.PacketIndex()))
	}

	if pkt1.Priority() != pkt2.Priority() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("Priority: %d != %d", pkt1.Priority(), pkt2.Priority()))
	}

	if pkt1.Deadline() != pkt2.Deadline() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("Deadline: %d != %d", pkt1.Deadline(), pkt2.Deadline()))
	}

	if len(pkt1.Route()) != len(pkt2.Route()) {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("Route: %v != %v", spew.Sprint(pkt1.Route()), spew.Sprint(pkt2.Route())))
	}

	for i := 0; i < len(pkt1.Route()); i++ {
		if pkt1.Route()[i] != pkt2.Route()[i] {
			return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("Route: %v != %v", spew.Sprint(pkt1.Route()), spew.Sprint(pkt2.Route())))
		}
	}

	if pkt1.BodySize() != pkt2.BodySize() {
		return errors.Join(domain.ErrPacketsNotEqual, errors.New("BodySize"))
	}

	return nil
}
