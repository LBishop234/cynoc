package topology

import (
	"testing"

	"main/src/domain"

	"github.com/stretchr/testify/require"
)

type nodeSpec struct {
	id string
	x  int
	y  int
}

type edgeSpec struct {
	id  string
	src string
	dst string
}

func ThreeHorozontalLine(t testing.TB) *Topology {
	nodeSpecs := []nodeSpec{
		{"n0", 0, 0},
		{"n1", 1, 0},
		{"n2", 2, 0},
	}

	edgeSpecs := []edgeSpec{
		{"e0", "n0", "n1"},
		{"e1", "n1", "n2"},
	}

	return constructTopology(t, nodeSpecs, edgeSpecs)
}

func ThreeByThreeMesh(t testing.TB) *Topology {
	nodeSpec := []nodeSpec{
		{"n0", 0, 0},
		{"n1", 1, 0},
		{"n2", 2, 0},
		{"n3", 0, 1},
		{"n4", 1, 1},
		{"n5", 2, 1},
		{"n6", 0, 2},
		{"n7", 1, 2},
		{"n8", 2, 2},
	}

	edgeSpec := []edgeSpec{
		{"e0", "n0", "n1"},
		{"e1", "n0", "n3"},
		{"e2", "n1", "n2"},
		{"e3", "n1", "n4"},
		{"e4", "n2", "n5"},
		{"e5", "n3", "n4"},
		{"e6", "n3", "n6"},
		{"e7", "n4", "n5"},
		{"e8", "n4", "n7"},
		{"e9", "n5", "n8"},
		{"e10", "n6", "n7"},
		{"e11", "n7", "n8"},
	}

	return constructTopology(t, nodeSpec, edgeSpec)
}

func FourByFourMesh(t testing.TB) *Topology {
	nodeSpec := []nodeSpec{
		{"n0", 0, 0},
		{"n1", 1, 0},
		{"n2", 2, 0},
		{"n3", 3, 0},
		{"n4", 0, 1},
		{"n5", 1, 1},
		{"n6", 2, 1},
		{"n7", 3, 1},
		{"n8", 0, 2},
		{"n9", 1, 2},
		{"n10", 2, 2},
		{"n11", 3, 2},
		{"n12", 0, 3},
		{"n13", 1, 3},
		{"n14", 2, 3},
		{"n15", 3, 3},
	}

	edgeSpec := []edgeSpec{
		{"e0", "n0", "n1"},
		{"e1", "n0", "n4"},
		{"e2", "n1", "n2"},
		{"e3", "n1", "n5"},
		{"e4", "n2", "n3"},
		{"e5", "n2", "n6"},
		{"e6", "n3", "n7"},
		{"e7", "n4", "n5"},
		{"e8", "n4", "n8"},
		{"e9", "n5", "n6"},
		{"e10", "n5", "n9"},
		{"e11", "n6", "n7"},
		{"e12", "n6", "n10"},
		{"e13", "n7", "n11"},
		{"e14", "n8", "n9"},
		{"e15", "n8", "n12"},
		{"e16", "n9", "n10"},
		{"e17", "n9", "n13"},
		{"e18", "n10", "n11"},
		{"e19", "n10", "n14"},
		{"e20", "n11", "n15"},
		{"e21", "n12", "n13"},
		{"e22", "n13", "n14"},
		{"e23", "n14", "n15"},
	}

	return constructTopology(t, nodeSpec, edgeSpec)
}

func constructTopology(tb testing.TB, nodeSpecs []nodeSpec, edgeSpecs []edgeSpec) *Topology {
	nodes := make(map[string]*Node, len(nodeSpecs))
	for i := 0; i < len(nodeSpecs); i++ {
		aNode, err := NewNode(nodeSpecs[i].id, domain.NewPosition(nodeSpecs[i].x, nodeSpecs[i].y))
		require.NoError(tb, err)

		nodes[aNode.NodeID().ID] = aNode
	}

	edges := make(map[string]*Edge, len(edgeSpecs))
	for i := 0; i < len(edgeSpecs); i++ {
		aEdge, err := NewEdge(edgeSpecs[i].id, edgeSpecs[i].src, edgeSpecs[i].dst)
		require.NoError(tb, err)

		edges[aEdge.ID()] = aEdge
	}

	top, err := NewTopology(nodes, edges)
	require.NoError(tb, err)

	return top
}
