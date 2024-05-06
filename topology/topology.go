package topology

import (
	"path/filepath"

	"main/domain"
	"main/log"
)

type Topology struct {
	nodes     map[string]*Node
	nodeByPos map[domain.Position]*Node
	edges     map[string]*Edge
}

type Node struct {
	nodeID domain.NodeID
}

type Edge struct {
	id string
	a  string
	b  string
}

func ReadTopology(fPath string) (*Topology, error) {
	var top *Topology
	var err error

	log.Log.Debug().Msg("reading topology file")

	switch filepath.Ext(fPath) {
	case ".xml":
		top, err = graphML(fPath)
	default:
		log.Log.Error().Err(domain.ErrInvalidFilepath).Str("ext", filepath.Ext(fPath)).Msg("invalid topology file extension")
		return nil, domain.ErrInvalidFilepath
	}

	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error reading topology file")
		return nil, err
	}

	log.Log.Info().Msg("loaded topology from file")
	return top, nil
}

func NewTopology(nodes map[string]*Node, edges map[string]*Edge) (*Topology, error) {
	top := &Topology{
		nodes: nodes,
		edges: edges,
	}

	nodesByPos := make(map[domain.Position]*Node)
	for id := range nodes {
		nodesByPos[nodes[id].NodeID().Pos] = nodes[id]
	}
	top.nodeByPos = nodesByPos

	return top, nil
}

func (t *Topology) Nodes() map[string]*Node {
	return t.nodes
}

func (t *Topology) Node(id string) (*Node, bool) {
	node, ok := t.nodes[id]
	return node, ok
}

func (t *Topology) Edges() map[string]*Edge {
	return t.edges
}

func (t *Topology) Edge(id string) (*Edge, bool) {
	edge, ok := t.edges[id]
	return edge, ok
}

func NewNode(id string, pos domain.Position) (*Node, error) {
	log.Log.Trace().Str("id", id).Msg("new node")

	return &Node{
		nodeID: domain.NodeID{
			ID:  id,
			Pos: pos,
		},
	}, nil
}

func (n *Node) NodeID() domain.NodeID {
	return n.nodeID
}

func NewEdge(id string, source, target string) (*Edge, error) {
	log.Log.Trace().Str("id", id).Msg("new edge")

	return &Edge{
		id: id,
		a:  source,
		b:  target,
	}, nil
}

func (e *Edge) ID() string {
	return e.id
}

func (e *Edge) A() string {
	return e.a
}

func (e *Edge) B() string {
	return e.b
}
