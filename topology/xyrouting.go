package topology

import (
	"math"

	"main/domain"
)

func (t *Topology) XYRoute(src, dst domain.NodeID) (domain.Route, error) {
	var route domain.Route = domain.Route{src}

	current := src
	for current != dst {
		if current.Pos.X() != dst.Pos.X() {
			baseXDiff := math.Abs(float64(dst.Pos.X() - current.Pos.X()))

			leftXDiff := math.Abs(float64(dst.Pos.X() - (current.Pos.X() - 1)))
			rightXDiff := math.Abs(float64(dst.Pos.X() - (current.Pos.X() + 1)))

			var diff int
			if leftXDiff < baseXDiff {
				diff = -1
			} else if rightXDiff < baseXDiff {
				diff = +1
			}

			nextNode, exists := t.nodeByPos[domain.NewPosition(current.Pos.X()+diff, current.Pos.Y())]
			if !exists {
				nextNode, exists = t.nodeByPos[domain.NewPosition(current.Pos.X()-diff, current.Pos.Y())]
				if !exists {
					return nil, domain.ErrMissingRouter
				}
			}
			route = append(route, nextNode.nodeID)
		} else if current.Pos.Y() != dst.Pos.Y() {
			baseYDiff := math.Abs(float64(dst.Pos.Y() - current.Pos.Y()))

			upYDiff := math.Abs(float64(dst.Pos.Y() - (current.Pos.Y() - 1)))
			downYDiff := math.Abs(float64(dst.Pos.Y() - (current.Pos.Y() + 1)))

			var diff int
			if upYDiff < baseYDiff {
				diff = -1
			} else if downYDiff < baseYDiff {
				diff = +1
			}

			nextNode, exists := t.nodeByPos[domain.NewPosition(current.Pos.X(), current.Pos.Y()+diff)]
			if !exists {
				nextNode, exists = t.nodeByPos[domain.NewPosition(current.Pos.X(), current.Pos.Y()-diff)]
				if !exists {
					return nil, domain.ErrMissingRouter
				}
			}
			route = append(route, nextNode.nodeID)
		}

		current = route[len(route)-1]
	}

	return route, nil
}
