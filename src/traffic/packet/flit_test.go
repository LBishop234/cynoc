package packet

import (
	"math"
	"math/rand"
	"testing"

	"main/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testDummyRoute(t *testing.T) (src, dst domain.NodeID, route domain.Route) {
	src = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route = domain.Route{src, dst}

	return src, dst, route
}

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
	packetIndex := 0
	priority := 1
	deadline := 100
	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit(trafficFlowID, packetUUID, packetIndex, priority, deadline, route)
	assert.Equal(t, trafficFlowID, flit.trafficFlowID)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, packetIndex, flit.pktIndex)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, route, flit.route)

	assert.Implements(t, (*HeaderFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewBodyFlit(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	packetIndex := 1
	priority := 1
	dataSize := 4

	flit := NewBodyFlit("t", packetUUID, packetIndex, priority, dataSize)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, packetIndex, flit.pktIndex)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, dataSize, flit.dataSize)

	assert.Implements(t, (*BodyFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewTailFlit(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	packetIndex := 2
	priority := 1

	flit := NewTailFlit("t", packetUUID, packetIndex, priority)
	assert.Equal(t, packetUUID, flit.packetUUID)
	assert.Equal(t, packetIndex, flit.pktIndex)
	assert.Equal(t, priority, flit.priority)

	assert.Implements(t, (*TailFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestFlitUUID(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit("t", uuid.New(), 0, 1, 100, route)
	assert.Equal(t, flit.id, flit.ID())
}

func TestHeaderFlitType(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	assert.Equal(t, HeaderFlitType, NewHeaderFlit("t", uuid.New(), 0, 1, 100, route).Type())
}

func TestHeaderFlitTrafficFlowID(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	_, _, route := testDummyRoute(t)

	assert.Equal(t, trafficFlowID, NewHeaderFlit(trafficFlowID, uuid.New(), 0, 1, 100, route).TrafficFlowID())
}

func TestHeaderFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	_, _, route := testDummyRoute(t)

	assert.Equal(t, packetUUID, NewHeaderFlit("t", packetUUID, 0, 1, 100, route).PacketUUID())
}

func TestHeaderFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	_, _, route := testDummyRoute(t)

	assert.Equal(t, packetIndex, NewHeaderFlit("t", uuid.New(), packetIndex, 1, 100, route).PacketIndex())
}

func TestHeaderFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	_, _, route := testDummyRoute(t)

	assert.Equal(t, priority, NewHeaderFlit("t", uuid.New(), 0, priority, 100, route).Priority())
}

func TestHeaderFlitDeadline(t *testing.T) {
	t.Parallel()

	deadline := 100
	_, _, route := testDummyRoute(t)

	assert.Equal(t, deadline, NewHeaderFlit("t", uuid.New(), 0, 1, deadline, route).Deadline())
}

func TestHeaderFlitRoute(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit("t", uuid.New(), 0, 1, 100, route)
	assert.Equal(t, route, flit.Route())
}

func TestBodyFlitUUID(t *testing.T) {
	t.Parallel()

	flit := NewBodyFlit("t", uuid.New(), 1, 1, 1)
	assert.Equal(t, flit.id, flit.ID())
}

func TestBodyFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, BodyFlitType, NewBodyFlit("t", uuid.New(), 1, 1, 1).Type())
}

func TestBodyFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	assert.Equal(t, packetUUID, NewBodyFlit("t", packetUUID, 1, 1, 1).PacketUUID())
}

func TestBodyFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, packetIndex, NewBodyFlit("t", uuid.New(), packetIndex, 1, 1).PacketIndex())
}

func TestBodyFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewBodyFlit("t", uuid.New(), 1, 1, 1).Priority())
}

func TestBodyFlitDataSize(t *testing.T) {
	t.Parallel()

	dataSize := 4

	flit := NewBodyFlit("t", uuid.New(), 1, 1, dataSize)
	assert.Equal(t, dataSize, flit.DataSize())
}

func TestTailFlitUUID(t *testing.T) {
	t.Parallel()

	flit := NewTailFlit("t", uuid.New(), 2, 1)
	assert.Equal(t, flit.id, flit.ID())
}

func TestTailFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, TailFlitType, NewTailFlit("t", uuid.New(), 2, 1).Type())
}

func TestTailFlitPacketUUID(t *testing.T) {
	t.Parallel()

	packetUUID := uuid.New()
	assert.Equal(t, packetUUID, NewTailFlit("t", packetUUID, 2, 1).PacketUUID())
}

func TestTailFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, packetIndex, NewTailFlit("t", uuid.New(), packetIndex, 1).PacketIndex())
}

func TestTailFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewTailFlit("t", uuid.New(), 2, 1).Priority())
}
