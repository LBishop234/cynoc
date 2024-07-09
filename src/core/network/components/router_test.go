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

func testRouter(t *testing.T) *routerImpl {
	router, err := newRouter(
		RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n",
				Pos: domain.NewPosition(0, 0),
			},
			SimConfig: domain.SimConfig{
				BufferSize:      1,
				MaxPriority:     1,
				ProcessingDelay: 1,
				LinkBandwidth:   1,
			},
		},
		zerolog.New(io.Discard).With().Logger(),
	)
	require.NoError(t, err)

	return router
}

type testRouterPair struct {
	rA   *routerImpl
	niA  *networkInterfaceImpl
	rB   *routerImpl
	niB  *networkInterfaceImpl
	AtoB *connectionImpl
	BtoA *connectionImpl
}

func newTestRouterPair(t *testing.T, bufferSize, processingDelay, maxPriority, linkBandwidth int) testRouterPair {
	aPos := domain.NewPosition(0, 0)
	bPos := domain.NewPosition(1, 0)

	rA, err := newRouter(
		RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n-a",
				Pos: aPos,
			},
			SimConfig: domain.SimConfig{
				BufferSize:      bufferSize,
				ProcessingDelay: processingDelay,
				MaxPriority:     maxPriority,
				LinkBandwidth:   linkBandwidth,
			},
		},
		zerolog.New(io.Discard).With().Logger(),
	)
	require.NoError(t, err)

	niA, err := newNetworkInterface(domain.NodeID{ID: "i-a", Pos: aPos}, bufferSize, maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	err = rA.SetNetworkInterface(niA)
	require.NoError(t, err)

	rB, err := newRouter(
		RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n-b",
				Pos: bPos,
			},
			SimConfig: domain.SimConfig{
				BufferSize:      bufferSize,
				ProcessingDelay: processingDelay,
				MaxPriority:     maxPriority,
				LinkBandwidth:   linkBandwidth,
			},
		},
		zerolog.New(io.Discard).With().Logger(),
	)
	require.NoError(t, err)

	niB, err := newNetworkInterface(domain.NodeID{ID: "i-b", Pos: bPos}, bufferSize, maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	err = rB.SetNetworkInterface(niB)
	require.NoError(t, err)

	AtoB, err := NewConnection(maxPriority, linkBandwidth, zerolog.New(io.Discard))
	require.NoError(t, err)

	rA.RegisterOutputPort(AtoB)
	rB.RegisterInputPort(AtoB)

	BtoA, err := NewConnection(maxPriority, linkBandwidth, zerolog.New(io.Discard))
	require.NoError(t, err)

	rB.RegisterOutputPort(BtoA)
	rA.RegisterInputPort(BtoA)

	rA.UpdateOutputMap()
	rB.UpdateOutputMap()

	return testRouterPair{
		rA:   rA,
		niA:  niA,
		rB:   rB,
		niB:  niB,
		AtoB: AtoB,
		BtoA: BtoA,
	}
}

func TestNewRouter(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		router := testRouter(t)
		assert.Implements(t, (*Router)(nil), router)
	})

	t.Run("Valid", func(t *testing.T) {
		conf := RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n",
				Pos: domain.NewPosition(0, 0),
			},
			SimConfig: domain.SimConfig{
				BufferSize:      1,
				ProcessingDelay: 1,
				MaxPriority:     1,
				LinkBandwidth:   1,
			},
		}

		router, err := newRouter(conf, zerolog.New(io.Discard).With().Logger())
		require.NoError(t, err)

		assert.Equal(t, conf.NodeID, router.nodeID)

		assert.NotNil(t, router.inputPorts)
		assert.NotNil(t, router.outputPorts)

		assert.Equal(t, conf.SimConfig, router.simConf)

		assert.NotNil(t, router.headerFlitsProcessings)
	})

	t.Run("InvalidBufferSize", func(t *testing.T) {
		_, err := newRouter(
			RouterConfig{
				SimConfig: domain.SimConfig{
					BufferSize:      0,
					ProcessingDelay: 1,
					MaxPriority:     1,
					LinkBandwidth:   1,
				},
			},
			zerolog.New(io.Discard).With().Logger(),
		)
		require.Error(t, err)
	})

	t.Run("InvalidProcessingDelay", func(t *testing.T) {
		_, err := newRouter(
			RouterConfig{
				SimConfig: domain.SimConfig{
					BufferSize:      1,
					ProcessingDelay: 0,
					MaxPriority:     1,
					LinkBandwidth:   1,
				},
			},
			zerolog.New(io.Discard).With().Logger(),
		)
		require.ErrorIs(t, err, domain.ErrInvalidParameter)
	})
}

func TestRouterNodeID(t *testing.T) {
	t.Parallel()

	var nodeID domain.NodeID = domain.NodeID{ID: "n", Pos: domain.NewPosition(0, 0)}

	router, err := newRouter(
		RouterConfig{
			NodeID: nodeID,
			SimConfig: domain.SimConfig{
				BufferSize:      1,
				ProcessingDelay: 1,
				MaxPriority:     1,
				LinkBandwidth:   1,
			},
		},
		zerolog.New(io.Discard).With().Logger(),
	)
	require.NoError(t, err)

	assert.Equal(t, nodeID, router.NodeID())
}

func TestRouterRegisterInputPort(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		router := testRouter(t)

		conn, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = router.RegisterInputPort(conn)
		require.NoError(t, err)
		assert.Equal(t, router.inputPorts[0].connection(), conn)
	})

	t.Run("NewInputPortError", func(t *testing.T) {
		router := testRouter(t)

		err := router.RegisterInputPort(nil)
		require.Error(t, err)
	})

	t.Run("NewBufferError", func(t *testing.T) {
		router := testRouter(t)

		router.simConf.BufferSize = 0

		conn, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = router.RegisterInputPort(conn)
		require.Error(t, err)
	})
}

