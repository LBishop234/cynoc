package packet

import (
	"fmt"
	"io"
	"math"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPacket(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	bodySize := 4

	packet := NewPacket("t", "AA", 1, 100, route, bodySize, zerolog.New(io.Discard))
	assert.Equal(t, route, packet.route)
	assert.Equal(t, bodySize, packet.bodySize)

	assert.Implements(t, (*Packet)(nil), packet)
}

func TestPacketID(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	packet := NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, packet.packetID, packet.PacketID())
}

func TestPacketTrafficFlowID(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	trafficFlowID := "t"

	packet := NewPacket(trafficFlowID, "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, trafficFlowID, packet.TrafficFlowID())
}

func TestPacketPriority(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	priority := 1

	packet := NewPacket("t", "AA", priority, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, priority, packet.Priority())
}

func TestPacketRoute(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	packet := NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))
	assert.Equal(t, route, packet.Route())
}

func TestPacketBodySize(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	bodySize := 4

	packet := NewPacket("t", "AA", 1, 100, route, bodySize, zerolog.New(io.Discard))
	assert.Equal(t, bodySize, packet.BodySize())
}

func TestPacketFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		trafficFlowID string
		id            string
		priority      int
		deadline      int
		src           domain.NodeID
		dst           domain.NodeID
		bodySize      int
		flitSize      int
		expected      []Flit
	}{
		{
			trafficFlowID: "t",
			id:            "AA",
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 1, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 2, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 3, 1, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "AA", 4, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "BB",
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n3", Pos: domain.NewPosition(1, 1)},
			dst:           domain.NodeID{ID: "n4", Pos: domain.NewPosition(1, 2)},
			bodySize:      5,
			flitSize:      3,
			expected: []Flit{
				NewHeaderFlit("t", "BB", 0, 1, 100, domain.Route{domain.NodeID{ID: "n3", Pos: domain.NewPosition(1, 1)}, domain.NodeID{ID: "n4", Pos: domain.NewPosition(1, 2)}}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 1, 3, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 2, 2, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "BB", 3, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "CC",
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n5", Pos: domain.NewPosition(2, 2)},
			dst:           domain.NodeID{ID: "n6", Pos: domain.NewPosition(2, 3)},
			bodySize:      0,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", "CC", 0, 1, 100, domain.Route{domain.NodeID{ID: "n5", Pos: domain.NewPosition(2, 2)}, domain.NodeID{ID: "n6", Pos: domain.NewPosition(2, 3)}}, zerolog.New(io.Discard)),
				NewTailFlit("t", "CC", 1, 1, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "DD",
			priority:      5,
			deadline:      100,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", "DD", 0, 5, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 1, 5, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 2, 5, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "DD", 3, 5, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "DD", 4, 5, zerolog.New(io.Discard)),
			},
		},
		{
			trafficFlowID: "t",
			id:            "EE",
			priority:      1,
			deadline:      500,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", "EE", 0, 1, 500, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 1, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 2, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "EE", 3, 1, 1, zerolog.New(io.Discard)),
				NewTailFlit("t", "EE", 4, 1, zerolog.New(io.Discard)),
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		index := i
		testCase := testCases[index]

		t.Run(fmt.Sprintf("TestCase-%d", index), func(t *testing.T) {
			packet := NewPacket(testCase.trafficFlowID, testCase.id, testCase.priority, testCase.deadline, domain.Route{testCase.src, testCase.dst}, testCase.bodySize, zerolog.New(io.Discard))

			gotFlits := packet.Flits(testCase.flitSize)
			assert.Equal(t, len(testCase.expected), len(gotFlits))
			for i := 0; i < len(gotFlits); i++ {
				assert.Equal(t, testCase.expected[i].Type(), gotFlits[i].Type())
				assert.Equal(t, testCase.expected[i].PacketID(), gotFlits[i].PacketID())
			}
		})
	}
}

func TestPacketBodyFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id       string
		bodySize int
		flitSize int
		expected []BodyFlit
	}{
		{
			id:       "AA",
			bodySize: 3,
			flitSize: 1,
			expected: []BodyFlit{
				NewBodyFlit("t", "AA", 1, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 2, 1, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "AA", 3, 1, 1, zerolog.New(io.Discard)),
			},
		},
		{
			id:       "BB",
			bodySize: 7,
			flitSize: 2,
			expected: []BodyFlit{
				NewBodyFlit("t", "BB", 1, 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 2, 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 3, 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "BB", 4, 1, 1, zerolog.New(io.Discard)),
			},
		},
		{
			id:       "CC",
			bodySize: 4,
			flitSize: 2,
			expected: []BodyFlit{
				NewBodyFlit("t", "CC", 1, 2, 1, zerolog.New(io.Discard)),
				NewBodyFlit("t", "CC", 2, 2, 1, zerolog.New(io.Discard)),
			},
		},
		{
			bodySize: 0,
			flitSize: 1,
			expected: []BodyFlit{},
		},
	}

	for i := 0; i < len(testCases); i++ {
		index := i
		testCase := testCases[index]

		t.Run(fmt.Sprintf("TestCase-%d", index), func(t *testing.T) {
			src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
			dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
			route := domain.Route{src, dst}

			packet := NewPacket("t", testCase.id, 1, 100, route, testCase.bodySize, zerolog.New(io.Discard))
			assert.Equal(t, len(testCase.expected), len(packet.bodyFlits(testCase.flitSize)))
		})
	}
}

