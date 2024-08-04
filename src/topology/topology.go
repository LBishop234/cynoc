package topology

import (
	"path/filepath"

	"main/log"
	"main/src/domain"
)

type Topology struct {
	nodes map[string]*Node
	edges map[string]*Edge
}

type Node string

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

func NewTopology(nodes map[string]*Node, edges map[string]*Edge) *Topology {
	top := &Topology{
		nodes: nodes,
		edges: edges,
	}

	return top
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

func (t *Topology) Route(nodes []string) (domain.Route, error) {
	if len(nodes) < 2 {
		return nil, domain.ErrInvalidRoute
	}

	route := make(domain.Route, len(nodes))
	for i := 0; i < len(nodes); i++ {
		nodeID, exists := t.Node(nodes[i])
		if !exists {
			return nil, domain.ErrInvalidRoute
		}

		route[i] = nodeID.NodeID()
	}

	return route, nil
}

func NewNode(id string) *Node {
	log.Log.Trace().Str("id", id).Msg("new node")

	node := Node(id)
	return &node
}

func (n *Node) NodeID() string {
	return string(*n)
}

func NewEdge(id string, source, target string) *Edge {
	log.Log.Trace().Str("id", id).Msg("new edge")

	return &Edge{
		id: id,
		a:  source,
		b:  target,
	}
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
