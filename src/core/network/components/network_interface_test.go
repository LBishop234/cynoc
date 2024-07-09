package components

import (
	"io"
	"testing"

	"main/src/domain"
	"main/src/traffic/packet"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkInterface(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.Implements(t, (*NetworkInterface)(nil), netIntfc)
	})

	t.Run("Valid", func(t *testing.T) {
		var nodeID domain.NodeID = domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}
		var bufferSize int = 1
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(nodeID, bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.NotNil(t, netIntfc)
		assert.Equal(t, nodeID, netIntfc.NodeID())
		assert.Equal(t, bufferSize, netIntfc.bufferSize)
		assert.Equal(t, maxPriority, netIntfc.maxPriority)
		assert.NotNil(t, netIntfc.flitsInTransit)
		assert.NotNil(t, netIntfc.flitsArriving)
		assert.NotNil(t, netIntfc.arrivedPackets)
	})

	t.Run("InvalidBufferSize", func(t *testing.T) {
		_, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 0, 1, zerolog.New(io.Discard))
		require.Error(t, err)
	})
}

func TestNetworkInterfaceNodeID(t *testing.T) {
	t.Parallel()

	var nodeID domain.NodeID = domain.NodeID{ID: "n", Pos: domain.NewPosition(0, 0)}

	netIntfc, err := newNetworkInterface(nodeID, 1, 1, zerolog.New(io.Discard))
	require.NoError(t, err)

	assert.Equal(t, nodeID, netIntfc.NodeID())
}

func TestNetworkInterfaceSetInputPort(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var bufferSize int = 1
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetInputPort(conn)
		require.NoError(t, err)

		assert.Equal(t, conn, netIntfc.inputPort.connection())
		assert.Equal(t, netIntfc.NodeID(), netIntfc.inputPort.connection().GetDstRouter())
	})

	t.Run("NilConnection", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetInputPort(nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})

	t.Run("InvalidBufferSize", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		netIntfc.bufferSize = 0

		err = netIntfc.SetInputPort(&connectionImpl{})
		require.Error(t, err)
	})
}

func TestNetworkInterfaceSetOutputPort(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var bufferSize int = 1
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetOutputPort(conn)
		require.NoError(t, err)

		assert.Equal(t, conn, netIntfc.outputPort.connection())
		assert.Equal(t, netIntfc.NodeID(), netIntfc.outputPort.connection().GetSrcRouter())
	})

	t.Run("NilConnection", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetOutputPort(nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestNetworkInterfaceRoutePacket(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		var src domain.NodeID = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		var dst domain.NodeID = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		var route domain.Route = domain.Route{src, dst}

		pkt := packet.NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))

		err = netIntfc.RoutePacket(0, pkt)
		require.NoError(t, err)

		for i := 0; i < len(pkt.Flits()); i++ {
			assert.Equal(t, netIntfc.flitsInTransit[pkt.Priority()][i].PacketIndex(), pkt.Flits()[i].PacketIndex())
			assert.Equal(t, netIntfc.flitsInTransit[pkt.Priority()][i].Type(), pkt.Flits()[i].Type())
		}
	})

	t.Run("NilPacket", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.RoutePacket(0, nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestNetworkInterfacePopArrivedPackets(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		var src domain.NodeID = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		var dst domain.NodeID = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		var route domain.Route = domain.Route{src, dst}

		pkts := []packet.Packet{
			packet.NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard)),
			packet.NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard)),
			packet.NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard)),
		}
		netIntfc.arrivedPackets = append(netIntfc.arrivedPackets, pkts...)

		gotPkts := netIntfc.PopArrivedPackets(0)
		assert.Equal(t, pkts, gotPkts)
		assert.Empty(t, netIntfc.arrivedPackets)
	})
}

func TestNetworkInterfaceHandleArrivingFlits(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var src domain.NodeID = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		var dst domain.NodeID = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		var route domain.Route = domain.Route{src, dst}

		var bufferSize int = 1
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(src, bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		pkt := packet.NewPacket("t", "AA", 1, 100, route, 2, zerolog.New(io.Discard))
		flits := pkt.Flits()

		inConn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetInputPort(inConn)
		require.NoError(t, err)

		for i := 1; i <= maxPriority; i++ {
			<-inConn.creditChannel(i)
		}

		for i := 0; i < len(flits); i++ {
			inConn.flitChannel() <- flits[i]

			err = netIntfc.HandleArrivingFlits(0)
			require.NoError(t, err)

			<-inConn.creditChannel(flits[i].Priority())
		}

		require.Len(t, netIntfc.arrivedPackets, 1)
		require.NoError(t, packet.EqualPackets(pkt, netIntfc.arrivedPackets[0]))
	})
}

