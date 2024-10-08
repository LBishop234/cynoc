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

func testInputPort(t *testing.T, bufferSize, maxPriority int) *inputPortImpl {
	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)
	buff, err := newBuffer(bufferSize, 1, zerolog.New(io.Discard))
	require.NoError(t, err)

	port, err := newInputPort(conn, buff, zerolog.New(io.Discard))
	require.NoError(t, err)

	for _, credChan := range conn.creditChannels() {
		<-credChan
	}

	return port
}

func testOutputPort(t *testing.T, maxPriority int) *outputPortImpl {
	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	port, err := newOutputPort(conn, maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	return port
}

func TestNewInputPort(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		port := testInputPort(t, 1, 1)

		assert.Implements(t, (*inputPort)(nil), port)
	})

	t.Run("Valid", func(t *testing.T) {
		var bufferSize int = 32
		var maxPriority int = 4

		buff, err := newBuffer(bufferSize, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		port, err := newInputPort(conn, buff, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.Equal(t, conn, port.conn)
		assert.Equal(t, buff, port.buff)
	})

	t.Run("NilConnection", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		_, err = newInputPort(nil, buff, zerolog.New(io.Discard))
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})

	t.Run("NilBuffer", func(t *testing.T) {
		var maxPriority int = 1

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		_, err = newInputPort(conn, nil, zerolog.New(io.Discard))
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestNewOutputPort(t *testing.T) {
	t.Parallel()

	t.Run("ValidInterface", func(t *testing.T) {
		port := testOutputPort(t, 1)

		assert.Implements(t, (*outputPort)(nil), port)
	})

	t.Run("Valid", func(t *testing.T) {
		var priority int = 1
		var maxPriority int = 1

		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		outputPort, err := newOutputPort(conn, priority, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.Equal(t, conn, outputPort.conn)
	})

	t.Run("NilConnection", func(t *testing.T) {
		_, err := newOutputPort(nil, 1, zerolog.New(io.Discard))
		require.ErrorIs(t, err, domain.ErrNilParameter)
	})
}

func TestInputPortConnection(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)
	buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
	require.NoError(t, err)

	port, err := newInputPort(conn, buff, zerolog.New(io.Discard))
	require.NoError(t, err)
	assert.Equal(t, conn, port.connection())
}

func TestInputPortReadIntoBuffer(t *testing.T) {
	t.Parallel()

	t.Run("FlitInChannel", func(t *testing.T) {
		var packetID string = "AA"
		var bufferSize int = 3
		var maxPriority int = 1
		var flitPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		flit := packet.NewHeaderFlit("t0", packetID, 0, flitPriority, 100, domain.Route{"n1", "n2"}, zerolog.New(io.Discard))
		port.conn.flitChannel() <- flit

		err := port.readIntoBuffer(0)
		require.NoError(t, err)

		gotFlit, exists := port.buff.popFlit(flitPriority)
		assert.True(t, exists)
		assert.Equal(t, flit, gotFlit)
	})

	t.Run("NoFlitInChannel", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		err := port.readIntoBuffer(0)
		require.NoError(t, err)
	})

	t.Run("FullBuffer", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		flit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{"n1", "n2"}, zerolog.New(io.Discard))

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

		port := testInputPort(t, bufferSize, maxPriority)

		flit := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{"n1", "n2", "n1", "n2"}, zerolog.New(io.Discard))
		port.buff.addFlit(flit)

		gotFlit, exists := port.peakBuffer(flit.Priority())
		assert.True(t, exists)
		assert.Equal(t, flit, gotFlit)
	})

	t.Run("EmptyBuffer", func(t *testing.T) {
		bufferSize := 1
		var maxPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		flit, exists := port.peakBuffer(1)
		assert.False(t, exists)
		assert.Nil(t, flit)
	})
}

func TestInputPortReadOutOfBuffer(t *testing.T) {
	t.Parallel()

	t.Run("FlitInBuffer", func(t *testing.T) {
		var bufferSize int = 3
		var maxPriority int = 1
		var flitPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		flit := packet.NewHeaderFlit("t", "AA", 0, flitPriority, 100, domain.Route{"n1", "n2"}, zerolog.New(io.Discard))
		port.buff.addFlit(flit)

		gotFlit, exists := port.readOutOfBuffer(0, flitPriority)
		assert.True(t, exists)
		assert.Equal(t, flit, gotFlit)
		assert.Equal(t, 1, <-port.conn.creditChannel(flit.Priority()))
	})

	t.Run("NoFlitInBuffer", func(t *testing.T) {
		var bufferSize int = 1
		var maxPriority int = 1

		port := testInputPort(t, bufferSize, maxPriority)

		gotFlit, exists := port.readOutOfBuffer(0, 1)
		assert.False(t, exists)
		assert.Nil(t, gotFlit)
	})
}

func TestOutputPortConnection(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)
	port, err := newOutputPort(conn, 1, zerolog.New(io.Discard))
	require.NoError(t, err)

	assert.Equal(t, conn, port.connection())
}

func TestOutputPortAllowedToSend(t *testing.T) {
	t.Parallel()

	t.Run("Allowed", func(t *testing.T) {
		var priority int = 1

		port := testOutputPort(t, priority)

		port.credits[priority] = 1

		assert.True(t, port.allowedToSend(priority))
		port.credits[priority]--
		port.conn.flitChannel() <- packet.NewTailFlit("t", "AA", 2, priority, zerolog.New(io.Discard))
	})

	t.Run("NotAllowedLackingCredits", func(t *testing.T) {
		var priority int = 1

		port := testOutputPort(t, priority)

		port.credits[priority] = 0

		assert.False(t, port.allowedToSend(priority))
	})
}

func TestOutputPortSendFlit(t *testing.T) {
	t.Parallel()

	t.Run("AllowedToSend", func(t *testing.T) {
		var priority int = 1

		port := testOutputPort(t, priority)

		for i := 0; i < priority; i++ {
			port.credits[i] = 1
		}

		for i := 0; i < priority; i++ {
			flit := packet.NewTailFlit("t", "AA", 2, i, zerolog.New(io.Discard))

			err := port.sendFlit(0, flit)
			require.NoError(t, err)
			assert.Equal(t, 0, port.credits[i])
			assert.Equal(t, flit, <-port.conn.flitChannel())
		}
	})

	t.Run("NotAllowedToSend", func(t *testing.T) {
		var priority int = 1

		port := testOutputPort(t, priority)

		port.credits[priority] = 0

		err := port.sendFlit(0, packet.NewTailFlit("t", "AA", 2, priority, zerolog.New(io.Discard)))
		require.ErrorIs(t, err, domain.ErrPortNoCredit)
		assert.Empty(t, port.conn.flitChannel())
	})
}

func TestOutputPortUpdateCredits(t *testing.T) {
	t.Parallel()

	t.Run("PendingCredits", func(t *testing.T) {
		var credits int = 3
		var priority int = 1

		port := testOutputPort(t, priority)

		for i := 0; i < credits; i++ {
			port.conn.creditChannel(priority) <- 1
			port.updateCredits()
		}

		assert.Equal(t, credits, port.credits[priority])
	})
}
