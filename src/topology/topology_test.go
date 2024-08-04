package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTop(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		nodes := map[string]*Node{
			"1": NewNode("1"),
			"2": NewNode("2"),
		}
		edges := map[string]*Edge{
			"1": {
				id: "1",
				a:  nodes["1"].NodeID(),
				b:  nodes["2"].NodeID(),
			},
		}

		topology := NewTopology(nodes, edges)
		assert.Equal(t, nodes, topology.nodes)
		assert.Equal(t, edges, topology.edges)
	})
}

func TestTopologyNodes(t *testing.T) {
	t.Parallel()

	nodes := map[string]*Node{
		"1": NewNode("1"),
		"2": NewNode("2"),
	}

	topology := NewTopology(nodes, nil)
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

	topology := NewTopology(nil, edges)
	assert.Equal(t, edges, topology.Edges())
}

func TestNewNode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var id string = "1"

		node := NewNode(id)
		assert.Equal(t, id, node.NodeID())
	})
}

func TestNodeID(t *testing.T) {
	t.Parallel()

	node := NewNode("1")
	assert.Equal(t, string("1"), node.NodeID())
}

func TestNewEdge(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		var id string = "1"
		var source string = "n1"
		var target string = "n2"

		edge := NewEdge(id, source, target)
		assert.Equal(t, id, edge.id)
		assert.Equal(t, source, edge.a)
		assert.Equal(t, target, edge.b)
	})
}

func TestEdgeID(t *testing.T) {
	t.Parallel()

	edge := NewEdge("1", "n1", "n2")
	assert.Equal(t, string("1"), edge.ID())
}

func TestEdgeSource(t *testing.T) {
	t.Parallel()

	source := "n1"
	edge := NewEdge("", source, "n2")
	assert.Equal(t, source, edge.A())
}

func TestEdgeTarget(t *testing.T) {
	t.Parallel()

	target := "n2"
	edge := NewEdge("", "n1", target)
	assert.Equal(t, target, edge.B())
}
