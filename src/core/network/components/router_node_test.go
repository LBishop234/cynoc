package components

import (
	"testing"

	"main/src/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouterNode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		conf := RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n1",
				Pos: domain.NewPosition(1, 1),
			},
			SimConfig: domain.SimConfig{
				RoutingAlgorithm: domain.XYRouting,
				BufferSize:       16,
				FlitSize:         8,
				MaxPriority:      4,
				ProcessingDelay:  5,
				LinkBandwidth:    1,
			},
		}

		routerNode, err := NewRouterNode(conf)
		require.NoError(t, err)
		assert.Equal(t, conf.NodeID, routerNode.NodeID())
		assert.NotNil(t, routerNode.Router)
		assert.NotNil(t, routerNode.NetworkInterface)
	})

	t.Run("NewRouterError", func(t *testing.T) {
		conf := RouterConfig{
			NodeID: domain.NodeID{
				ID:  "n1",
				Pos: domain.NewPosition(1, 1),
			},
			SimConfig: domain.SimConfig{
				RoutingAlgorithm: domain.XYRouting,
				BufferSize:       -1,
				FlitSize:         8,
				MaxPriority:      3,
				ProcessingDelay:  5,
				LinkBandwidth:    1,
			},
		}

		_, err := NewRouterNode(conf)
		require.Error(t, err)
	})

	t.Run("NewNetworkInterfaceError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("SetNetworkInterfaceError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})
}
