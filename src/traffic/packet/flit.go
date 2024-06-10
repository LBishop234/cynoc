package packet

import (
	"fmt"
	"main/src/domain"

	"github.com/rs/zerolog"
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
	PacketIndex() string
	FlitIndex() int
	Priority() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type HeaderFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() string
	FlitIndex() int
	Priority() int
	Deadline() int
	Route() domain.Route
	RecordEvent(cycle int, event FlitEvent, location string)
}

type headerFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   string
	flitIndex     int
	priority      int
	deadline      int
	route         domain.Route
	logger        zerolog.Logger
}

type BodyFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() string
	FlitIndex() int
	Priority() int
	DataSize() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type bodyFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   string
	flitIndex     int
	priority      int
	dataSize      int
	logger        zerolog.Logger
}

type TailFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketID() string
	PacketIndex() string
	FlitIndex() int
	Priority() int
	RecordEvent(cycle int, event FlitEvent, location string)
}

type tailFlit struct {
	id            string
	trafficFlowID string
	packetID      string
	packetIndex   string
	flitIndex     int
	priority      int
	logger        zerolog.Logger
}

func newFlitID(trafficFlowID, packetIndex string, flitIndex int) string {
	return fmt.Sprintf("%s-%d", newPacketID(trafficFlowID, packetIndex), flitIndex)
}

func NewHeaderFlit(trafficFlowID string, packetIndex string, flitIndex int, priority, deadline int, route domain.Route, logger zerolog.Logger) *headerFlit {
	id := newFlitID(trafficFlowID, packetIndex, flitIndex)
	packetID := newPacketID(trafficFlowID, packetIndex)

	logger = logger.With().Str("flit", id).Logger()
	logger.Trace().Msg("new header flit")

	return &headerFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetID:      packetID,
		packetIndex:   packetIndex,
		flitIndex:     flitIndex,
		priority:      priority,
		deadline:      deadline,
		route:         route,
	}
}

func NewBodyFlit(trafficFlowID string, packetIndex string, flitIndex int, priority, dataSize int, logger zerolog.Logger) *bodyFlit {
	id := newFlitID(trafficFlowID, packetIndex, flitIndex)
	packetID := newPacketID(trafficFlowID, packetIndex)

	logger = logger.With().Str("flit", id).Logger()
	logger.Trace().Msg("new header flit")

	return &bodyFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetID:      packetID,
		packetIndex:   packetIndex,
		flitIndex:     flitIndex,
		priority:      priority,
		dataSize:      dataSize,
	}
}

func NewTailFlit(trafficFlowID string, packetIndex string, flitIndex int, priority int, logger zerolog.Logger) *tailFlit {
	id := newFlitID(trafficFlowID, packetIndex, flitIndex)
	packetID := newPacketID(trafficFlowID, packetIndex)

	logger = logger.With().Str("flit", id).Logger()
	logger.Trace().Msg("new header flit")

	return &tailFlit{
		id:            id,
		packetID:      packetID,
		trafficFlowID: trafficFlowID,
		packetIndex:   packetIndex,
		flitIndex:     flitIndex,
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

func (f *headerFlit) PacketIndex() string {
	return f.packetIndex
}

func (f *headerFlit) FlitIndex() int {
	return f.flitIndex
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
	recordEvent(&f.logger, f, cycle, event, location)
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

func (f *bodyFlit) PacketIndex() string {
	return f.packetIndex
}

func (f *bodyFlit) FlitIndex() int {
	return f.flitIndex
}

func (f *bodyFlit) Priority() int {
	return f.priority
}

func (f *bodyFlit) DataSize() int {
	return f.dataSize
}

func (f *bodyFlit) RecordEvent(cycle int, event FlitEvent, location string) {
	recordEvent(&f.logger, f, cycle, event, location)
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

func (f *tailFlit) PacketIndex() string {
	return f.packetIndex
}

func (f *tailFlit) FlitIndex() int {
	return f.flitIndex
}

func (f *tailFlit) Priority() int {
	return f.priority
}

func (f *tailFlit) RecordEvent(cycle int, event FlitEvent, location string) {
	recordEvent(&f.logger, f, cycle, event, location)
}

func recordEvent(logger *zerolog.Logger, f Flit, cycle int, event FlitEvent, location string) {
	logger.Trace().Int("cycle", cycle).Str("flit", f.ID()).Str("event", event.String()).Str("location", location).Msgf("%s at %s", event.String(), location)
}
