package packet

import (
	"testing"

	"main/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFlitTypeString(t *testing.T) {
	t.Parallel()

	t.Run("Header", func(t *testing.T) {
		assert.Equal(t, "header", HeaderFlitType.String())
	})

	t.Run("Body", func(t *testing.T) {
		assert.Equal(t, "body", BodyFlitType.String())
	})

	t.Run("Tail", func(t *testing.T) {
		assert.Equal(t, "tail", TailFlitType.String())
	})
}

func TestNewHeaderFlit(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	packetUUID := uuid.New()
	priority := 1
	deadline := 100
	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	flit := NewHeaderFlit(trafficFlowID, packetUUID, priority, deadline, route)
	assert.Equal(t, trafficFlowID, flit.trafficFlowID)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, route, flit.route)

	assert.Implements(t, (*HeaderFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewBodyFlit(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	priority := 1
	dataSize := 4

	flit := NewBodyFlit(packetUUID, priority, dataSize)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, dataSize, flit.dataSize)

	assert.Implements(t, (*BodyFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewTailFlit(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	priority := 1

	flit := NewTailFlit(packetUUID, priority)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, priority, flit.priority)

	assert.Implements(t, (*TailFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestFlitUUID(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	flit := NewHeaderFlit("t", uuid.New(), 1, 100, route)
	assert.Equal(t, flit.uuid, flit.UUID())
}

func TestHeaderFlitType(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	assert.Equal(t, HeaderFlitType, NewHeaderFlit("t", uuid.New(), 1, 100, route).Type())
}

func TestHeaderFlitTrafficFlowID(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	assert.Equal(t, trafficFlowID, NewHeaderFlit(trafficFlowID, uuid.New(), 1, 100, route).TrafficFlowID())
}

func TestHeaderFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	assert.Equal(t, packetUUID, NewHeaderFlit("t", packetUUID, 1, 100, route).PacketUUID())
}

func TestHeaderFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	assert.Equal(t, priority, NewHeaderFlit("t", uuid.New(), priority, 100, route).Priority())
}

func TestHeaderFlitDeadline(t *testing.T) {
	t.Parallel()

	deadline := 100
	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	assert.Equal(t, deadline, NewHeaderFlit("t", uuid.New(), 1, deadline, route).Deadline())
}

func TestHeaderFlitRoute(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	flit := NewHeaderFlit("t", uuid.New(), 1, 100, route)
	assert.Equal(t, route, flit.Route())
}

func TestBodyFlitUUID(t *testing.T) {
	t.Parallel()

	flit := NewBodyFlit(uuid.New(), 1, 1)
	assert.Equal(t, flit.uuid, flit.UUID())
}

func TestBodyFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, BodyFlitType, NewBodyFlit(uuid.New(), 1, 1).Type())
}

func TestBodyFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	assert.Equal(t, packetUUID, NewBodyFlit(packetUUID, 1, 1).PacketUUID())
}

func TestBodyFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewBodyFlit(uuid.New(), 1, 1).Priority())
}

func TestBodyFlitDataSize(t *testing.T) {
	t.Parallel()

	dataSize := 4

	flit := NewBodyFlit(uuid.New(), 1, dataSize)
	assert.Equal(t, dataSize, flit.DataSize())
}

func TestTailFlitUUID(t *testing.T) {
	t.Parallel()

	flit := NewTailFlit(uuid.New(), 1)
	assert.Equal(t, flit.uuid, flit.UUID())
}

func TestTailFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, TailFlitType, NewTailFlit(uuid.New(), 1).Type())
}

func TestTailFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	assert.Equal(t, packetUUID, NewTailFlit(packetUUID, 1).PacketUUID())
}

func TestTailFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewTailFlit(uuid.New(), 1).Priority())
}
