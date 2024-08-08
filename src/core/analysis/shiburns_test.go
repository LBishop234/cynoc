package analysis

import (
	"context"
	"strconv"
	"testing"

	"main/src/domain"
	"main/src/topology"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const XiongEtAl = "XiongEtAlNoC"

func testCasesTrafficFlowAndRoutes(tb testing.TB) map[string][]analysisTF {
	fourByFourTop := topology.FourByFourMesh(tb)

	tfAndRoutes := map[string][]analysisTF{
		XiongEtAl: {
			{
				TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
					TrafficFlowConfig: domain.TrafficFlowConfig{
						ID:         "t1",
						Priority:   1,
						Deadline:   100,
						Period:     500,
						Jitter:     26,
						PacketSize: 20,
						Route:      "[n3,n2,n1]",
					},
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n3"].NodeID(),
					fourByFourTop.Nodes()["n2"].NodeID(),
					fourByFourTop.Nodes()["n1"].NodeID(),
				},
			},
			{
				TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
					TrafficFlowConfig: domain.TrafficFlowConfig{
						ID:         "t2",
						Priority:   2,
						Deadline:   107,
						Period:     407,
						Jitter:     33,
						PacketSize: 97,
						Route:      "[n8,n12]",
					},
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n8"].NodeID(),
					fourByFourTop.Nodes()["n12"].NodeID(),
				},
			},
			{
				TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
					TrafficFlowConfig: domain.TrafficFlowConfig{
						ID:         "t3",
						Priority:   3,
						Deadline:   95,
						Period:     628,
						Jitter:     14,
						PacketSize: 36,
						Route:      "[n2,n1,n0,n4,n8,n12]",
					},
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n2"].NodeID(),
					fourByFourTop.Nodes()["n1"].NodeID(),
					fourByFourTop.Nodes()["n0"].NodeID(),
					fourByFourTop.Nodes()["n4"].NodeID(),
					fourByFourTop.Nodes()["n8"].NodeID(),
					fourByFourTop.Nodes()["n12"].NodeID(),
				},
			},
			{
				TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
					TrafficFlowConfig: domain.TrafficFlowConfig{
						ID:         "t4",
						Priority:   4,
						Deadline:   124,
						Period:     1506,
						Jitter:     8,
						PacketSize: 58,
						Route:      "[n8,n12]",
					},
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n8"].NodeID(),
					fourByFourTop.Nodes()["n12"].NodeID(),
				},
			},
			{
				TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
					TrafficFlowConfig: domain.TrafficFlowConfig{
						ID:         "t5",
						Priority:   5,
						Deadline:   189,
						Period:     689,
						Jitter:     27,
						PacketSize: 124,
						Route:      "[n1,n0,n4,n8]",
					},
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n1"].NodeID(),
					fourByFourTop.Nodes()["n0"].NodeID(),
					fourByFourTop.Nodes()["n4"].NodeID(),
					fourByFourTop.Nodes()["n8"].NodeID(),
				},
			},
		},
	}

	return tfAndRoutes
}

func TestNewShiBurns(t *testing.T) {
	t.Parallel()

	type testCase struct {
		conf     domain.SimConfig
		tfsMapID string
		expected []int
	}

	testCases := []testCase{
		{
			conf: domain.SimConfig{
				CycleLimit:      2000,
				MaxPriority:     5,
				BufferSize:      25,
				ProcessingDelay: 6,
			},
			tfsMapID: XiongEtAl,
			expected: []int{38, 109, 219, 251, 220},
		},
	}

	for tcIndex, tc := range testCases {
		t.Run(strconv.Itoa(tcIndex), func(t *testing.T) {
			aTfs, err := basicLatency(context.TODO(), tc.conf, testCasesTrafficFlowAndRoutes(t)[tc.tfsMapID])
			require.NoError(t, err)

			aTFs, err := shiBurns(context.TODO(), aTfs)
			require.NoError(t, err)

			assert.Len(t, aTFs, len(tc.expected))
			for i := 0; i < len(aTFs); i++ {
				assert.Equal(t, tc.expected[i], aTFs[i].ShiAndBurns)
			}
		})
	}
}

