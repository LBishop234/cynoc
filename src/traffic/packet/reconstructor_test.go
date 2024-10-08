package packet

import (
	"io"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dummyHeaderFlit = NewHeaderFlit("t", "AA", 0, 1, 100, domain.Route{}, zerolog.New(io.Discard))

func TestNewReconstructor(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)
		assert.Implements(t, (*Reconstructor)(nil), reconstructor)
	})

	t.Run("Valid", func(t *testing.T) {
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}

		headerFlit := NewHeaderFlit("t", "AA", 0, 1, 100, route, zerolog.New(io.Discard))

		reconstructor, err := NewReconstructor(headerFlit, zerolog.New(io.Discard))
		require.NoError(t, err)
		assert.NotNil(t, reconstructor)
		assert.Equal(t, headerFlit, reconstructor.headerFlit)
		assert.Empty(t, reconstructor.bodyFlits)
		assert.Nil(t, reconstructor.tailFlit)
	})

	t.Run("NilHeaderFlit", func(t *testing.T) {
		_, err := NewReconstructor(nil, zerolog.New(io.Discard))
		require.ErrorIs(t, domain.ErrNilParameter, err)
	})
}

func TestReconstructorAddBody(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		bodyFlit := NewBodyFlit("t", "AA", 1, 2, zerolog.New(io.Discard))

		err = reconstructor.AddBody(bodyFlit)
		require.NoError(t, err)
		assert.Equal(t, []BodyFlit{bodyFlit}, reconstructor.bodyFlits)
	})

	t.Run("NilBodyFlit", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		err = reconstructor.AddBody(nil)
		require.ErrorIs(t, domain.ErrNilParameter, err)
	})
}

func TestReconstructorSetTail(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		tailFlit := NewTailFlit("t", "AA", 2, 1, zerolog.New(io.Discard))

		err = reconstructor.SetTail(tailFlit)
		require.NoError(t, err)
		assert.Equal(t, tailFlit, reconstructor.tailFlit)
	})

	t.Run("InvalidAlreadySet", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		tailFlit := NewTailFlit("t", "AA", 2, 1, zerolog.New(io.Discard))

		err = reconstructor.SetTail(tailFlit)
		require.NoError(t, err)

		err = reconstructor.SetTail(tailFlit)
		require.ErrorIs(t, domain.ErrFlitAlreadySet, err)
	})
}

func TestReconstructorReconstruct(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var trafficFlowID string = "t"
		var packetID string = "AA"
		var priority int = 1
		var deadline int = 100
		src := "n1"
		dst := "n2"
		route := domain.Route{src, dst}
		var bodySize int = 4

		packet := NewPacket(trafficFlowID, packetID, priority, deadline, route, bodySize, zerolog.New(io.Discard))
		flits := packet.Flits()

		headerFlit, ok := flits[0].(HeaderFlit)
		require.True(t, ok)

		reconstructor, err := NewReconstructor(headerFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		for i := 1; i < len(flits)-1; i++ {
			bodyFlit, ok := flits[i].(BodyFlit)
			require.True(t, ok)
			err = reconstructor.AddBody(bodyFlit)
			require.NoError(t, err)
		}

		tailFlit, ok := flits[len(flits)-1].(TailFlit)
		require.True(t, ok)
		err = reconstructor.SetTail(tailFlit)
		require.NoError(t, err)

		gotPacket, err := reconstructor.Reconstruct()
		require.NoError(t, err)
		assert.Equal(t, packet, gotPacket)
	})

	t.Run("TailUnset", func(t *testing.T) {
		reconstructor, err := NewReconstructor(dummyHeaderFlit, zerolog.New(io.Discard))
		require.NoError(t, err)

		pkt, err := reconstructor.Reconstruct()
		require.ErrorIs(t, domain.ErrFlitUnset, err)
		assert.Nil(t, pkt)
	})
}
