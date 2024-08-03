package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTop(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		nodes := map[string]string{
			"1": "1",
			"2": "2",
		}
		edges := map[string]*Edge{
			"1": {
				id: "1",
				a:  nodes["1"],
				b:  nodes["2"],
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

	nodes := map[string]string{
		"1": "1",
		"2": "2",
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

		node, err := NewNode(id)
		require.NoError(t, err)
		assert.Equal(t, id, node)
	})
}

func TestNodestring(t *testing.T) {
	t.Parallel()

	node, err := NewNode("1")
	require.NoError(t, err)
	assert.Equal(t, string("1"), node)
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