func TestRouterRegisterOutputPort(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		router := testRouter(t)

		conn, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = router.RegisterOutputPort(conn)
		require.NoError(t, err)
		assert.Equal(t, router.outputPorts[0].connection(), conn)
	})

	t.Run("NewOutputPortError", func(t *testing.T) {
		router := testRouter(t)

		err := router.RegisterOutputPort(nil)
		require.Error(t, err)
	})
}

func TestRouterUpdateOutputMap(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		router := testRouter(t)

		conn1, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
		require.NoError(t, err)
		port1, err := newOutputPort(conn1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)
		nodeID1 := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 1)}
		conn1.SetDstRouter(nodeID1)
		router.outputPorts = append(router.outputPorts, port1)

		conn2, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
		require.NoError(t, err)
		port2, err := newOutputPort(conn2, 1, zerolog.New(io.Discard))
		require.NoError(t, err)
		nodeID2 := domain.NodeID{ID: "n2", Pos: domain.NewPosition(1, 0)}
		conn2.SetDstRouter(nodeID2)
		router.outputPorts = append(router.outputPorts, port2)

		router.UpdateOutputMap()
		assert.Equal(t, port1, router.outputMap[nodeID1])
		assert.Equal(t, port2, router.outputMap[nodeID2])
	})
}

func TestRouterSetNetworkInterface(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		router := testRouter(t)

		netIntfc, err := newNetworkInterface(domain.NodeID{ID: "i", Pos: domain.NewPosition(1, 1)}, 1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = router.SetNetworkInterface(netIntfc)
		require.NoError(t, err)
	})

	t.Run("NilNetworkInterface", func(t *testing.T) {
		router := testRouter(t)

		err := router.SetNetworkInterface(nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})

	t.Run("InConnNewConnectionError", func(t *testing.T) {
		router := testRouter(t)

		router.simConf.BufferSize = 0
		err := router.SetNetworkInterface(&networkInterfaceImpl{})
		require.Error(t, err)
	})

	t.Run("InConnSetOutputPortError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("InConnRegisterInputPortError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("OutConnNewConnectionError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("OutConnSetOutputPortError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("OutConnRegisterInputPortError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})
}

func TestRouterUpdateOutputPortsCredit(t *testing.T) {
	t.Parallel()

	router := testRouter(t)

	conn, err := NewConnection(router.simConf.MaxPriority, router.simConf.LinkBandwidth, zerolog.New(io.Discard))
	require.NoError(t, err)
	router.RegisterOutputPort(conn)

	conn.creditChannel(1) <- 1

	err = router.UpdateOutputPortsCredit()
	require.NoError(t, err)
}

func TestRouterRouteBufferedFlits(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		testRouterPair := newTestRouterPair(t, 1, 1, 1, 1)

		pkt := packet.NewPacket("t", "AA", 1, 100, domain.Route{testRouterPair.rA.NodeID(), testRouterPair.rB.NodeID()}, 10, zerolog.New(io.Discard))

		err := testRouterPair.niA.RoutePacket(0, pkt)
		require.NoError(t, err)

		flits := pkt.Flits()
		for i := 0; i < len(flits); i++ {
			err = testRouterPair.rA.UpdateOutputPortsCredit()
			require.NoError(t, err)

			err = testRouterPair.niA.TransmitPendingPackets(0)
			require.NoError(t, err)

			err = testRouterPair.rA.ReadFromInputPorts(0)
			require.NoError(t, err)

			err = testRouterPair.rA.RouteBufferedFlits(0)
			require.NoError(t, err)

			gotFlit := <-testRouterPair.AtoB.flitChannel()
			assert.Equal(t, pkt.Flits()[i].PacketIndex(), gotFlit.PacketIndex())
			assert.Equal(t, pkt.Flits()[i].Type(), gotFlit.Type())

			testRouterPair.AtoB.creditChannel(flits[i].Priority()) <- 1
		}
	})
}

func TestRouterReadFromInputPorts(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		testRouterPair := newTestRouterPair(t, 1, 1, 1, 1)

		flit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{testRouterPair.rA.NodeID(), testRouterPair.rB.NodeID()}, zerolog.New(io.Discard))
		testRouterPair.AtoB.flitChannel() <- flit

		err := testRouterPair.rB.ReadFromInputPorts(0)
		require.NoError(t, err)

		gotFlit, exists := testRouterPair.rB.inputPorts[1].peakBuffer(flit.Priority())
		assert.True(t, exists)
		assert.Equal(t, flit, gotFlit)
	})

	t.Run("ReadIntoBufferError", func(t *testing.T) {
		testRouterPair := newTestRouterPair(t, 1, 1, 1, 1)

		flit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{testRouterPair.rA.NodeID(), testRouterPair.rB.NodeID()}, zerolog.New(io.Discard))
		testRouterPair.AtoB.flitChannel() <- flit

		err := testRouterPair.rB.ReadFromInputPorts(0)
		require.NoError(t, err)

		testRouterPair.AtoB.flitChannel() <- flit

		err = testRouterPair.rB.ReadFromInputPorts(1)
		require.Error(t, err)
	})
}

func TestRouterRouteFlit(t *testing.T) {
	t.Parallel()

	t.Run("MissingPacketsTransmitting", func(t *testing.T) {
		router := testRouter(t)

		flit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{}, zerolog.New(io.Discard))
		_, err := router.routeFlit(flit)
		require.ErrorIs(t, err, domain.ErrNoPort)
	})
}
