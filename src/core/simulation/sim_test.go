package simulation

import (
	"context"
	"math"
	"testing"

	"main/src/core/network"
	"main/src/domain"
	"main/src/topology"
	"main/src/traffic"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	OnePriorityConfig = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      1,
		BufferSize:       2,
		FlitSize:         1,
		ProcessingDelay:  3,
		LinkBandwidth:    1,
	}

	TwoPriorityConfig = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      2,
		BufferSize:       4,
		FlitSize:         1,
		ProcessingDelay:  3,
		LinkBandwidth:    1,
	}

	TwoPriorityConfig2LinkBandwidth = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      2,
		BufferSize:       4,
		FlitSize:         1,
		ProcessingDelay:  3,
		LinkBandwidth:    2,
	}

	FourPriorityConfig = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      4,
		BufferSize:       32,
		FlitSize:         4,
		ProcessingDelay:  6,
		LinkBandwidth:    1,
	}

	TenPriorityConfig = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      10,
		BufferSize:       80,
		FlitSize:         4,
		ProcessingDelay:  6,
		LinkBandwidth:    1,
	}

	TwentyPriorityConfig = domain.SimConfig{
		RoutingAlgorithm: domain.XYRouting,
		MaxPriority:      20,
		BufferSize:       160,
		FlitSize:         4,
		ProcessingDelay:  6,
		LinkBandwidth:    1,
	}
)

type templateTestCase struct {
	templateRun  bool
	cycles       int
	topologyFunc func(t testing.TB) *topology.Topology
	networkConf  domain.SimConfig
	traffic      []domain.TrafficFlowConfig
}

