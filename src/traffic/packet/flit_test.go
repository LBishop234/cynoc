package packet

import (
	"io"
	"math"
	"math/rand"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
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
	packetID := "AABBCCDD"
	packetIndex := 0
	priority := 1
	deadline := 100
	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit(trafficFlowID, packetID, packetIndex, priority, deadline, route, zerolog.New(io.Discard))
	assert.Equal(t, trafficFlowID, flit.trafficFlowID)
	assert.Equal(t, packetID, flit.packetID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, route, flit.route)

	assert.Implements(t, (*HeaderFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewBodyFlit(t *testing.T) {
	t.Parallel()

	packetID := "AABBCCDD"
	packetIndex := 1
	priority := 1
	dataSize := 4

	flit := NewBodyFlit("t", packetID, packetIndex, priority, dataSize, zerolog.New(io.Discard))
	assert.Equal(t, packetID, flit.packetID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, dataSize, flit.dataSize)

	assert.Implements(t, (*BodyFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewTailFlit(t *testing.T) {
	t.Parallel()

	packetID := "AABBCCDD"
	packetIndex := 2
	priority := 1

	flit := NewTailFlit("t", packetID, packetIndex, priority, zerolog.New(io.Discard))
	assert.Equal(t, packetID, flit.packetID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, priority, flit.priority)

	assert.Implements(t, (*TailFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestFlitID(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit("t", "AABBCCDD", 0, 1, 100, route, zerolog.New(io.Discard))
	assert.Equal(t, flit.id, flit.ID())
}

func TestHeaderFlitType(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	assert.Equal(t, HeaderFlitType, NewHeaderFlit("t", "AABBCCDD", 0, 1, 100, route, zerolog.New(io.Discard)).Type())
}

func TestHeaderFlitTrafficFlowID(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	_, _, route := testDummyRoute(t)

	assert.Equal(t, trafficFlowID, NewHeaderFlit(trafficFlowID, "AABBCCDD", 0, 1, 100, route, zerolog.New(io.Discard)).TrafficFlowID())
}

func TestHeaderFlitPacketID(t *testing.T) {
	t.Parallel()

	packetID := "AABBCCDD"
	_, _, route := testDummyRoute(t)

	assert.Equal(t, packetID, NewHeaderFlit("t", packetID, 0, 1, 100, route, zerolog.New(io.Discard)).PacketID())
}

func TestHeaderFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	_, _, route := testDummyRoute(t)

	assert.Equal(t, packetIndex, NewHeaderFlit("t", "AABBCCDD", packetIndex, 1, 100, route, zerolog.New(io.Discard)).PacketIndex())
}

func TestHeaderFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	_, _, route := testDummyRoute(t)

	assert.Equal(t, priority, NewHeaderFlit("t", "AABBCCDD", 0, priority, 100, route, zerolog.New(io.Discard)).Priority())
}

func TestHeaderFlitDeadline(t *testing.T) {
	t.Parallel()

	deadline := 100
	_, _, route := testDummyRoute(t)

	assert.Equal(t, deadline, NewHeaderFlit("t", "AABBCCDD", 0, 1, deadline, route, zerolog.New(io.Discard)).Deadline())
}

func TestHeaderFlitRoute(t *testing.T) {
	t.Parallel()

	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit("t", "AABBCCDD", 0, 1, 100, route, zerolog.New(io.Discard))
	assert.Equal(t, route, flit.Route())
}

func TestBodyFlitID(t *testing.T) {
	t.Parallel()

	flit := NewBodyFlit("t", "AABBCCDD", 1, 1, 1, zerolog.New(io.Discard))
	assert.Equal(t, flit.id, flit.ID())
}

func TestBodyFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, BodyFlitType, NewBodyFlit("t", "AABBCCDD", 1, 1, 1, zerolog.New(io.Discard)).Type())
}

func TestBodyFlitPacketID(t *testing.T) {
	t.Parallel()

	packetID := "AABBCCDD"
	assert.Equal(t, packetID, NewBodyFlit("t", packetID, 1, 1, 1, zerolog.New(io.Discard)).PacketID())
}

func TestBodyFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, packetIndex, NewBodyFlit("t", "AABBCCDD", packetIndex, 1, 1, zerolog.New(io.Discard)).PacketIndex())
}

func TestBodyFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewBodyFlit("t", "AABBCCDD", 1, 1, 1, zerolog.New(io.Discard)).Priority())
}

func TestBodyFlitDataSize(t *testing.T) {
	t.Parallel()

	dataSize := 4

	flit := NewBodyFlit("t", "AABBCCDD", 1, 1, dataSize, zerolog.New(io.Discard))
	assert.Equal(t, dataSize, flit.DataSize())
}

func TestTailFlitID(t *testing.T) {
	t.Parallel()

	flit := NewTailFlit("t", "AABBCCDD", 2, 1, zerolog.New(io.Discard))
	assert.Equal(t, flit.id, flit.ID())
}

func TestTailFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, TailFlitType, NewTailFlit("t", "AABBCCDD", 2, 1, zerolog.New(io.Discard)).Type())
}

func TestTailFlitPacketID(t *testing.T) {
	t.Parallel()

	packetID := "AABBCCDD"
	assert.Equal(t, packetID, NewTailFlit("t", packetID, 2, 1, zerolog.New(io.Discard)).PacketID())
}

func TestTailFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, packetIndex, NewTailFlit("t", "AABBCCDD", packetIndex, 1, zerolog.New(io.Discard)).PacketIndex())
}

func TestTailFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewTailFlit("t", "AABBCCDD", 2, 1, zerolog.New(io.Discard)).Priority())
}
