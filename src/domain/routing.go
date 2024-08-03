package domain

type RoutingAlgorithm string

const (
	XYRouting RoutingAlgorithm = "XY"
)

type Route []string

// Returns an array of all valid routing algorithms.
func RoutingAlgorithms() []RoutingAlgorithm {
	return []RoutingAlgorithm{
		XYRouting,
	}
}
