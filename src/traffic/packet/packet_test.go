package packet

import (
	"fmt"
	"math"
	"testing"

	"main/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPacket(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	bodySize := 4

	packet := NewPacket("t", 1, 100, route, bodySize)
	assert.Equal(t, route, packet.route)
	assert.Equal(t, bodySize, packet.bodySize)

	assert.Implements(t, (*Packet)(nil), packet)
}

func TestPacketUUID(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	packet := NewPacket("t", 1, 100, route, 4)
	assert.Equal(t, packet.uuid, packet.UUID())
}

func TestPacketTrafficFlowID(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	trafficFlowID := "t"

	packet := NewPacket(trafficFlowID, 1, 100, route, 4)
	assert.Equal(t, trafficFlowID, packet.TrafficFlowID())
}

func TestPacketPriority(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	priority := 1

	packet := NewPacket("t", priority, 100, route, 4)
	assert.Equal(t, priority, packet.Priority())
}

func TestPacketRoute(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}

	packet := NewPacket("t", 1, 100, route, 4)
	assert.Equal(t, route, packet.Route())
}

func TestPacketBodySize(t *testing.T) {
	t.Parallel()

	src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
	dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
	route := domain.Route{src, dst}
	bodySize := 4

	packet := NewPacket("t", 1, 100, route, bodySize)
	assert.Equal(t, bodySize, packet.BodySize())
}

func TestPacketFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		trafficFlowID string
		uuid          uuid.UUID
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
			uuid:          uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"),
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 1, 1, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 2, 1, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 3, 1, 1),
				NewTailFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 4, 1),
			},
		},
		{
			trafficFlowID: "t",
			uuid:          uuid.MustParse("25ce3554-f19b-4d7c-8153-cee9023fb291"),
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n3", Pos: domain.NewPosition(1, 1)},
			dst:           domain.NodeID{ID: "n4", Pos: domain.NewPosition(1, 2)},
			bodySize:      5,
			flitSize:      3,
			expected: []Flit{
				NewHeaderFlit("t", uuid.MustParse("25ce3554-f19b-4d7c-8153-cee9023fb291"), 0, 1, 100, domain.Route{domain.NodeID{ID: "n3", Pos: domain.NewPosition(1, 1)}, domain.NodeID{ID: "n4", Pos: domain.NewPosition(1, 2)}}),
				NewBodyFlit("t", uuid.MustParse("25ce3554-f19b-4d7c-8153-cee9023fb291"), 1, 3, 1),
				NewBodyFlit("t", uuid.MustParse("25ce3554-f19b-4d7c-8153-cee9023fb291"), 2, 2, 1),
				NewTailFlit("t", uuid.MustParse("25ce3554-f19b-4d7c-8153-cee9023fb291"), 3, 1),
			},
		},
		{
			trafficFlowID: "t",
			uuid:          uuid.MustParse("e45c2547-0d59-4586-b87f-fcacdf507983"),
			priority:      1,
			deadline:      100,
			src:           domain.NodeID{ID: "n5", Pos: domain.NewPosition(2, 2)},
			dst:           domain.NodeID{ID: "n6", Pos: domain.NewPosition(2, 3)},
			bodySize:      0,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", uuid.MustParse("e45c2547-0d59-4586-b87f-fcacdf507983"), 0, 1, 100, domain.Route{domain.NodeID{ID: "n5", Pos: domain.NewPosition(2, 2)}, domain.NodeID{ID: "n6", Pos: domain.NewPosition(2, 3)}}),
				NewTailFlit("t", uuid.MustParse("e45c2547-0d59-4586-b87f-fcacdf507983"), 1, 1),
			},
		},
		{
			trafficFlowID: "t",
			uuid:          uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"),
			priority:      5,
			deadline:      100,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 0, 5, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 1, 5, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 2, 5, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 3, 5, 1),
				NewTailFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 4, 5),
			},
		},
		{
			trafficFlowID: "t",
			uuid:          uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"),
			priority:      1,
			deadline:      500,
			src:           domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)},
			dst:           domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)},
			bodySize:      3,
			flitSize:      1,
			expected: []Flit{
				NewHeaderFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 0, 1, 500, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 1, 1, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 2, 1, 1),
				NewBodyFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 3, 1, 1),
				NewTailFlit("t", uuid.MustParse("87a6823b-28e6-4148-91e1-371924a7e20b"), 4, 1),
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		index := i
		testCase := testCases[index]

		t.Run(fmt.Sprintf("TestCase-%d", index), func(t *testing.T) {
			packet := newPacketWithUUID(testCase.trafficFlowID, testCase.uuid, testCase.priority, testCase.deadline, domain.Route{testCase.src, testCase.dst}, testCase.bodySize)

			gotFlits := packet.Flits(testCase.flitSize)
			assert.Equal(t, len(testCase.expected), len(gotFlits))
			for i := 0; i < len(gotFlits); i++ {
				assert.Equal(t, testCase.expected[i].Type(), gotFlits[i].Type())
				assert.Equal(t, testCase.expected[i].PacketUUID(), gotFlits[i].PacketUUID())
			}
		})
	}
}

