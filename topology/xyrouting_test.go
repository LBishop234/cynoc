package topology

import (
	"testing"

	"main/domain"
)

func TestTopologyXYRoute(t *testing.T) {
	type testCase struct {
		src   domain.NodeID
		dst   domain.NodeID
		route []domain.NodeID
	}

	testCases := []testCase{
		{
			src: domain.NodeID{
				ID:  "n3",
				Pos: domain.NewPosition(0, 1),
			},
			dst: domain.NodeID{
				ID:  "n5",
				Pos: domain.NewPosition(2, 1),
			},
			route: []domain.NodeID{
				{
					ID:  "n3",
					Pos: domain.NewPosition(0, 1),
				},
				{
					ID:  "n4",
					Pos: domain.NewPosition(1, 1),
				},
				{
					ID:  "n5",
					Pos: domain.NewPosition(2, 1),
				},
			},
		},
		{
			src: domain.NodeID{
				ID:  "n1",
				Pos: domain.NewPosition(1, 0),
			},
			dst: domain.NodeID{
				ID:  "n7",
				Pos: domain.NewPosition(1, 2),
			},
			route: []domain.NodeID{
				{
					ID:  "n1",
					Pos: domain.NewPosition(1, 0),
				},
				{
					ID:  "n4",
					Pos: domain.NewPosition(1, 1),
				},
				{
					ID:  "n7",
					Pos: domain.NewPosition(1, 2),
				},
			},
		},
		{
			src: domain.NodeID{
				ID:  "n0",
				Pos: domain.NewPosition(0, 0),
			},
			dst: domain.NodeID{
				ID:  "n8",
				Pos: domain.NewPosition(2, 2),
			},
			route: []domain.NodeID{
				{
					ID:  "n0",
					Pos: domain.NewPosition(0, 0),
				},
				{
					ID:  "n1",
					Pos: domain.NewPosition(1, 0),
				},
				{
					ID:  "n2",
					Pos: domain.NewPosition(2, 0),
				},
				{
					ID:  "n5",
					Pos: domain.NewPosition(2, 1),
				},
				{
					ID:  "n8",
					Pos: domain.NewPosition(2, 2),
				},
			},
		},
	}

	top := ThreeByThreeMesh(t)

	for _, tc := range testCases {
		route, err := top.XYRoute(tc.src, tc.dst)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(route) != len(tc.route) {
			t.Errorf("Expected route length %d, got %d", len(tc.route), len(route))
		}

		for i := range route {
			if route[i] != tc.route[i] {
				t.Errorf("Expected route[%d] to be %v, got %v", i, tc.route[i], route[i])
			}
		}
	}
}
