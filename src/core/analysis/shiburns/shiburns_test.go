package shiburns

import (
	"testing"

	"main/src/core/analysis/util"
	"main/src/domain"
	"main/src/topology"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const XiongEtAl = "XiongEtAlNoC"

func testCasesTrafficFlowAndRoutes(tb testing.TB) map[string]map[string]util.TrafficFlowAndRoute {
	fourByFourTop := topology.FourByFourMesh(tb)

	tfAndRoutes := map[string]map[string]util.TrafficFlowAndRoute{
		XiongEtAl: {
			"t1": {
				TrafficFlowConfig: domain.TrafficFlowConfig{
					ID:         "t1",
					Priority:   1,
					Deadline:   100,
					Period:     500,
					Jitter:     26,
					PacketSize: 20,
					Route:      "[n3,n2,n1]",
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n3"].NodeID(),
					fourByFourTop.Nodes()["n2"].NodeID(),
					fourByFourTop.Nodes()["n1"].NodeID(),
				},
			},
			"t2": {
				TrafficFlowConfig: domain.TrafficFlowConfig{
					ID:         "t2",
					Priority:   2,
					Deadline:   107,
					Period:     407,
					Jitter:     33,
					PacketSize: 97,
					Route:      "[n8,n12]",
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n8"].NodeID(),
					fourByFourTop.Nodes()["n12"].NodeID(),
				},
			},
			"t3": {
				TrafficFlowConfig: domain.TrafficFlowConfig{
					ID:         "t3",
					Priority:   3,
					Deadline:   95,
					Period:     628,
					Jitter:     14,
					PacketSize: 36,
					Route:      "[n2,n1,n0,n4,n8,n12]",
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
			"t4": {
				TrafficFlowConfig: domain.TrafficFlowConfig{
					ID:         "t4",
					Priority:   4,
					Deadline:   124,
					Period:     1506,
					Jitter:     8,
					PacketSize: 58,
					Route:      "[n8,n12]",
				},
				Route: domain.Route{
					fourByFourTop.Nodes()["n8"].NodeID(),
					fourByFourTop.Nodes()["n12"].NodeID(),
				},
			},
			"t5": {
				TrafficFlowConfig: domain.TrafficFlowConfig{
					ID:         "t5",
					Priority:   5,
					Deadline:   189,
					Period:     689,
					Jitter:     27,
					PacketSize: 124,
					Route:      "[n1,n0,n4,n8]",
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

func TestShiBurns(t *testing.T) {
	t.Parallel()

	tfAndRoutes := testCasesTrafficFlowAndRoutes(t)

	type tfStruct struct {
		tfr util.TrafficFlowAndRoute
		res ShiBurnsResults
	}

	testCases := map[string]struct {
		simConf domain.SimConfig
		tfs     map[string]tfStruct
	}{
		XiongEtAl: {
			simConf: domain.SimConfig{
				CycleLimit:      2000,
				MaxPriority:     5,
				BufferSize:      25,
				ProcessingDelay: 6,
			},
			tfs: map[string]tfStruct{
				"t1": {
					tfr: tfAndRoutes[XiongEtAl]["t1"],
					res: ShiBurnsResults{
						DirectInterferenceCount:   0,
						IndirectInterferenceCount: 0,
						Latency:                   40,
					},
				},
				"t2": {
					tfr: tfAndRoutes[XiongEtAl]["t2"],
					res: ShiBurnsResults{
						DirectInterferenceCount:   0,
						IndirectInterferenceCount: 0,
						Latency:                   111,
					},
				},
				"t3": {
					tfr: tfAndRoutes[XiongEtAl]["t3"],
					res: ShiBurnsResults{
						DirectInterferenceCount:   2,
						IndirectInterferenceCount: 0,
						Latency:                   225,
					},
				},
				"t4": {
					tfr: tfAndRoutes[XiongEtAl]["t4"],
					res: ShiBurnsResults{
						DirectInterferenceCount:   2,
						IndirectInterferenceCount: 1,
						Latency:                   257,
					},
				},
				"t5": {
					tfr: tfAndRoutes[XiongEtAl]["t5"],
					res: ShiBurnsResults{
						DirectInterferenceCount:   1,
						IndirectInterferenceCount: 2,
						Latency:                   224,
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tfrs := make(map[string]util.TrafficFlowAndRoute, len(testCase.tfs))
			for tfKey := range testCase.tfs {
				tfrs[tfKey] = testCase.tfs[tfKey].tfr
			}

			for tfKey := range testCase.tfs {
				t.Run(tfKey, func(t *testing.T) {
					got, err := ShiBurns(testCase.simConf, tfrs, tfKey)

					require.NoError(t, err)
					assert.Equal(t, testCase.tfs[tfKey].res.DirectInterferenceCount, got.DirectInterferenceCount)
					assert.Equal(t, testCase.tfs[tfKey].res.IndirectInterferenceCount, got.IndirectInterferenceCount)
					assert.Equal(t, testCase.tfs[tfKey].res.Latency, got.Latency)
				})
			}
		})
	}
}

func BenchmarkShiBurns(b *testing.B) {
	tfAndRoutes := testCasesTrafficFlowAndRoutes(b)

	testCases := map[string]struct {
		simConf domain.SimConfig
		tfs     map[string]util.TrafficFlowAndRoute
	}{
		XiongEtAl: {
			simConf: domain.SimConfig{
				CycleLimit:      2000,
				MaxPriority:     5,
				BufferSize:      25,
				ProcessingDelay: 6,
			},
			tfs: tfAndRoutes[XiongEtAl],
		},
	}

	for name, testc := range testCases {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for tfKey := range testc.tfs {
					_, err := ShiBurns(testc.simConf, testc.tfs, tfKey)
					assert.NoError(b, err)
				}
			}
		})
	}
}

func TestFindInterferenceSet(t *testing.T) {
	t.Parallel()

	tfAndRoutes := testCasesTrafficFlowAndRoutes(t)

	testCases := map[string]map[string]struct {
		tfr         util.TrafficFlowAndRoute
		directInt   map[string]util.TrafficFlowAndRoute
		indirectInt map[string]util.TrafficFlowAndRoute
	}{
		XiongEtAl: {
			"t1": {
				tfr:         tfAndRoutes[XiongEtAl]["t1"],
				directInt:   map[string]util.TrafficFlowAndRoute{},
				indirectInt: map[string]util.TrafficFlowAndRoute{},
			},
			"t2": {
				tfr:         tfAndRoutes[XiongEtAl]["t2"],
				directInt:   map[string]util.TrafficFlowAndRoute{},
				indirectInt: map[string]util.TrafficFlowAndRoute{},
			},
			"t3": {
				tfr: tfAndRoutes[XiongEtAl]["t3"],
				directInt: map[string]util.TrafficFlowAndRoute{
					"t1": tfAndRoutes[XiongEtAl]["t1"],
					"t2": tfAndRoutes[XiongEtAl]["t2"],
				},
				indirectInt: map[string]util.TrafficFlowAndRoute{},
			},
			"t4": {
				tfr: tfAndRoutes[XiongEtAl]["t4"],
				directInt: map[string]util.TrafficFlowAndRoute{
					"t2": tfAndRoutes[XiongEtAl]["t2"],
					"t3": tfAndRoutes[XiongEtAl]["t3"],
				},
				indirectInt: map[string]util.TrafficFlowAndRoute{
					"t1": tfAndRoutes[XiongEtAl]["t1"],
				},
			},
			"t5": {
				tfr: tfAndRoutes[XiongEtAl]["t5"],
				directInt: map[string]util.TrafficFlowAndRoute{
					"t3": tfAndRoutes[XiongEtAl]["t3"],
				},
				indirectInt: map[string]util.TrafficFlowAndRoute{
					"t1": tfAndRoutes[XiongEtAl]["t1"],
					"t2": tfAndRoutes[XiongEtAl]["t2"],
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tfrs := make(map[string]util.TrafficFlowAndRoute, len(testCase))
			for tfKey := range testCase {
				tfrs[tfKey] = testCase[tfKey].tfr
			}

			for tfKey := range testCase {
				gotDIntSet, gotIIntSet := findInterferenceSets(tfrs, tfKey)

				assert.Equal(t, testCase[tfKey].directInt, gotDIntSet)
				assert.Equal(t, testCase[tfKey].indirectInt, gotIIntSet)
			}
		})
	}
}

func TestFilterByPriority(t *testing.T) {
	t.Parallel()

	var (
		tf1 = util.TrafficFlowAndRoute{
			TrafficFlowConfig: domain.TrafficFlowConfig{
				ID:       "tf1",
				Priority: 1,
			},
			Route: domain.Route{},
		}
		tf2 = util.TrafficFlowAndRoute{
			TrafficFlowConfig: domain.TrafficFlowConfig{
				ID:       "tf2",
				Priority: 2,
			},
			Route: domain.Route{},
		}
		tf3 = util.TrafficFlowAndRoute{
			TrafficFlowConfig: domain.TrafficFlowConfig{
				ID:       "tf3",
				Priority: 3,
			},
			Route: domain.Route{},
		}
	)

	tfs := map[string]util.TrafficFlowAndRoute{
		tf1.ID: tf1,
		tf2.ID: tf2,
		tf3.ID: tf3,
	}

	expectedFiltered := map[string]util.TrafficFlowAndRoute{
		tf1.ID: tf1,
		tf2.ID: tf2,
	}

	filtered := filterByPriority(tfs, tf2.Priority)

	assert.Equal(t, expectedFiltered, filtered)
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
