package packet

import (
	"fmt"
	"io"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPacket(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}
	packetSize := 4

	packet := NewPacket("t", "AA", 1, 100, route, packetSize, zerolog.New(io.Discard))
	assert.Equal(t, route, packet.route)
	assert.Equal(t, packetSize, packet.packetSize)

	assert.Implements(t, (*Packet)(nil), packet)
}

func TestPacketID(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}

	packet := NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, packet.packetIndex, packet.PacketIndex())
}

func TestPacketTrafficFlowID(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}
	trafficFlowID := "t"

	packet := NewPacket(trafficFlowID, "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, trafficFlowID, packet.TrafficFlowID())
}

func TestPacketPriority(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}
	priority := 1

	packet := NewPacket("t", "AA", priority, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, priority, packet.Priority())
}

func TestPacketRoute(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}

	packet := NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, route, packet.Route())
}

func TestPacketPacketSize(t *testing.T) {
	t.Parallel()

	src := "n1"
	dst := "n2"
	route := domain.Route{src, dst}
	packetSize := 4

	packet := NewPacket("t", "AA", 1, 100, route, packetSize, zerolog.New(io.Discard))
	assert.Equal(t, packetSize, packet.PacketSize())
}

func TestPacketFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		trafficFlowID string
		id            string
		priority      int
		deadline      int
		src           string
		dst           string
		packetSize    int
		expected      []Flit
	}{
		{
			trafficFlowID: "t",
			id:            "AA",
			priority:      1,
			deadline:      100,
			src:           "n1",
			dst:           "n2",
			packetSize:    5,
			expected: []Flit{
				NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{"n1", "n2"}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 3, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "AA", 4, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "BB",
			priority:      1,
			deadline:      100,
			src:           "n3",
			dst:           "n4",
			packetSize:    4,
			expected: []Flit{
				NewHeaderFlit("t", "BB", 0, 1, 100, domain.Route{"n3", "n4"}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 1, 3, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 2, 2, zerolog.New(io.Discard)),
				NewTailFlit("t", "BB", 3, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "CC",
			priority:      1,
			deadline:      100,
			src:           "n5",
			dst:           "n6",
			packetSize:    2,
			expected: []Flit{
				NewHeaderFlit("t", "CC", 0, 1, 100, domain.Route{"n5", "n6"}, zerolog.New(io.Discard)),
				NewTailFlit("t", "CC", 1, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "DD",
			priority:      5,
			deadline:      100,
			src:           "n1",
			dst:           "n2",
			packetSize:    5,
			expected: []Flit{
				NewHeaderFlit("t", "DD", 0, 5, 100, domain.Route{"n1", "n2"}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 1, 5, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 2, 5, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 3, 5, zerolog.New(io.Discard)),
				NewTailFlit("t", "DD", 4, 5, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "EE",
			priority:      1,
			deadline:      500,
			src:           "n1",
			dst:           "n2",
			packetSize:    5,
			expected: []Flit{
				NewHeaderFlit("t", "EE", 0, 1, 500, domain.Route{"n1", "n2"}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 3, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "EE", 4, 1, zerolog.New(io.Discard)),
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		index := i
		testCase := testCases[index]

		t.Run(fmt.Sprintf("TestCase-%d", index), func(t *testing.T) {
			packet := NewPacket(testCase.trafficFlowID, testCase.id, testCase.priority, testCase.deadline, domain.Route{testCase.src, testCase.dst}, testCase.packetSize, zerolog.New(io.Discard))

			gotFlits := packet.Flits()
			assert.Equal(t, len(testCase.expected), len(gotFlits))
			for i := 0; i < len(gotFlits); i++ {
				assert.Equal(t, testCase.expected[i].Type(), gotFlits[i].Type())
				assert.Equal(t, testCase.expected[i].PacketIndex(), gotFlits[i].PacketIndex())
			}
		})
	}
}

func TestPacketBodyFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id         string
		packetSize int
		expected   []BodyFlit
	}{
		{
			id:         "3BodyFlits",
			packetSize: 5,
			expected: []BodyFlit{
				NewBodyFlit("t", "AA", 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 3, 1, zerolog.New(io.Discard)),
			},
		},
		{
			id:         "NoBody",
			packetSize: 2,
			expected:   []BodyFlit{},
		},
	}

	for i := 0; i < len(testCases); i++ {
		index := i
		testCase := testCases[index]

		t.Run("TestCase-"+testCase.id, func(t *testing.T) {
			src := "n1"
			dst := "n2"
			route := domain.Route{src, dst}

			packet := NewPacket("t", testCase.id, 1, 100, route, testCase.packetSize, zerolog.New(io.Discard))
			assert.Equal(t, len(testCase.expected), len(packet.bodyFlits()))
		})
	}
}

func TestEqualPackets(t *testing.T) {
	t.Parallel()

	t.Run("Equal", func(t *testing.T) {
		trafficFlowID := "t"
		packetIndex := "AA"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, packetIndex, priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, priority, deadline, route, data, zerolog.New(io.Discard))

		require.NoError(t, EqualPackets(pkt1, pkt2))
	})

	t.Run("IDNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, "AA", priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, "BB", priority, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("TrafficFlowIDNotEqual", func(t *testing.T) {
		packetIndex := "AA"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket("t1", packetIndex, priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket("t2", packetIndex, priority, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("PriorityNotEqual", func(t *testing.T) {
		packetIndex := "AA"
		trafficFlowID := "t"
		deadline := 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, packetIndex, 1, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, 2, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("DeadlineNotEqual", func(t *testing.T) {
		packetIndex := "AA"
		trafficFlowID := "t"
		priority := 1
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, packetIndex, priority, 1, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, priority, 2, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualLen", func(t *testing.T) {
		trafficFlowID := "t"
		packetIndex := "AA"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		data := 4

		pkt1 := NewPacket(trafficFlowID, packetIndex, priority, deadline, domain.Route{src, dst}, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, priority, deadline, domain.Route{}, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualContent", func(t *testing.T) {
		trafficFlowID := "t"
		packetIndex := "AA"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		data := 4

		pkt1 := NewPacket(trafficFlowID, packetIndex, priority, deadline, domain.Route{src, dst}, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, priority, deadline, domain.Route{dst, src}, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("BodySizeNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		packetIndex := "AA"
		priority := 1
		deadline := 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}

		pkt1 := NewPacket(trafficFlowID, packetIndex, priority, deadline, route, 4, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, packetIndex, priority, deadline, route, 11, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("NilParameters", func(t *testing.T) {
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}

		require.ErrorIs(t, EqualPackets(nil, nil), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(nil, NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(NewPacket("t", "BB", 1, 100, route, 4, zerolog.New(io.Discard)), nil), domain.ErrNilParameter)
	})
}
