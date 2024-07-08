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
	packetIndex := "AABBCCDD"
	flitIndex := 0
	priority := 1
	deadline := 100
	_, _, route := testDummyRoute(t)

	flit := NewHeaderFlit(trafficFlowID, packetIndex, flitIndex, priority, deadline, route, zerolog.New(io.Discard))
	assert.Equal(t, trafficFlowID, flit.trafficFlowID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, flitIndex, flit.flitIndex)
	assert.Equal(t, priority, flit.priority)
	assert.Equal(t, route, flit.route)

	assert.Implements(t, (*HeaderFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewBodyFlit(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	packetIndex := "AABBCCDD"
	flitIndex := 1
	priority := 1

	flit := NewBodyFlit(trafficFlowID, packetIndex, flitIndex, priority, zerolog.New(io.Discard))
	assert.Equal(t, newPacketID(trafficFlowID, packetIndex), flit.packetID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, flitIndex, flit.flitIndex)
	assert.Equal(t, priority, flit.priority)

	assert.Implements(t, (*BodyFlit)(nil), flit)
	assert.Implements(t, (*Flit)(nil), flit)
}

func TestNewTailFlit(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	packetIndex := "AABBCCDD"
	flitIndex := 2
	priority := 1

	flit := NewTailFlit(trafficFlowID, packetIndex, flitIndex, priority, zerolog.New(io.Discard))
	assert.Equal(t, newPacketID(trafficFlowID, packetIndex), flit.packetID)
	assert.Equal(t, packetIndex, flit.packetIndex)
	assert.Equal(t, flitIndex, flit.flitIndex)
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

	trafficFlowID := "t"
	packetIndex := "AABBCCDD"
	_, _, route := testDummyRoute(t)

	expectedPacketID := "t-AABBCCDD"

	assert.Equal(t, expectedPacketID, NewHeaderFlit(trafficFlowID, packetIndex, 0, 1, 100, route, zerolog.New(io.Discard)).PacketID())
}

func TestHeaderFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := "AABBCCDD"
	_, _, route := testDummyRoute(t)

	assert.Equal(t, packetIndex, NewHeaderFlit("t", packetIndex, 0, 1, 100, route, zerolog.New(io.Discard)).PacketIndex())
}

func TestHeaderFlitIndex(t *testing.T) {
	t.Parallel()

	flitIndex := rand.Intn(math.MaxInt)
	_, _, route := testDummyRoute(t)

	assert.Equal(t, flitIndex, NewHeaderFlit("t", "AABBCCDD", flitIndex, 1, 100, route, zerolog.New(io.Discard)).FlitIndex())
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

	flit := NewBodyFlit("t", "AABBCCDD", 1, 1, zerolog.New(io.Discard))
	assert.Equal(t, flit.id, flit.ID())
}

func TestBodyFlitType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, BodyFlitType, NewBodyFlit("t", "AABBCCDD", 1, 1, zerolog.New(io.Discard)).Type())
}

func TestBodyFlitPacketID(t *testing.T) {
	t.Parallel()

	trafficFlowID := "t"
	packetIndex := "AABBCCDD"

	expectedPacketID := "t-AABBCCDD"

	assert.Equal(t, expectedPacketID, NewBodyFlit(trafficFlowID, packetIndex, 1, 1, zerolog.New(io.Discard)).PacketID())
}

func TestBodyFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := "AABBCCDD"

	assert.Equal(t, packetIndex, NewBodyFlit("t", packetIndex, 1, 1, zerolog.New(io.Discard)).PacketIndex())
}

func TestBodyFlitIndex(t *testing.T) {
	t.Parallel()

	flitIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, flitIndex, NewBodyFlit("t", "AABBCCDD", flitIndex, 1, zerolog.New(io.Discard)).FlitIndex())
}

func TestBodyFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewBodyFlit("t", "AABBCCDD", 1, 1, zerolog.New(io.Discard)).Priority())
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

	trafficFlowID := "t"
	packetIndex := "AABBCCDD"

	expectedPacketID := "t-AABBCCDD"

	assert.Equal(t, expectedPacketID, NewTailFlit(trafficFlowID, packetIndex, 2, 1, zerolog.New(io.Discard)).PacketID())
}

func TestTailFlitPacketIndex(t *testing.T) {
	t.Parallel()

	packetIndex := "AABBCCDD"
	assert.Equal(t, packetIndex, NewTailFlit("t", packetIndex, 2, 1, zerolog.New(io.Discard)).PacketIndex())
}

func TestTailFlitflitIndex(t *testing.T) {
	t.Parallel()

	flitIndex := rand.Intn(math.MaxInt)
	assert.Equal(t, flitIndex, NewTailFlit("t", "AABBCCDD", flitIndex, 1, zerolog.New(io.Discard)).FlitIndex())
}

func TestTailFlitPriority(t *testing.T) {
	t.Parallel()

	priority := 1
	assert.Equal(t, priority, NewTailFlit("t", "AABBCCDD", 2, 1, zerolog.New(io.Discard)).Priority())
}
