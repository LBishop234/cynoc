package components

import (
	"fmt"
	"testing"

	"main/src/domain"
	"main/src/traffic/packet"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInputPort(t *testing.T, bufferSize, maxPriority, linkBandwidth int) *inputPortImpl {
	conn, err := NewConnection(maxPriority, linkBandwidth)
	require.NoError(t, err)
	buff, err := newBuffer(bufferSize, 1)
	require.NoError(t, err)

	port, err := newInputPort(conn, buff)
	require.NoError(t, err)

	for _, credChan := range conn.creditChannels() {
		<-credChan
	}

	return port
}

func testOutputPort(t *testing.T, maxPriority, linkBandwidth int) *outputPortImpl {
	conn, err := NewConnection(maxPriority, linkBandwidth)
	require.NoError(t, err)

	port, err := newOutputPort(conn, maxPriority)
	require.NoError(t, err)

	return port
}

func TestNewInputPort(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		port := testInputPort(t, 1, 1, 1)

		assert.Implements(t, (*inputPort)(nil), port)
	})

	t.Run("Valid", func(t *testing.T) {
		var bufferSize int = 32
		var maxPriority int = 4
		var linkBandwidth int = 4

		buff, err := newBuffer(bufferSize, maxPriority)
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, linkBandwidth)
		require.NoError(t, err)

		port, err := newInputPort(conn, buff)
		require.NoError(t, err)

		assert.Equal(t, conn, port.conn)
		assert.Equal(t, buff, port.buff)
	})

	t.Run("NilConnection", func(t *testing.T) {
		buff, err := newBuffer(1, 1)
		require.NoError(t, err)

		_, err = newInputPort(nil, buff)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})

	t.Run("NilBuffer", func(t *testing.T) {
		var maxPriority int = 1
		var linkBandwidth int = 1

		conn, err := NewConnection(maxPriority, linkBandwidth)
		require.NoError(t, err)

		_, err = newInputPort(conn, nil)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestNewOutputPort(t *testing.T) {
	t.Parallel()

	t.Run("ValidInterface", func(t *testing.T) {
		port := testOutputPort(t, 1, 1)

		assert.Implements(t, (*outputPort)(nil), port)
	})

	t.Run("Valid", func(t *testing.T) {
		var priority int = 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		conn, err := NewConnection(maxPriority, linkBandwidth)
		require.NoError(t, err)

		outputPort, err := newOutputPort(conn, priority)
		require.NoError(t, err)

		assert.Equal(t, conn, outputPort.conn)
	})

	t.Run("NilConnection", func(t *testing.T) {
		_, err := newOutputPort(nil, 1)
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestInputPortConnection(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1
	var linkBandwidth int = 1

	conn, err := NewConnection(maxPriority, linkBandwidth)
	require.NoError(t, err)
	buff, err := newBuffer(1, 1)
	require.NoError(t, err)

	port, err := newInputPort(conn, buff)
	require.NoError(t, err)
	assert.Equal(t, conn, port.connection())
}

func TestInputPortReadIntoBuffer(t *testing.T) {
	t.Parallel()

	t.Run("FlitInChannel", func(t *testing.T) {
		bufferSize := 3
		var maxPriority int = 1
		var linkBandwidth int = 3

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		flits := make([]packet.Flit, linkBandwidth)
		for i := 0; i < linkBandwidth; i++ {
			flits[i] = packet.NewHeaderFlit(fmt.Sprintf("t%d", i), uuid.New(), 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}})
			port.conn.flitChannel() <- flits[i]
		}

		err := port.readIntoBuffer(0)
		require.NoError(t, err)

		for i := 0; i < linkBandwidth; i++ {
			gotFlit, exists := port.buff.popFlit(flits[i].Priority())
			assert.True(t, exists)
			assert.Equal(t, flits[i], gotFlit)
		}
	})

	t.Run("NoFlitInChannel", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		err := port.readIntoBuffer(0)
		require.NoError(t, err)
	})

	t.Run("FullBuffer", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		flit := packet.NewHeaderFlit("t", uuid.New(), 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}})

		port.conn.flitChannel() <- flit
		err := port.readIntoBuffer(0)
		require.NoError(t, err)

		port.conn.flitChannel() <- flit
		err = port.readIntoBuffer(0)
		require.Error(t, err)
	})
}

func TestInputPortPeakBuffer(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		flit := packet.NewHeaderFlit("t", uuid.New(), 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}, domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}})
		port.buff.addFlit(flit)

		gotFlit, exists := port.peakBuffer(flit.Priority())
		assert.True(t, exists)
		assert.Equal(t, flit, gotFlit)
	})

	t.Run("EmptyBuffer", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		flit, exists := port.peakBuffer(1)
		assert.False(t, exists)
		assert.Nil(t, flit)
	})
}

