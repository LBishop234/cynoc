package packet

import (
	"testing"

	"main/src/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReconstructor(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		reconstructor := NewReconstructor()
		assert.Implements(t, (*Reconstructor)(nil), reconstructor)
	})

	t.Run("Valid", func(t *testing.T) {
		reconstructor := NewReconstructor()
		assert.NotNil(t, reconstructor)
		assert.Nil(t, reconstructor.headerFlit)
		assert.Empty(t, reconstructor.bodyFlits)
		assert.Nil(t, reconstructor.tailFlit)
	})
}

func TestReconstructorSetHeader(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		reconstructor := NewReconstructor()

		headerFlit := NewHeaderFlit("t", "AA", 0, 1, 100, route)

		err := reconstructor.SetHeader(headerFlit)
		require.NoError(t, err)
		assert.Equal(t, headerFlit, reconstructor.headerFlit)
	})

	t.Run("NilHeaderFlit", func(t *testing.T) {
		reconstructor := NewReconstructor()

		err := reconstructor.SetHeader(nil)
		require.ErrorIs(t, domain.ErrNilParameter, err)
	})

	t.Run("InvalidAlreadySet", func(t *testing.T) {
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		reconstructor := NewReconstructor()

		headerFlit := NewHeaderFlit("t", "AA", 0, 1, 100, route)

		err := reconstructor.SetHeader(headerFlit)
		require.NoError(t, err)

		err = reconstructor.SetHeader(headerFlit)
		require.ErrorIs(t, domain.ErrFlitAlreadySet, err)
	})
}

func TestReconstructorAddBody(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		reconstructor := NewReconstructor()
		bodyFlit := NewBodyFlit("t", "AA", 1, 2, 1)

		err := reconstructor.AddBody(bodyFlit)
		require.NoError(t, err)
		assert.Equal(t, []BodyFlit{bodyFlit}, reconstructor.bodyFlits)
	})

	t.Run("NilHeaderFlit", func(t *testing.T) {
		reconstructor := NewReconstructor()

		err := reconstructor.AddBody(nil)
		require.ErrorIs(t, domain.ErrNilParameter, err)
	})
}

func TestReconstructorSetTail(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		reconstructor := NewReconstructor()
		tailFlit := NewTailFlit("t", "AA", 2, 1)

		err := reconstructor.SetTail(tailFlit)
		require.NoError(t, err)
		assert.Equal(t, tailFlit, reconstructor.tailFlit)
	})

	t.Run("NilHeaderFlit", func(t *testing.T) {
		reconstructor := NewReconstructor()

		err := reconstructor.SetTail(nil)
		require.ErrorIs(t, domain.ErrNilParameter, err)
	})

	t.Run("InvalidAlreadySet", func(t *testing.T) {
		reconstructor := NewReconstructor()
		tailFlit := NewTailFlit("t", "AA", 2, 1)

		err := reconstructor.SetTail(tailFlit)
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
		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}
		var bodySize int = 4

		packet := NewPacket(trafficFlowID, packetID, priority, deadline, route, bodySize)
		flits := packet.Flits(1)

		reconstructor := NewReconstructor()

		headerFlit, ok := flits[0].(HeaderFlit)
		require.True(t, ok)
		err := reconstructor.SetHeader(headerFlit)
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

	t.Run("HeaderUnset", func(t *testing.T) {
		reconstructor := NewReconstructor()

		reconstructor.SetTail(NewTailFlit("t", "AA", 2, 1))

		pkt, err := reconstructor.Reconstruct()
		require.ErrorIs(t, domain.ErrFlitUnset, err)
		assert.Nil(t, pkt)
	})

	t.Run("TailUnset", func(t *testing.T) {
		reconstructor := NewReconstructor()

		src := domain.NodeID{ID: "n1", Pos: domain.NewPosition(0, 0)}
		dst := domain.NodeID{ID: "n2", Pos: domain.NewPosition(0, 1)}
		route := domain.Route{src, dst}

		reconstructor.SetHeader(NewHeaderFlit("t", "AA", 0, 1, 100, route))

		pkt, err := reconstructor.Reconstruct()
		require.ErrorIs(t, domain.ErrFlitUnset, err)
		assert.Nil(t, pkt)
	})
}
