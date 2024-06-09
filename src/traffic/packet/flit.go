package packet

import (
	"fmt"
	"main/log"
	"main/src/domain"

	"github.com/google/uuid"
)

type FlitType string

const (
	HeaderFlitType FlitType = "header"
	BodyFlitType   FlitType = "body"
	TailFlitType   FlitType = "tail"
)

func (f FlitType) String() string {
	return string(f)
}

type Flit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
}

type HeaderFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
	Deadline() int
	Route() domain.Route
}

type headerFlit struct {
	id            string
	trafficFlowID string
	packetUUID    uuid.UUID
	packetIndex   int
	priority      int
	deadline      int
	route         domain.Route
}

type BodyFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
	DataSize() int
}

type bodyFlit struct {
	id            string
	trafficFlowID string
	packetUUID    uuid.UUID
	packetIndex   int
	priority      int
	dataSize      int
}

type TailFlit interface {
	ID() string
	Type() FlitType
	TrafficFlowID() string
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
}

type tailFlit struct {
	id            string
	trafficFlowID string
	packetUUID    uuid.UUID
	packetIndex   int
	priority      int
}

func id(trafficFlowID string, packetUUID uuid.UUID, packetIndex int) string {
	return fmt.Sprintf("%s_%s_%d", trafficFlowID, packetUUID.String(), packetIndex)
}

func NewHeaderFlit(trafficFlowID string, packetUUID uuid.UUID, packetIndex int, priority, deadline int, route domain.Route) *headerFlit {
	id := id(trafficFlowID, packetUUID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &headerFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetUUID:    packetUUID,
		packetIndex:   packetIndex,
		priority:      priority,
		deadline:      deadline,
		route:         route,
	}
}

func NewBodyFlit(trafficFlowID string, packetUUID uuid.UUID, packetIndex int, priority, dataSize int) *bodyFlit {
	id := id(trafficFlowID, packetUUID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &bodyFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetUUID:    packetUUID,
		packetIndex:   packetIndex,
		priority:      priority,
		dataSize:      dataSize,
	}
}

func NewTailFlit(trafficFlowID string, packetUUID uuid.UUID, packetIndex int, priority int) *tailFlit {
	id := id(trafficFlowID, packetUUID, packetIndex)

	log.Log.Trace().Str("flit", id).Msg("new header flit")

	return &tailFlit{
		id:            id,
		trafficFlowID: trafficFlowID,
		packetUUID:    packetUUID,
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

func (f *headerFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
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

func (f *bodyFlit) ID() string {
	return f.id
}

func (f *bodyFlit) Type() FlitType {
	return BodyFlitType
}

func (f *bodyFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *bodyFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
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

func (f *tailFlit) ID() string {
	return f.id
}

func (f *tailFlit) Type() FlitType {
	return TailFlitType
}

func (f *tailFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *tailFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
}

func (f *tailFlit) PacketIndex() int {
	return f.packetIndex
}

func (f *tailFlit) Priority() int {
	return f.priority
}
