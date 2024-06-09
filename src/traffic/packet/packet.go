package packet

import (
	"errors"
	"fmt"
	"math"

	"main/log"
	"main/src/domain"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

// type PktEvent string

// const (
// 	PktCreated           PktEvent = "Packet created"
// 	PktReleasedToNetwork PktEvent = "Packet released to the network"
// 	PktSplitToFlits      PktEvent = "Packet split into flits"
// )

type Packet interface {
	UUID() uuid.UUID
	TrafficFlowID() string
	Priority() int
	Deadline() int
	Route() domain.Route
	BodySize() int
	Flits(flitSize int) []Flit
}

type packet struct {
	trafficFlowID string
	uuid          uuid.UUID
	priority      int
	deadline      int
	route         domain.Route
	bodySize      int
}

func newPacketWithUUID(trafficFlowID string, uuid uuid.UUID, priority, deadline int, route domain.Route, bodySize int) *packet {
	log.Log.Trace().Str("traffic_flow", trafficFlowID).Str("id", uuid.String()).Msg("new packet")

	return &packet{
		trafficFlowID: trafficFlowID,
		uuid:          uuid,
		priority:      priority,
		deadline:      deadline,
		route:         route,
		bodySize:      bodySize,
	}
}

func NewPacket(trafficFlowID string, priority, deadline int, router domain.Route, bodySize int) *packet {
	return newPacketWithUUID(trafficFlowID, uuid.New(), priority, deadline, router, bodySize)
}

func (p *packet) UUID() uuid.UUID {
	return p.uuid
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

func (p *packet) Flits(flitSize int) []Flit {
	flits := make([]Flit, 1+p.bodyFlitCount(flitSize)+1)

	flits[0] = NewHeaderFlit(p.TrafficFlowID(), p.UUID(), 0, p.priority, p.deadline, p.route)

	bodyFlits := p.bodyFlits(flitSize)
	for i := 0; i < len(bodyFlits); i++ {
		flits[i+1] = bodyFlits[i]
	}

	flits[len(flits)-1] = NewTailFlit(p.TrafficFlowID(), p.UUID(), len(flits)-1, p.priority)

	return flits
}

func (p *packet) bodyFlits(flitSize int) []BodyFlit {
	bodyFlits := make([]BodyFlit, p.bodyFlitCount(flitSize))
	for i := 0; i < p.bodyFlitCount(flitSize); i++ {
		if (i+1)*flitSize < p.bodySize {
			bodyFlits[i] = NewBodyFlit(p.TrafficFlowID(), p.UUID(), i+1, p.priority, flitSize)
		} else {
			bodyFlits[i] = NewBodyFlit(p.TrafficFlowID(), p.UUID(), i+1, p.priority, p.bodySize-(i*flitSize))
		}
	}

	return bodyFlits
}

func (p *packet) bodyFlitCount(flitSize int) int {
	return int(math.Ceil(float64(p.bodySize) / float64(flitSize)))
}

func EqualPackets(pkt1, pkt2 Packet) error {
	if pkt1 == nil || pkt2 == nil {
		return domain.ErrNilParameter
	}

	if pkt1.UUID() != pkt2.UUID() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("UUID: %s != %s", pkt1.UUID().String(), pkt2.UUID().String()))
	}

	if pkt1.TrafficFlowID() != pkt2.TrafficFlowID() {
		return errors.Join(domain.ErrPacketsNotEqual, fmt.Errorf("TrafficFlowID: %s != %s", pkt1.TrafficFlowID(), pkt2.TrafficFlowID()))
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
