package components

import (
	"io"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	t.Run("ImplementsInterface", func(t *testing.T) {
		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)
		assert.Implements(t, (*Connection)(nil), conn)
	})

	t.Run("Valid", func(t *testing.T) {
		conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)
		assert.NotNil(t, conn.flitChan)
		assert.NotNil(t, conn.creditChan)
	})
}

func TestConnectionFlitChannel(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)
	assert.Equal(t, conn.flitChan, conn.flitChannel())
}

func TestConnectionCreditChannel(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1
	var priority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	conn.creditChan[priority] = make(chan int, 1)

	assert.Equal(t, conn.creditChan[priority], conn.creditChannel(priority))
}

func TestConnectionGetDestRouterID(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	nodeID := domain.NodeID{ID: "n"}
	conn.destRouter = nodeID
	assert.Equal(t, nodeID, conn.GetDstRouter())
}

func TestConnectionGetSrcRouterID(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	nodeID := domain.NodeID{ID: "n"}
	conn.srcRouter = nodeID
	assert.Equal(t, nodeID, conn.GetSrcRouter())
}

func TestConnectionSetDestRouterID(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	nodeID := domain.NodeID{ID: "n"}
	conn.SetDstRouter(nodeID)
	assert.Equal(t, nodeID, conn.destRouter)
}

func TestConnectionSetSrcRouterID(t *testing.T) {
	t.Parallel()

	var maxPriority int = 1

	conn, err := NewConnection(maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	nodeID := domain.NodeID{ID: "n"}
	conn.SetSrcRouter(nodeID)
	assert.Equal(t, nodeID, conn.srcRouter)
}
