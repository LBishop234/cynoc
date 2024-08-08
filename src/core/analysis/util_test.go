package analysis

import (
	"strconv"
	"testing"

	"main/src/domain"
	"main/src/topology"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstructAnalysisTFs(t *testing.T) {
	t.Skip("Skipping to prevent replication of TestNewTrafficFlowAndRoute")
}

func TestNewAnalysisTF(t *testing.T) {
	t.Parallel()

	top := topology.ThreeByThreeMesh(t)

	tf1 := domain.TrafficFlowConfig{
		ID:    "tf1",
		Route: "[n3,n4,n5]",
	}

	expectedRoute := domain.Route{
		top.Nodes()["n3"].NodeID(),
		top.Nodes()["n4"].NodeID(),
		top.Nodes()["n5"].NodeID(),
	}

	tfr, err := newAnalysisTF(top, tf1)
	require.NoError(t, err)

	assert.Equal(t, tf1, tfr.TrafficFlowConfig)
	assert.Equal(t, expectedRoute, tfr.Route)
}

func TestSortTfsByPriority(t *testing.T) {
	t.Parallel()

	type testCase struct {
		tfs           []domain.TrafficFlowConfig
		expectedIndex []int
	}

	testCases := []testCase{
		{
			tfs: []domain.TrafficFlowConfig{
				{ID: "tf1", Priority: 1},
				{ID: "tf2", Priority: 2},
				{ID: "tf3", Priority: 3},
			},
			expectedIndex: []int{0, 1, 2},
		},
		{
			tfs: []domain.TrafficFlowConfig{
				{ID: "tf1", Priority: 3},
				{ID: "tf2", Priority: 2},
				{ID: "tf3", Priority: 1},
			},
			expectedIndex: []int{2, 1, 0},
		},
		{
			tfs: []domain.TrafficFlowConfig{
				{ID: "tf1", Priority: 2},
				{ID: "tf2", Priority: 3},
				{ID: "tf3", Priority: 1},
			},
			expectedIndex: []int{2, 0, 1},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			expectedTFs := make([]domain.TrafficFlowConfig, len(tc.tfs))
			for i, idx := range tc.expectedIndex {
				expectedTFs[i] = tc.tfs[idx]
			}

			tfs := sortTfsByPriority(tc.tfs)

			assert.Equal(t, expectedTFs, tfs)
		})
	}
}