func TestNetworkInterfaceArrivedHeaderFlit(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		headerFlit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}}, zerolog.New(io.Discard))

		err = netIntfc.arrivedHeaderFlit(headerFlit)
		require.NoError(t, err)

		assert.Contains(t, netIntfc.flitsArriving, headerFlit.PacketID())
	})

	t.Run("SetHeaderError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})
}

func TestNetworkInterfaceArrivedBodyFlit(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var pktID string = "AA"

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		bodyFlit := packet.NewBodyFlit("t", pktID, 1, 1, zerolog.New(io.Discard))

		netIntfc.flitsArriving[bodyFlit.PacketID()], err = packet.NewReconstructor(packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{}, zerolog.New(io.Discard)), zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.arrivedBodyFlit(bodyFlit)
		require.NoError(t, err)
	})

	t.Run("AddBodyError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})
}

func TestNetworkInterfaceArrivedTailFlit(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var trafficFlowID string = "t"
		var pktID string = "AA"
		var priority int = 1
		var deadline int = 100
		var src domain.NodeID = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		var dst domain.NodeID = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		var route domain.Route = domain.Route{src, dst}
		var bodySize int = 4

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		headerFlit := packet.NewHeaderFlit(trafficFlowID, pktID, 0, priority, deadline, route, zerolog.New(io.Discard))
		netIntfc.flitsArriving[headerFlit.PacketID()], err = packet.NewReconstructor(headerFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		for i := 0; i < bodySize; i++ {
			bodyFlit := packet.NewBodyFlit(trafficFlowID, pktID, i+1, priority, zerolog.New(io.Discard))
			err = netIntfc.flitsArriving[bodyFlit.PacketID()].AddBody(bodyFlit)
			require.NoError(t, err)
		}

		tailFlit := packet.NewTailFlit(trafficFlowID, pktID, bodySize+1, priority, zerolog.New(io.Discard))
		err = netIntfc.arrivedTailFlit(tailFlit)
		require.NoError(t, err)

		assert.Equal(t, pktID, netIntfc.arrivedPackets[0].PacketIndex())
		assert.Equal(t, trafficFlowID, netIntfc.arrivedPackets[0].TrafficFlowID())
		assert.Equal(t, priority, netIntfc.arrivedPackets[0].Priority())
		assert.Equal(t, deadline, netIntfc.arrivedPackets[0].Deadline())
		assert.Equal(t, route, netIntfc.arrivedPackets[0].Route())
		assert.Equal(t, bodySize, netIntfc.arrivedPackets[0].BodySize())

		assert.NotContains(t, netIntfc.flitsArriving, tailFlit.PacketID())
	})

	t.Run("SetTailError", func(t *testing.T) {
		var pktID string = "AA"

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		netIntfc.flitsArriving[pktID], err = packet.NewReconstructor(packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{}, zerolog.New(io.Discard)), zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.flitsArriving[pktID].SetTail(nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)

		tailFlit := packet.NewTailFlit("t", pktID, 2, 1, zerolog.New(io.Discard))
		err = netIntfc.arrivedTailFlit(tailFlit)
		require.Error(t, err)
	})

	t.Run("ReconstructError", func(t *testing.T) {
		var pktID string = "AA"

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		netIntfc.flitsArriving[pktID], err = packet.NewReconstructor(packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{}, zerolog.New(io.Discard)), zerolog.New(io.Discard))
		require.NoError(t, err)

		tailFlit := packet.NewTailFlit("t", pktID, 2, 1, zerolog.New(io.Discard))
		err = netIntfc.arrivedTailFlit(tailFlit)
		require.Error(t, err)
	})
}

func TestNetworkInterfaceTransmitPendingPackets(t *testing.T) {
	t.Parallel()

	t.Run("NoFlitsToTransit", func(t *testing.T) {
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, 1, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = netIntfc.SetOutputPort(conn)
		require.NoError(t, err)

		err = netIntfc.TransmitPendingPackets(0)
		require.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		var src domain.NodeID = domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		var dst domain.NodeID = domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		var route domain.Route = domain.Route{src, dst}

		var bufferSize int = 1
		var maxPriority int = 1

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(0, 0)}, bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		buff, err := newBuffer(bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)
		newInputPort(conn, buff, zerolog.New(io.Discard))

		err = netIntfc.SetOutputPort(conn)
		require.NoError(t, err)

		pkt := packet.NewPacket("t", "AA", 1, 100, route, 4, zerolog.New(io.Discard))

		err = netIntfc.RoutePacket(0, pkt)
		require.NoError(t, err)

		err = netIntfc.TransmitPendingPackets(0)
		require.NoError(t, err)

		gotFlit := <-conn.flitChan
		assert.Equal(t, pkt.Flits()[0].PacketIndex(), gotFlit.PacketIndex())
		assert.Equal(t, pkt.Flits()[0].Type(), gotFlit.Type())
	})
}
