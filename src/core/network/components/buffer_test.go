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

func TestNewBuffer(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.Implements(t, (*buffer)(nil), buff)
	})

	t.Run("Valid", func(t *testing.T) {
		var capacity int = 1
		var maxPriority int = 1

		buff, err := newBuffer(capacity, maxPriority, zerolog.New(io.Discard))
		require.NoError(t, err)

		assert.Equal(t, capacity, buff.bufferCap)
		assert.Equal(t, capacity/maxPriority, buff.vChanCap)
		assert.NotNil(t, buff.flits)
	})

	t.Run("InvalidCapacity", func(t *testing.T) {
		var capacity int = 0

		_, err := newBuffer(capacity, 1, zerolog.New(io.Discard))
		require.Error(t, err)
	})
}

func TestBufferTotalCapacity(t *testing.T) {
	t.Parallel()

	var capacity int = 1
	var maxPriority int = 1

	buff, err := newBuffer(capacity, maxPriority, zerolog.New(io.Discard))
	require.NoError(t, err)

	assert.Equal(t, capacity, buff.totalCapacity())
}

func TestBufferPopFlit(t *testing.T) {
	t.Parallel()

	t.Run("ValidEmpty", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		flit, exists := buff.popFlit(0)
		assert.Nil(t, flit)
		assert.False(t, exists)
	})

	t.Run("ValidFlit", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		flit := packet.NewTailFlit("t", "AA", 2, 1, zerolog.New(io.Discard))
		buff.flits[flit.Priority()] = append(buff.flits[flit.Priority()], flit)

		gotFlit, exists := buff.popFlit(flit.Priority())
		assert.Equal(t, flit, gotFlit)
		assert.True(t, exists)
		assert.NotContains(t, buff.flits, flit)
	})

	t.Run("ValidMultipleFlits", func(t *testing.T) {
		buff, err := newBuffer(2, 2, zerolog.New(io.Discard))
		require.NoError(t, err)

		flit1 := packet.NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{domain.NodeID{ID: "n1"}, domain.NodeID{ID: "n2"}}, zerolog.New(io.Discard))
		buff.flits[flit1.Priority()] = append(buff.flits[flit1.Priority()], flit1)

		flit2 := packet.NewTailFlit("t", "AA", 1, 1, zerolog.New(io.Discard))
		buff.flits[flit2.Priority()] = append(buff.flits[flit2.Priority()], flit2)

		gotFlit1, exists := buff.popFlit(flit1.Priority())
		assert.Equal(t, flit1, gotFlit1)
		assert.True(t, exists)
		assert.NotContains(t, buff.flits, flit1)

		gotFlit2, exists := buff.popFlit(flit2.Priority())
		assert.Equal(t, flit2, gotFlit2)
		assert.True(t, exists)
		assert.NotContains(t, buff.flits, flit2)
	})
}

func TestBufferAddFlit(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		flit := packet.NewTailFlit("t", "AA", 2, 1, zerolog.New(io.Discard))
		err = buff.addFlit(flit)
		require.NoError(t, err)
		assert.Contains(t, buff.flits[flit.Priority()], flit)
	})

	t.Run("NoCapacity", func(t *testing.T) {
		buff, err := newBuffer(1, 1, zerolog.New(io.Discard))
		require.NoError(t, err)

		flit := packet.NewTailFlit("t", "AA", 2, 1, zerolog.New(io.Discard))
		err = buff.addFlit(flit)
		require.NoError(t, err)

		err = buff.addFlit(flit)
		require.ErrorIs(t, err, domain.ErrBufferNoCapacity)
	})
}

func TestValidBufferSize(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var bufferSize int = 1
		var maxPriority int = 1
		err := validBufferSize(bufferSize, maxPriority)
		require.NoError(t, err)
	})

	t.Run("InvalidCapacity", func(t *testing.T) {
		var bufferSize int = 0
		var maxPriority int = 1
		err := validBufferSize(bufferSize, maxPriority)
		require.ErrorIs(t, err, domain.ErrInvalidParameter)
	})

	t.Run("InvalidMaxPriority", func(t *testing.T) {
		var bufferSize int = 3
		var maxPriority int = 2
		err := validBufferSize(bufferSize, maxPriority)
		require.ErrorIs(t, err, domain.ErrInvalidParameter)
	})
}

func TestBufferVChanCapacity(t *testing.T) {
	t.Parallel()

	type testCase struct {
		cap         int
		maxPrio     int
		expectedVal int
		err         error
	}

	testCases := []testCase{
		{1, 1, 1, nil},
		{2, 1, 2, nil},
		{2, 2, 1, nil},
		{3, 2, 0, domain.ErrInvalidParameter},
		{4, 3, 0, domain.ErrInvalidParameter},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			val, err := bufferVChanCapacity(tc.cap, tc.maxPrio)
			require.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.expectedVal, val)
		})
	}
}
