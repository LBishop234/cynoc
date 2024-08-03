package components

import (
	"io"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouterNode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		conf := RouterConfig{
			NodeID: domain.NodeID{
				ID: "n1",
			},
			SimConfig: domain.SimConfig{
				BufferSize:      16,
				MaxPriority:     4,
				ProcessingDelay: 5,
			},
		}

		routerNode, err := NewRouterNode(conf, zerolog.New(io.Discard).With().Logger())
		require.NoError(t, err)
		assert.Equal(t, conf.NodeID, routerNode.NodeID())
		assert.NotNil(t, routerNode.Router)
		assert.NotNil(t, routerNode.NetworkInterface)
	})

	t.Run("NewRouterError", func(t *testing.T) {
		conf := RouterConfig{
			NodeID: domain.NodeID{
				ID: "n1",
			},
			SimConfig: domain.SimConfig{
				BufferSize:      -1,
				MaxPriority:     3,
				ProcessingDelay: 5,
			},
		}

		_, err := NewRouterNode(conf, zerolog.New(io.Discard).With().Logger())
		require.Error(t, err)
	})

	t.Run("NewNetworkInterfaceError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})

	t.Run("SetNetworkInterfaceError", func(t *testing.T) {
		t.Skip("Cannot currently test, possible error cases cannot be met")
	})
}
