package topology

import (
	"bytes"
	"fmt"
	"testing"

	"main/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaricom/goGraphML/graphml"
)

func testGraphmlGraph(t *testing.T) (*graphml.Graph, []*Node, []*Edge) {
	nA, err := NewNode("nA", domain.NewPosition(0, 0))
	require.NoError(t, err)

	nB, err := NewNode("nB", domain.NewPosition(1, 0))
	require.NoError(t, err)

	e, err := NewEdge("e", nA.NodeID().ID, nB.NodeID().ID)
	require.NoError(t, err)

	const weight float64 = 1

	graphmlStr := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<graphml xmlns="http://graphml.graphdrawing.org/xmlns"  
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
	<graph id="G" edgedefault="undirected">
		<node id="%s">
			<data key="x">%d</data>
			<data key="y">%d</data>
		</node>
		<node id="%s">
			<data key="x">%d</data>
			<data key="y">%d</data>
		</node>
		<edge id="%s" source="%s" target="%s">
			<data key="weight">%f</data>
		</edge>
	</graph>
	</graphml>`,
		nA.NodeID().ID,
		nA.NodeID().Pos.X(),
		nA.NodeID().Pos.Y(),
		nB.NodeID().ID,
		nB.NodeID().Pos.X(),
		nB.NodeID().Pos.Y(),
		e.ID(),
		nA.NodeID().ID,
		nB.NodeID().ID,
		weight,
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

func TestParseGraphMLNode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		graph, nodes, _ := testGraphmlGraph(t)

		node, err := parseGraphMLNode(graph.Nodes[0])
		require.NoError(t, err)
		assert.Equal(t, nodes[0], node)
	})
}

func TestGraphMLEdges(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		graph, nodes, edges := testGraphmlGraph(t)

		nodeMap := make(map[string]*Node, len(nodes))
		for i := 0; i < len(nodes); i++ {
			nodeMap[nodes[i].NodeID().ID] = nodes[i]
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
			nodeMap[nodes[i].NodeID().ID] = nodes[i]
		}

		edge, err := parseGraphMLEdge(nodeMap, graph.Edges[0])
		require.NoError(t, err)
		assert.Equal(t, edges[0], edge)
	})
}

func TestGraphMLDataToAttributes(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		gmlData := []*graphml.Data{
			{
				Key:   "x",
				Value: "1",
			},
			{
				Key:   "y",
				Value: "2",
			},
		}

		attributes := map[string]interface{}{
			"x": "1",
			"y": "2",
		}

		gotAttributes := graphMLDataToAttributes(gmlData)
		assert.Equal(t, attributes, gotAttributes)
	})
}
