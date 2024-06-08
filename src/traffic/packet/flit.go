package packet

import (
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
	UUID() uuid.UUID
	Type() FlitType
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
}

type HeaderFlit interface {
	UUID() uuid.UUID
	Type() FlitType
	TrafficFlowID() string
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
	Deadline() int
	Route() domain.Route
}

type headerFlit struct {
	uuid          uuid.UUID
	trafficFlowID string
	packetUUID    uuid.UUID
	pktIndex      int
	priority      int
	deadline      int
	route         domain.Route
}

type BodyFlit interface {
	UUID() uuid.UUID
	Type() FlitType
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
	DataSize() int
}

type bodyFlit struct {
	uuid       uuid.UUID
	packetUUID uuid.UUID
	pktIndex   int
	priority   int
	dataSize   int
}

type TailFlit interface {
	UUID() uuid.UUID
	Type() FlitType
	PacketUUID() uuid.UUID
	PacketIndex() int
	Priority() int
}

type tailFlit struct {
	uuid       uuid.UUID
	packetUUID uuid.UUID
	pktIndex   int
	priority   int
}

func NewHeaderFlit(trafficFlowID string, packetUUID uuid.UUID, pktIndex int, priority, deadline int, route domain.Route) *headerFlit {
	flitID := uuid.New()

	log.Log.Trace().
		Str("traffic_flow", trafficFlowID).Str("packet", packetUUID.String()).Str("flit", flitID.String()).
		Msg("new header flit")

	return &headerFlit{
		trafficFlowID: trafficFlowID,
		uuid:          flitID,
		packetUUID:    packetUUID,
		pktIndex:      pktIndex,
		priority:      priority,
		deadline:      deadline,
		route:         route,
	}
}

func NewBodyFlit(packetUUID uuid.UUID, pktIndex int, priority, dataSize int) *bodyFlit {
	flitID := uuid.New()

	log.Log.Trace().
		Str("packet", packetUUID.String()).Str("flit", flitID.String()).
		Msg("new body flit")

	return &bodyFlit{
		uuid:       flitID,
		packetUUID: packetUUID,
		pktIndex:   pktIndex,
		priority:   priority,
		dataSize:   dataSize,
	}
}

func NewTailFlit(packetUUID uuid.UUID, pktIndex int, priority int) *tailFlit {
	flitID := uuid.New()

	log.Log.Trace().
		Str("packet", packetUUID.String()).Str("flit", flitID.String()).
		Msg("new tail flit")

	return &tailFlit{
		uuid:       flitID,
		packetUUID: packetUUID,
		pktIndex:   pktIndex,
		priority:   priority,
	}
}

func (f *headerFlit) UUID() uuid.UUID {
	return f.uuid
}

func (f *headerFlit) Type() FlitType {
	return HeaderFlitType
}

func (f *headerFlit) PacketIndex() int {
	return f.pktIndex
}

func (f *headerFlit) TrafficFlowID() string {
	return f.trafficFlowID
}

func (f *headerFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
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

func (f *bodyFlit) UUID() uuid.UUID {
	return f.uuid
}

func (f *bodyFlit) Type() FlitType {
	return BodyFlitType
}

func (f *bodyFlit) PacketIndex() int {
	return f.pktIndex
}

func (f *bodyFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
}

func (f *bodyFlit) Priority() int {
	return f.priority
}

func (f *bodyFlit) DataSize() int {
	return f.dataSize
}

func (f *tailFlit) UUID() uuid.UUID {
	return f.uuid
}

func (f *tailFlit) Type() FlitType {
	return TailFlitType
}

func (f *tailFlit) PacketIndex() int {
	return f.pktIndex
}

func (f *tailFlit) PacketUUID() uuid.UUID {
	return f.packetUUID
}

func (f *tailFlit) Priority() int {
	return f.priority
}