func BenchmarkShiBurns(b *testing.B) {
	type testCase struct {
		conf     domain.SimConfig
		tfsMapID string
	}

	testCases := []testCase{
		{
			conf: domain.SimConfig{
				CycleLimit:      2000,
				MaxPriority:     5,
				BufferSize:      25,
				ProcessingDelay: 6,
			},
			tfsMapID: XiongEtAl,
		},
	}

	for tcIndex, tc := range testCases {
		b.Run(strconv.Itoa(tcIndex), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				aTfs, err := basicLatency(context.TODO(), tc.conf, testCasesTrafficFlowAndRoutes(b)[tc.tfsMapID])
				require.NoError(b, err)

				_, err = shiBurns(context.TODO(), aTfs)
				require.NoError(b, err)
			}
		})
	}
}

func TestNewFindInterferenceSets(t *testing.T) {
	t.Parallel()

	type testCase struct {
		tfs             []analysisTF
		expectedDIntSet []map[string]int
		expectedIIntSet []map[string]int
	}

	testCases := []testCase{
		{
			tfs: []analysisTF{
				{
					TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
						TrafficFlowConfig: domain.TrafficFlowConfig{
							ID:       "tf1",
							Priority: 1,
							Route:    "[n0,n1,n2]",
						},
					},
					Route: domain.Route{"n0", "n1", "n2"},
				},
				{
					TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
						TrafficFlowConfig: domain.TrafficFlowConfig{
							ID:       "tf2",
							Priority: 2,
							Route:    "[n1,n2,n3]",
						},
					},
					Route: domain.Route{"n1", "n2", "n3"},
				},
				{
					TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
						TrafficFlowConfig: domain.TrafficFlowConfig{
							ID:       "tf3",
							Priority: 3,
							Route:    "[n2,n3,n4]",
						},
					},
					Route: domain.Route{"n2", "n3", "n4"},
				},
				{
					TrafficFlowAnalysisSet: domain.TrafficFlowAnalysisSet{
						TrafficFlowConfig: domain.TrafficFlowConfig{
							ID:       "tf4",
							Priority: 4,
							Route:    "[n3,n4,n5]",
						},
					},
					Route: domain.Route{"n3", "n4", "n5"},
				},
			},
			expectedDIntSet: []map[string]int{
				{},
				{"tf1": 0},
				{"tf2": 1},
				{"tf3": 2},
			},
			expectedIIntSet: []map[string]int{
				{},
				{},
				{"tf1": 0},
				{"tf1": 0, "tf2": 1},
			},
		},
	}

	for tcIndex, tc := range testCases {
		t.Run(strconv.Itoa(tcIndex), func(t *testing.T) {
			aTFs := findIntereferenceSets(tc.tfs)
			for i := 0; i < len(tc.tfs); i++ {
				assert.Equal(t, tc.expectedDIntSet[i], aTFs[i].directIntSet)
				assert.Equal(t, tc.expectedIIntSet[i], aTFs[i].indirectIntSet)
			}
		})
	}
}

func TestIntersectingRoutes(t *testing.T) {
	t.Run("NotIntersecting", func(t *testing.T) {
		t.Parallel()

		top := topology.ThreeByThreeMesh(t)

		r1 := domain.Route{
			top.Nodes()["n0"].NodeID(),
			top.Nodes()["n1"].NodeID(),
			top.Nodes()["n4"].NodeID(),
			top.Nodes()["n7"].NodeID(),
		}
		r2 := domain.Route{
			top.Nodes()["n3"].NodeID(),
			top.Nodes()["n4"].NodeID(),
			top.Nodes()["n5"].NodeID(),
		}

		assert.False(t, intersectingRoutes(r1, r2))
	})

	t.Run("Intersecting", func(t *testing.T) {
		t.Parallel()

		top := topology.ThreeByThreeMesh(t)

		r1 := domain.Route{
			top.Nodes()["n0"].NodeID(),
			top.Nodes()["n1"].NodeID(),
			top.Nodes()["n4"].NodeID(),
			top.Nodes()["n7"].NodeID(),
		}
		r2 := domain.Route{
			top.Nodes()["n2"].NodeID(),
			top.Nodes()["n1"].NodeID(),
			top.Nodes()["n4"].NodeID(),
		}

		assert.True(t, intersectingRoutes(r1, r2))
	})
}
