package domain

type RoutingAlgorithm string

const (
	XYRouting RoutingAlgorithm = "XY"
)

type NodeID struct {
	ID  string
	Pos Position
}

type Route []NodeID
