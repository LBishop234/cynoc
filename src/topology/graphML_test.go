package topology

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaricom/goGraphML/graphml"
)

func testGraphmlGraph(t *testing.T) (*graphml.Graph, []*Node, []*Edge) {
	nA, err := NewNode("nA")
	require.NoError(t, err)

	nB, err := NewNode("nB")
	require.NoError(t, err)

	e, err := NewEdge("e", nA.NodeID(), nB.NodeID())
	require.NoError(t, err)

	graphmlStr := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<graphml xmlns="http://graphml.graphdrawing.org/xmlns"  
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
	<graph id="G" edgedefault="undirected">
		<node id="%s">
		</node>
		<node id="%s">
		</node>
		<edge id="%s" source="%s" target="%s">
		</edge>
	</graph>
	</graphml>`,
		nA.NodeID(),
		nB.NodeID(),
		e.ID(),
		nA.NodeID(),
		nB.NodeID(),
	)

	gml := graphml.NewGraphML("topology")
	err = gml.Decode(bytes.NewReader([]byte(graphmlStr)))
	require.NoError(t, err)

	return gml.Graphs[0], []*Node{nA, nB}, []*Edge{e}
}

func TestGraphMLNodes(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		graph, nodes, _ := testGraphmlGraph(t)

		gotNodes, err := graphMLNodes(graph.Nodes)
		require.NoError(t, err)

		assert.Len(t, gotNodes, len(nodes))
	})
}

func TestGraphMLEdges(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		graph, nodes, edges := testGraphmlGraph(t)

		nodeMap := make(map[string]*Node, len(nodes))
		for i := 0; i < len(nodes); i++ {
			nodeMap[nodes[i].NodeID()] = nodes[i]
		}

		gotEdges, err := graphMLEdges(nodeMap, graph.Edges)
		require.NoError(t, err)

		assert.Len(t, gotEdges, len(edges))
	})
}

func TestGraphMLToEdge(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		graph, nodes, edges := testGraphmlGraph(t)

		nodeMap := make(map[string]*Node, len(nodes))
		for i := 0; i < len(nodes); i++ {
			nodeMap[nodes[i].NodeID()] = nodes[i]
		}

		edge, err := parseGraphMLEdge(nodeMap, graph.Edges[0])
		require.NoError(t, err)
		assert.Equal(t, edges[0], edge)
	})
}