var templateTestCases = map[string]templateTestCase{
	"3hLineOnePkt": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeHorozontalLine,
		networkConf:  OnePriorityConfig,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     50,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
			},
		},
	},
	"3hLineTwoPkts": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeHorozontalLine,
		networkConf:  TwoPriorityConfig,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     50,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
			},
			{
				ID:         "t2",
				Src:        "n0",
				Dst:        "n2",
				Priority:   2,
				Period:     math.MaxInt,
				Deadline:   math.MaxInt,
				Jitter:     0,
				PacketSize: 2,
			},
		},
	},
	"3hLineTwoPkts2LinkBandwidth": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeHorozontalLine,
		networkConf:  TwoPriorityConfig2LinkBandwidth,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     50,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
			},
			{
				ID:         "t2",
				Src:        "n0",
				Dst:        "n2",
				Priority:   2,
				Period:     math.MaxInt,
				Deadline:   math.MaxInt,
				Jitter:     0,
				PacketSize: 2,
			},
		},
	},
	"3x3Mesh4Pkts": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeByThreeMesh,
		networkConf:  TwoPriorityConfig,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     100,
				Deadline:   75,
				Jitter:     0,
				PacketSize: 32,
			},
			{
				ID:         "t2",
				Src:        "n1",
				Dst:        "n8",
				Priority:   2,
				Period:     120,
				Deadline:   105,
				Jitter:     0,
				PacketSize: 40,
			},
			{
				ID:         "t3",
				Src:        "n5",
				Dst:        "n1",
				Priority:   3,
				Period:     100,
				Deadline:   85,
				Jitter:     0,
				PacketSize: 32,
			},
			{
				ID:         "t4",
				Src:        "n3",
				Dst:        "n8",
				Priority:   4,
				Period:     100,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 32,
			},
		},
	},
	"3x3Mesh10Pkts": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeByThreeMesh,
		networkConf:  TenPriorityConfig,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     100,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 32,
			},
			{
				ID:         "t2",
				Src:        "n1",
				Dst:        "n8",
				Priority:   2,
				Period:     120,
				Deadline:   110,
				Jitter:     0,
				PacketSize: 40,
			},
			{
				ID:         "t3",
				Src:        "n3",
				Dst:        "n7",
				Priority:   3,
				Period:     150,
				Deadline:   120,
				Jitter:     0,
				PacketSize: 64,
			},
			{
				ID:         "t4",
				Src:        "n6",
				Dst:        "n5",
				Priority:   4,
				Period:     100,
				Deadline:   75,
				Jitter:     0,
				PacketSize: 32,
			},
			{
				ID:         "t5",
				Src:        "n2",
				Dst:        "n3",
				Priority:   5,
				Period:     75,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 8,
			},
			{
				ID:         "t6",
				Src:        "n4",
				Dst:        "n3",
				Priority:   6,
				Period:     100,
				Deadline:   100,
				Jitter:     0,
				PacketSize: 96,
			},
			{
				ID:         "t7",
				Src:        "n6",
				Dst:        "n0",
				Priority:   7,
				Period:     120,
				Deadline:   105,
				Jitter:     0,
				PacketSize: 16,
			},
			{
				ID:         "t8",
				Src:        "n7",
				Dst:        "n2",
				Priority:   8,
				Period:     120,
				Deadline:   100,
				Jitter:     0,
				PacketSize: 32,
			},
			{
				ID:         "t9",
				Src:        "n4",
				Dst:        "n2",
				Priority:   9,
				Period:     100,
				Deadline:   85,
				Jitter:     0,
				PacketSize: 48,
			},
			{
				ID:         "t10",
				Src:        "n8",
				Dst:        "n0",
				Priority:   10,
				Period:     110,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 40,
			},
		},
	},
	"3x3Mesh20Pkts": {
		templateRun:  true,
		cycles:       1000,
		topologyFunc: topology.ThreeByThreeMesh,
		networkConf:  TwentyPriorityConfig,
		traffic: []domain.TrafficFlowConfig{
			{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				Priority:   1,
				Period:     100,
				Deadline:   80,
				Jitter:     10,
				PacketSize: 32,
			},
			{
				ID:         "t2",
				Src:        "n1",
				Dst:        "n8",
				Priority:   2,
				Period:     120,
				Deadline:   115,
				Jitter:     3,
				PacketSize: 40,
			},
			{
				ID:         "t3",
				Src:        "n3",
				Dst:        "n7",
				Priority:   3,
				Period:     150,
				Deadline:   135,
				Jitter:     6,
				PacketSize: 64,
			},
			{
				ID:         "t4",
				Src:        "n6",
				Dst:        "n5",
				Priority:   4,
				Period:     100,
				Deadline:   90,
				Jitter:     5,
				PacketSize: 32,
			},
			{
				ID:         "t5",
				Src:        "n2",
				Dst:        "n3",
				Priority:   5,
				Period:     120,
				Deadline:   100,
				Jitter:     1,
				PacketSize: 8,
			},
			{
				ID:         "t6",
				Src:        "n4",
				Dst:        "n3",
				Priority:   6,
				Period:     200,
				Deadline:   180,
				Jitter:     8,
				PacketSize: 96,
			},
			{
				ID:         "t7",
				Src:        "n6",
				Dst:        "n0",
				Priority:   7,
				Period:     70,
				Deadline:   50,
				Jitter:     2,
				PacketSize: 16,
			},
			{
				ID:         "t8",
				Src:        "n7",
				Dst:        "n2",
				Priority:   8,
				Period:     50,
				Deadline:   45,
				Jitter:     5,
				PacketSize: 32,
			},
			{
				ID:         "t9",
				Src:        "n4",
				Dst:        "n2",
				Priority:   9,
				Period:     100,
				Deadline:   70,
				Jitter:     9,
				PacketSize: 48,
			},
			{
				ID:         "t10",
				Src:        "n8",
				Dst:        "n0",
				Priority:   10,
				Period:     120,
				Deadline:   90,
				Jitter:     14,
				PacketSize: 40,
			},
			{
				ID:         "t11",
				Src:        "n5",
				Dst:        "n1",
				Priority:   11,
				Period:     100,
				Deadline:   85,
				Jitter:     7,
				PacketSize: 64,
			},
			{
				ID:         "t12",
				Src:        "n7",
				Dst:        "n4",
				Priority:   12,
				Period:     120,
				Deadline:   115,
				Jitter:     4,
				PacketSize: 32,
			},
			{
				ID:         "t13",
				Src:        "n8",
				Dst:        "n6",
				Priority:   13,
				Period:     100,
				Deadline:   95,
				Jitter:     10,
				PacketSize: 48,
			},
			{
				ID:         "t14",
				Src:        "n5",
				Dst:        "n7",
				Priority:   14,
				Period:     75,
				Deadline:   70,
				Jitter:     13,
				PacketSize: 40,
			},
			{
				ID:         "t15",
				Src:        "n3",
				Dst:        "n5",
				Priority:   15,
				Period:     100,
				Deadline:   85,
				Jitter:     12,
				PacketSize: 64,
			},
			{
				ID:         "t16",
				Src:        "n1",
				Dst:        "n6",
				Priority:   16,
				Period:     120,
				Deadline:   100,
				Jitter:     11,
				PacketSize: 32,
			},
			{
				ID:         "t17",
				Src:        "n0",
				Dst:        "n1",
				Priority:   17,
				Period:     100,
				Deadline:   100,
				Jitter:     15,
				PacketSize: 48,
			},
			{
				ID:         "t18",
				Src:        "n2",
				Dst:        "n8",
				Priority:   18,
				Period:     115,
				Deadline:   100,
				Jitter:     16,
				PacketSize: 40,
			},
			{
				ID:         "t19",
				Src:        "n4",
				Dst:        "n5",
				Priority:   19,
				Period:     110,
				Deadline:   100,
				Jitter:     17,
				PacketSize: 64,
			},
			{
				ID:         "t20",
				Src:        "n7",
				Dst:        "n3",
				Priority:   20,
				Period:     110,
				Deadline:   95,
				Jitter:     18,
				PacketSize: 32,
			},
		},
	},
}