func TestPacketBodyFlits(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		uuid     uuid.UUID
		bodySize int
		flitSize int
		expected []BodyFlit
	}{
		{
			uuid:     uuid.MustParse("26612df6-b7ea-4315-8445-73e02a8bb19a"),
			bodySize: 3,
			flitSize: 1,
			expected: []BodyFlit{
				NewBodyFlit("t", uuid.MustParse("26612df6-b7ea-4315-8445-73e02a8bb19a"), 1, 1, 1),
				NewBodyFlit("t", uuid.MustParse("26612df6-b7ea-4315-8445-73e02a8bb19a"), 2, 1, 1),
				NewBodyFlit("t", uuid.MustParse("26612df6-b7ea-4315-8445-73e02a8bb19a"), 3, 1, 1),
			},
		},
		{
			uuid:     uuid.MustParse("ed8d87c1-e1e0-4189-a658-e3017e9d1ebe"),
			bodySize: 7,
			flitSize: 2,
			expected: []BodyFlit{
				NewBodyFlit("t", uuid.MustParse("ed8d87c1-e1e0-4189-a658-e3017e9d1ebe"), 1, 2, 1),
				NewBodyFlit("t", uuid.MustParse("ed8d87c1-e1e0-4189-a658-e3017e9d1ebe"), 2, 2, 1),
				NewBodyFlit("t", uuid.MustParse("ed8d87c1-e1e0-4189-a658-e3017e9d1ebe"), 3, 2, 1),
				NewBodyFlit("t", uuid.MustParse("ed8d87c1-e1e0-4189-a658-e3017e9d1ebe"), 4, 1, 1),
			},
		},
		{
			uuid:     uuid.MustParse("64687f40-6a02-43d4-8c34-5131b3e26f9e"),
			bodySize: 4,
			flitSize: 2,
			expected: []BodyFlit{
				NewBodyFlit("t", uuid.MustParse("64687f40-6a02-43d4-8c34-5131b3e26f9e"), 1, 2, 1),
				NewBodyFlit("t", uuid.MustParse("64687f40-6a02-43d4-8c34-5131b3e26f9e"), 2, 2, 1),
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

			packet := NewPacket("t", 1, 100, route, testCase.bodySize)
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

		packet := NewPacket("t", 1, 100, route, bodySize)

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
		pktUUID := uuid.New()
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, route, data)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, route, data)

		require.NoError(t, EqualPackets(pkt1, pkt2))
	})

	t.Run("UUIDNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, uuid.New(), priority, deadline, route, data)
		pkt2 := newPacketWithUUID(trafficFlowID, uuid.New(), priority, deadline, route, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("TrafficFlowIDNotEqual", func(t *testing.T) {
		pktUUID := uuid.New()
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := newPacketWithUUID("t1", pktUUID, priority, deadline, route, data)
		pkt2 := newPacketWithUUID("t2", pktUUID, priority, deadline, route, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("PriorityNotEqual", func(t *testing.T) {
		pktUUID := uuid.New()
		trafficFlowID := "t"
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, 1, deadline, route, data)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, 2, deadline, route, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("DeadlineNotEqual", func(t *testing.T) {
		pktUUID := uuid.New()
		trafficFlowID := "t"
		priority := 1
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, priority, 1, route, data)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, priority, 2, route, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualLen", func(t *testing.T) {
		trafficFlowID := "t"
		pktUUID := uuid.New()
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, domain.Route{src, dst}, data)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, domain.Route{}, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("RouteNotEqualContent", func(t *testing.T) {
		trafficFlowID := "t"
		pktUUID := uuid.New()
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		data := 4

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, domain.Route{src, dst}, data)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, domain.Route{dst, src}, data)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("BodySizeNotEqual", func(t *testing.T) {
		trafficFlowID := "t"
		pktUUID := uuid.New()
		priority := 1
		deadline := 100
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		pkt1 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, route, 4)
		pkt2 := newPacketWithUUID(trafficFlowID, pktUUID, priority, deadline, route, 11)

		require.ErrorIs(t, EqualPackets(pkt1, pkt2), domain.ErrPacketsNotEqual)
	})

	t.Run("NilParameters", func(t *testing.T) {
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		require.ErrorIs(t, EqualPackets(nil, nil), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(nil, newPacketWithUUID("t", uuid.New(), 1, 100, route, 4)), domain.ErrNilParameter)
		require.ErrorIs(t, EqualPackets(newPacketWithUUID("t", uuid.New(), 1, 100, route, 4), nil), domain.ErrNilParameter)
	})
}
