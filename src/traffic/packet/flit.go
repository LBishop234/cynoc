package packet

import (
	"fmt"
	"main/log"
	"main/src/domain"
)

type FlitType string

const (
	HeaderFlitType FlitType = "header"
	BodyFlitType   FlitType = "body"
	TailFlitType   FlitType = "tail"
)

func (t FlitType) String() string {
	return string(t)
}

type FlitEvent string

const (
	FlitCreated     FlitEvent = "flit created"
	FlitTransmitted FlitEvent = "flit transmitted"
	FlitArrived     FlitEvent = "flit arrived"
)

func (e FlitEvent) String() string {
	return string(e)
}

type Flit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() int
	Priority() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type HeaderFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() int
	Priority() int
	Deadline() int
	Route() domain.Route
	RecordEvent(cycle int, event FlitEvent, location string)
}

type headerFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   int
	priority      int
	deadline      int
	route         domain.Route
}

type BodyFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() int
	Priority() int
	DataSize() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type bodyFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   int
	priority      int
	dataSize      int
}

type TailFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() int
	Priority() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type tailFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   int
	priority      int
}

func id(trafficFlowID string, packetID string, packetIndex int) string {
	return fmt.Sprintf("%s_%s_%d", trafficFlowID, packetID, packetIndex)
}

func NewHeaderFlit(trafficFlowID string, packetID string, packetIndex int, priority, deadline int, route domain.Route) *headerFlit {
	id := id(trafficFlowID, packetID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &headerFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetID:      packetID,
		packetIndex:   packetIndex,
		priority:      priority,
		deadline:      deadline,
		route:         route,
	}
}

func NewBodyFlit(trafficFlowID string, packetID string, packetIndex int, priority, dataSize int) *bodyFlit {
	id := id(trafficFlowID, packetID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &bodyFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetID:      packetID,
		packetIndex:   packetIndex,
		priority:      priority,
		dataSize:      dataSize,
	}
}

func NewTailFlit(trafficFlowID string, packetID string, packetIndex int, priority int) *tailFlit {
	id := id(trafficFlowID, packetID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &tailFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetID:      packetID,
		packetIndex:   packetIndex,
		priority:      priority,
	}
}

func (f *headerFlit) ID() string {
	return f.id
}

func (f *headerFlit) Type() FlitType {
	return HeaderFlitType
}

func (f *headerFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *headerFlit) PacketID() string {
	return f.packetID
}

func (f *headerFlit) PacketIndex() int {
	return f.packetIndex
}

func (f *headerFlit) Priority() int {
	return f.priority
}

func (f *headerFlit) Deadline() int {
	return f.deadline
}

func (f *headerFlit) Route() domain.Route {
	return f.route
}

func (f *headerFlit) RecordEvent(cycle int, event FlitEvent, location string) {
	recordEvent(f, cycle, event, location)
}

func (f *bodyFlit) ID() string {
	return f.id
}

func (f *bodyFlit) Type() FlitType {
	return BodyFlitType
}

func (f *bodyFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *bodyFlit) PacketID() string {
	return f.packetID
}

func (f *bodyFlit) PacketIndex() int {
	return f.packetIndex
}

func (f *bodyFlit) Priority() int {
	return f.priority
}

func (f *bodyFlit) DataSize() int {
	return f.dataSize
}

func (f *bodyFlit) RecordEvent(cycle int, event FlitEvent, location string) {
	recordEvent(f, cycle, event, location)
}

func (f *tailFlit) ID() string {
	return f.id
}

func (f *tailFlit) Type() FlitType {
	return TailFlitType
}

func (f *tailFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *tailFlit) PacketID() string {
	return f.packetID
}

func (f *tailFlit) PacketIndex() int {
	return f.packetIndex
}

func (f *tailFlit) Priority() int {
	return f.priority
}

func (f *tailFlit) RecordEvent(cycle int, event FlitEvent, location string) {
	recordEvent(f, cycle, event, location)
}

func recordEvent(f Flit, cycle int, event FlitEvent, location string) {
	log.Log.Trace().Int("cycle", cycle).Str("flit", f.ID()).Str("event", event.String()).Str("location", location).Msgf("%s at %s", event.String(), location)
}