func TestInputPortReadOutOfBuffer(t *testing.T) {
	t.Parallel()

	t.Run("FlitInBuffer", func(t *testing.T) {
		bufferSize := 3
		var maxPriority int = 1
		var linkBandwidth int = 3

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		flits := make([]packet.Flit, linkBandwidth)
		for i := 0; i < linkBandwidth; i++ {
			flits[i] = packet.NewHeaderFlit("t", uuid.New(), 0, 1, 100, domain.Route{domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}, domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}})
			port.buff.addFlit(flits[i])
		}

		for i := 0; i < linkBandwidth; i++ {
			gotFlit, exists := port.readOutOfBuffer(flits[i].Priority())
			assert.True(t, exists)
			assert.Equal(t, flits[i], gotFlit)
			assert.Equal(t, 1, <-port.conn.creditChannel(flits[i].Priority()))
		}
	})

	t.Run("NoFlitInBuffer", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1
		var linkBandwidth int = 1

		port := testInputPort(t, bufferSize, maxPriority, linkBandwidth)

		gotFlit, exists := port.readOutOfBuffer(1)
		assert.False(t, exists)
		assert.Nil(t, gotFlit)
	})
}

func TestOutputPortConnection(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1
	var linkBandwidth int = 1

	conn, err := NewConnection(maxPriority, linkBandwidth)
	require.NoError(t, err)
	port, err := newOutputPort(conn, 1)
	require.NoError(t, err)

	assert.Equal(t, conn, port.connection())
}

func TestOutputPortAllowedToSend(t *testing.T) {
	t.Parallel()

	t.Run("Allowed", func(t *testing.T) {
		var priority int = 1
		var linkBandwidth int = 3

		port := testOutputPort(t, priority, linkBandwidth)

		port.credits[priority] = linkBandwidth

		for i := 0; i < linkBandwidth; i++ {
			assert.True(t, port.allowedToSend(priority))
			port.credits[priority]--
			port.conn.flitChannel() <- packet.NewTailFlit("t", uuid.New(), 2, priority)
		}
	})

	t.Run("NotAllowedLackingCredits", func(t *testing.T) {
		var priority int = 1
		var linkBandwidth int = 1

		port := testOutputPort(t, priority, linkBandwidth)

		port.credits[priority] = 0

		assert.False(t, port.allowedToSend(priority))
	})

	t.Run("NotAllowedReachedBandwidth", func(t *testing.T) {
		var priority int = 1
		var linkBandwidth int = 1

		port := testOutputPort(t, priority, linkBandwidth)

		port.credits[priority] = linkBandwidth + 1

		for i := 0; i < linkBandwidth; i++ {
			port.conn.flitChannel() <- packet.NewTailFlit("t", uuid.New(), 2, priority)
		}

		assert.False(t, port.allowedToSend(priority))
	})
}

func TestOutputPortSendFlit(t *testing.T) {
	t.Parallel()

	t.Run("AllowedToSend", func(t *testing.T) {
		var priority int = 1
		var linkBandwidth int = 3

		port := testOutputPort(t, priority, linkBandwidth)

		for i := 0; i < priority; i++ {
			port.credits[i] = linkBandwidth
		}

		for i := 0; i < priority; i++ {
			flit := packet.NewTailFlit("t", uuid.New(), 2, i)

			err := port.sendFlit(0, flit)
			require.NoError(t, err)
			assert.Equal(t, linkBandwidth-1, port.credits[i])
			assert.Equal(t, flit, <-port.conn.flitChannel())
		}
	})

	t.Run("NotAllowedToSend", func(t *testing.T) {
		var priority int = 1
		var linkBandwidth int = 1

		port := testOutputPort(t, priority, linkBandwidth)

		port.credits[priority] = 0

		err := port.sendFlit(0, packet.NewTailFlit("t", uuid.New(), 2, priority))
		require.ErrorIs(t, err, domain.ErrPortNoCredit)
		assert.Empty(t, port.conn.flitChannel())
	})
}

func TestOutputPortUpdateCredits(t *testing.T) {
	t.Parallel()

	t.Run("PendingCredits", func(t *testing.T) {
		var credits int = 3
		var priority int = 1
		var linkBandwidth int = 3

		port := testOutputPort(t, priority, linkBandwidth)

		for i := 0; i < credits; i++ {
			port.conn.creditChannel(priority) <- 1
			port.updateCredits()
		}

		assert.Equal(t, credits, port.credits[priority])
	})
}
