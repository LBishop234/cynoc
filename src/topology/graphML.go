package topology

import (
	"os"

	"main/log"
	"main/src/domain"

	"github.com/yaricom/goGraphML/graphml"
)

func graphML(filepath string) (*Topology, error) {
	log.Log.Debug().Msg("reading GraphML topology file")

	f, err := os.Open(filepath)
	if err != nil {
		log.Log.Error().Err(err).Str("path", filepath).Msg("error opening GraphML topology file")
		return nil, err
	}
	log.Log.Debug().Msg("opened GraphML topology file")

	gml := graphml.NewGraphML("topology")
	if err = gml.Decode(f); err != nil {
		log.Log.Error().Err(err).Str("path", filepath).Msg("error parsing GraphML topology file")
		return nil, err
	}
	log.Log.Debug().Msg("decoded GraphML topology file")

	if len(gml.Graphs) != 1 {
		log.Log.Error().Err(domain.ErrInvalidTopology).Str("path", filepath).Msg("invalid GraphML topology file")
		return nil, domain.ErrInvalidTopology
	}
	graph := gml.Graphs[0]

	log.Log.Debug().Msg("parsing GraphML topology file")

	nodes, err := graphMLNodes(graph.Nodes)
	if err != nil {
		log.Log.Error().Err(err).Str("path", filepath).Msg("error parsing GraphML nodes")
		return nil, err
	}

	edges, err := graphMLEdges(nodes, graph.Edges)
	if err != nil {
		log.Log.Error().Err(err).Str("path", filepath).Msg("error parsing GraphML edges")
		return nil, err
	}

	log.Log.Debug().Msg("parsed GraphML topology file")

	top := NewTopology(nodes, edges)

	log.Log.Debug().Msg("loaded topology from GraphML file")
	return top, nil
}

func graphMLNodes(gmlNodes []*graphml.Node) (map[string]*Node, error) {
	var nodes map[string]*Node = make(map[string]*Node, len(gmlNodes))
	for i := 0; i < len(gmlNodes); i++ {
		node := NewNode(gmlNodes[i].ID)
		nodes[node.NodeID()] = node
	}

	log.Log.Debug().Msg("parsed GraphML nodes")
	return nodes, nil
}

func graphMLEdges(nodes map[string]*Node, gmlEdges []*graphml.Edge) (map[string]*Edge, error) {
	var edges map[string]*Edge = make(map[string]*Edge, len(gmlEdges))
	for i := 0; i < len(gmlEdges); i++ {
		edge, err := parseGraphMLEdge(nodes, gmlEdges[i])
		if err != nil {
			log.Log.Error().Err(err).Str("id", gmlEdges[i].ID).Msg("error parsing GraphML edge")
			return nil, err
		}

		edges[edge.ID()] = edge
	}

	log.Log.Debug().Msg("parsed GraphML edges")
	return edges, nil
}

func parseGraphMLEdge(nodes map[string]*Node, gmlEdge *graphml.Edge) (*Edge, error) {
	aNode, exists := nodes[gmlEdge.Source]
	if !exists {
		log.Log.Error().Err(domain.ErrInvalidTopology).Str("id", gmlEdge.ID).Msg("GraphML edge missing source node")
		return nil, domain.ErrInvalidTopology
	}

	bNode, exists := nodes[gmlEdge.Target]
	if !exists {
		log.Log.Error().Err(domain.ErrInvalidTopology).Str("id", gmlEdge.ID).Msg("GraphML edge missing target node")
		return nil, domain.ErrInvalidTopology
	}

	log.Log.Trace().Str("id", gmlEdge.ID).Msg("parsed GraphML edge")
	return NewEdge(gmlEdge.ID, aNode.NodeID(), bNode.NodeID()), nil
}
