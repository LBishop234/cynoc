package domain

type RoutingAlgorithm string

const (
	XYRouting RoutingAlgorithm = "XY"
)

type NodeID struct {
	ID string
}

type Route []NodeID

// Returns an array of all valid routing algorithms.
func RoutingAlgorithms() []RoutingAlgorithm {
	return []RoutingAlgorithm{
		XYRouting,
	}
}