func BenchmarkNewSimulator(b *testing.B) {
	for name, testCase := range templateTestCases {
		b.Run(name, func(b *testing.B) {
			if !testCase.templateRun {
				b.Skip()
			}

			// Setup & Run Simulation
			network, err := network.NewNetwork(testCase.topologyFunc(b), testCase.networkConf)
			require.NoError(b, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i])
				require.NoError(b, err)

				trafficFlows[i] = tf
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := newSimulator(network, trafficFlows, domain.XYRouting, testCase.cycles)
				require.NoError(b, err)
			}
		})
	}
}

func TestRunSimulation(t *testing.T) {
	type (
		expectedPkts struct {
			trafficFlowID string
			cycle         float64
		}

		testCase struct {
			run bool
			templateTestCase
			expected []expectedPkts
		}
	)

	testCases := map[string]testCase{
		"3hLineOnePkt": {
			run:              false,
			templateTestCase: templateTestCases["3hLineOnePkt"],
			expected: []expectedPkts{
				{
					trafficFlowID: "t1",
					cycle:         12,
				},
			},
		},
		"3hLineTwoPkts": {
			run:              false,
			templateTestCase: templateTestCases["3hLineTwoPkts"],
			expected: []expectedPkts{
				{
					trafficFlowID: "t1",
					cycle:         12,
				},
				{
					trafficFlowID: "t2",
					cycle:         16,
				},
			},
		},
		"3hLineTwoPkts2LinkBandwidth": {
			run:              true,
			templateTestCase: templateTestCases["3hLineTwoPkts2LinkBandwidth"],
			expected: []expectedPkts{
				{
					trafficFlowID: "t1",
					cycle:         11,
				},
				{
					trafficFlowID: "t2",
					cycle:         12,
				},
			},
		},
	}

	// Restrict cycle limit for expected packets
	for name, testCase := range testCases {
		testCase.cycles = 20
		testCases[name] = testCase
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if !testCase.run || !testCase.templateRun {
				t.Skip()
			}

			t.Parallel()

			// Setup & Run Simulation
			network, err := network.NewNetwork(testCase.topologyFunc(t), testCase.networkConf)
			require.NoError(t, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i])
				require.NoError(t, err)

				trafficFlows[i] = tf
			}

			simulator, err := newSimulator(network, trafficFlows, domain.XYRouting, testCase.cycles)
			require.NoError(t, err)

			_, records, err := simulator.runSimulation(context.Background())
			require.NoError(t, err)

			// Check Expected Packets
			for i := 0; i < len(testCase.expected); i++ {
				expctTfID := testCase.expected[i].trafficFlowID

				assert.NotEmpty(t, records.ArrivedByTF[expctTfID], 0)

				found := false
				for id, pkt := range records.ArrivedByTF[expctTfID] {
					if pkt.ReceivedCycle == testCase.expected[i].cycle {
						found = true
						delete(records.ArrivedByTF[expctTfID], id)
						break
					}
				}

				assert.True(t, found, "expected packet:\n%s", spew.Sdump(testCase.expected[i]))
			}
		})
	}
}

func BenchmarkRunSimulation(b *testing.B) {
	type testCase struct {
		run bool
		templateTestCase
	}

	testCases := map[string]testCase{
		"3hLineOnePkt": {
			run:              true,
			templateTestCase: templateTestCases["3hLineOnePkt"],
		},
		"3x3Mesh4Pkts": {
			run:              true,
			templateTestCase: templateTestCases["3x3Mesh4Pkts"],
		},
		"3x3Mesh20Pkts": {
			run:              true,
			templateTestCase: templateTestCases["3x3Mesh20Pkts"],
		},
	}

	for name, testCase := range testCases {
		b.Run(name, func(b *testing.B) {
			if !testCase.run || !testCase.templateRun {
				b.Skip()
			}

			// Setup & Run Simulation
			network, err := network.NewNetwork(testCase.topologyFunc(b), testCase.networkConf)
			require.NoError(b, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i])
				require.NoError(b, err)

				trafficFlows[i] = tf
			}

			simulator, err := newSimulator(network, trafficFlows, domain.XYRouting, testCase.cycles)
			require.NoError(b, err)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _, err := simulator.runSimulation(context.Background())
				require.NoError(b, err)
			}
		})
	}
}
