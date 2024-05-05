package topology

import (
	"testing"

	"main/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTop(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		nodes := map[string]*Node{
			"1": {
				nodeID: domain.NodeID{
					ID:  "1",
					Pos: domain.NewPosition(0, 0),
				},
			},
			"2": {
				nodeID: domain.NodeID{
					ID:  "2",
					Pos: domain.NewPosition(0, 1),
				},
			},
		}
		edges := map[string]*Edge{
			"1": {
				id: "1",
				a:  nodes["1"].NodeID().ID,
				b:  nodes["2"].NodeID().ID,
			},
		}

		topology, err := NewTopology(nodes, edges)
		require.NoError(t, err)
		assert.Equal(t, nodes, topology.nodes)
		assert.Equal(t, edges, topology.edges)
	})
}

func TestTopologyNodes(t *testing.T) {
	t.Parallel()

	nodes := map[string]*Node{
		"1": {
			nodeID: domain.NodeID{
				ID:  "1",
				Pos: domain.NewPosition(0, 0),
			},
		},
		"2": {
			nodeID: domain.NodeID{
				ID:  "2",
				Pos: domain.NewPosition(0, 1),
			},
		},
	}

	topology, err := NewTopology(nodes, nil)
	require.NoError(t, err)
	assert.Equal(t, nodes, topology.Nodes())
}

func TestTopologyEdges(t *testing.T) {
	t.Parallel()

	edges := map[string]*Edge{
		"1": {
			id: "1",
			a:  "n1",
			b:  "n2",
		},
	}

	topology, err := NewTopology(nil, edges)
	require.NoError(t, err)
	assert.Equal(t, edges, topology.Edges())
}

func TestNewNode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var id string = "1"
		var pos domain.Position = domain.NewPosition(0, 0)

		node, err := NewNode(id, pos)
		require.NoError(t, err)
		assert.Equal(t, id, node.NodeID().ID)
		assert.Equal(t, pos, node.NodeID().Pos)
	})
}

func TestNodestring(t *testing.T) {
	t.Parallel()

	node, err := NewNode("1", domain.Position{})
	require.NoError(t, err)
	assert.Equal(t, string("1"), node.NodeID().ID)
}

func TestNodePosition(t *testing.T) {
	t.Parallel()

	pos := domain.NewPosition(0, 0)
	node, err := NewNode("", pos)
	require.NoError(t, err)
	assert.Equal(t, pos, node.NodeID().Pos)
}

func TestNewEdge(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var id string = "1"
		var source string = "n1"
		var target string = "n2"

		edge, err := NewEdge(id, source, target)
		require.NoError(t, err)
		assert.Equal(t, id, edge.id)
		assert.Equal(t, source, edge.a)
		assert.Equal(t, target, edge.b)
	})
}

func TestEdgestring(t *testing.T) {
	t.Parallel()

	edge, err := NewEdge("1", "n1", "n2")
	require.NoError(t, err)
	assert.Equal(t, string("1"), edge.ID())
}

func TestEdgeSource(t *testing.T) {
	t.Parallel()

	source := "n1"
	edge, err := NewEdge("", source, "n2")
	require.NoError(t, err)
	assert.Equal(t, source, edge.A())
}

func TestEdgeTarget(t *testing.T) {
	t.Parallel()

	target := "n2"
	edge, err := NewEdge("", "n1", target)
	require.NoError(t, err)
	assert.Equal(t, target, edge.B())
}