func TestPacketBodyFlitCount(t *testing.T) {
	t.Parallel()

	t.Run("PopulatedBody", func(t *testing.T) {
		bodySize := 10
		testCases := map[int]int{
			0:  int(math.Inf(1)),
			1:  10,
			2:  5,
			3:  4,
			4:  3,
			5:  2,
			6:  2,
			7:  2,
			8:  2,
			9:  2,
			10: 1,
		}

		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		packet := NewPacket("t", "AA", 1, 100, route, bodySize, zerolog.New(io.Discard))

		for flitSize, flitCount := range testCases {
			t.Run(fmt.Sprintf("FlitSize-%d", flitSize), func(t *testing.T) {
				assert.Equal(t, flitCount, packet.bodyFlitCount(flitSize))
			})
		}
	})
}

func TestEqualPackets(t *testing.T) {
	t.Parallel()

	t.Run("Equal", func(t *testing.T) {
		trafficFlowID := "t"
		pktID := "AA"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, pktID, priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, priority, deadline, route, data, zerolog.New(io.Discard))

		require.NoError(t, EqualPackets(pkt1, pkt2))
	})

	t.Run("IDNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, "AA", priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, "BB", priority, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("TrafficFlowIDNotEqual", func(t *testing.T) {
		pktID := "AA"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket("t1", pktID, priority, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket("t2", pktID, priority, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("PriorityNotEqual", func(t *testing.T) {
		pktID := "AA"
		trafficFlowID := "t"
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, pktID, 1, deadline, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, 2, deadline, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("DeadlineNotEqual", func(t *testing.T) {
		pktID := "AA"
		trafficFlowID := "t"
		priority := 1
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := NewPacket(trafficFlowID, pktID, priority, 1, route, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, priority, 2, route, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualLen", func(t *testing.T) {
		trafficFlowID := "t"
		pktID := "AA"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		data := 4

		pkt1 := NewPacket(trafficFlowID, pktID, priority, deadline, domain.Route{src, dst}, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, priority, deadline, domain.Route{}, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualContent", func(t *testing.T) {
		trafficFlowID := "t"
		pktID := "AA"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		data := 4

		pkt1 := NewPacket(trafficFlowID, pktID, priority, deadline, domain.Route{src, dst}, data, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, priority, deadline, domain.Route{dst, src}, data, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("BodySizeNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		pktID := "AA"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		pkt1 := NewPacket(trafficFlowID, pktID, priority, deadline, route, 4, zerolog.New(io.Discard))
		pkt2 := NewPacket(trafficFlowID, pktID, priority, deadline, route, 11, zerolog.New(io.Discard))

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("NilParameters", func(t *testing.T) {
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		require.ErrorIs(t, EqualPackets(nil, nil), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(nil, NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(NewPacket("t", "BB", 1, 100, route, 4, zerolog.New(io.Discard)), nil), domain.ErrNilParameter)
	})
}
