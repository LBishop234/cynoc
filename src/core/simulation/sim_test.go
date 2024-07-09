package simulation

import (
	"context"
	"io"
	"math"
	"testing"

	"main/src/core/network"
	"main/src/domain"
	"main/src/topology"
	"main/src/traffic"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	OnePriorityConfig = domain.SimConfig{
		MaxPriority:     1,
		BufferSize:      2,
		ProcessingDelay: 3,
	}

	TwoPriorityConfig = domain.SimConfig{
		MaxPriority:     2,
		BufferSize:      4,
		ProcessingDelay: 3,
	}

	FourPriorityConfig = domain.SimConfig{
		MaxPriority:     4,
		BufferSize:      8,
		ProcessingDelay: 6,
	}

	TenPriorityConfig = domain.SimConfig{
		MaxPriority:     10,
		BufferSize:      20,
		ProcessingDelay: 6,
	}

	TwentyPriorityConfig = domain.SimConfig{
		MaxPriority:     20,
		BufferSize:      40,
		ProcessingDelay: 6,
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
				Priority:   1,
				Period:     50,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
				Route:      "[n0,n1,n2]",
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
				Priority:   1,
				Period:     50,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
				Route:      "[n0,n1,n2]",
			},
			{
				ID:         "t2",
				Priority:   2,
				Period:     math.MaxInt,
				Deadline:   math.MaxInt,
				Jitter:     0,
				PacketSize: 2,
				Route:      "[n0,n1,n2]",
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
				Priority:   1,
				Period:     100,
				Deadline:   75,
				Jitter:     0,
				PacketSize: 32,
				Route:      "[n0,n1,n2]",
			},
			{
				ID:         "t2",
				Priority:   2,
				Period:     120,
				Deadline:   105,
				Jitter:     0,
				PacketSize: 40,
				Route:      "[n1,n2,n5,n8]",
			},
			{
				ID:         "t3",
				Priority:   3,
				Period:     100,
				Deadline:   85,
				Jitter:     0,
				PacketSize: 32,
				Route:      "[n5,n2,n1]",
			},
			{
				ID:         "t4",
				Priority:   4,
				Period:     100,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 32,
				Route:      "[n3,n4,n5,n8]",
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
				Priority:   1,
				Period:     100,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 8,
				Route:      "[n0,n1,n2]",
			},
			{
				ID:         "t2",
				Priority:   2,
				Period:     120,
				Deadline:   110,
				Jitter:     0,
				PacketSize: 10,
				Route:      "[n1,n2,n5,n8]",
			},
			{
				ID:         "t3",
				Priority:   3,
				Period:     150,
				Deadline:   120,
				Jitter:     0,
				PacketSize: 16,
				Route:      "[n3,n4,n7]",
			},
			{
				ID:         "t4",
				Priority:   4,
				Period:     100,
				Deadline:   75,
				Jitter:     0,
				PacketSize: 8,
				Route:      "[n6,n7,n8,n5]",
			},
			{
				ID:         "t5",
				Priority:   5,
				Period:     75,
				Deadline:   50,
				Jitter:     0,
				PacketSize: 2,
				Route:      "[n2,n1,n0,n3]",
			},
			{
				ID:         "t6",
				Priority:   6,
				Period:     100,
				Deadline:   100,
				Jitter:     0,
				PacketSize: 24,
				Route:      "[n4,n3]",
			},
			{
				ID:         "t7",
				Priority:   7,
				Period:     120,
				Deadline:   105,
				Jitter:     0,
				PacketSize: 4,
				Route:      "[n6,n3,n0]",
			},
			{
				ID:         "t8",
				Priority:   8,
				Period:     120,
				Deadline:   100,
				Jitter:     0,
				PacketSize: 8,
				Route:      "[n7,n8,n5,n2]",
			},
			{
				ID:         "t9",
				Priority:   9,
				Period:     100,
				Deadline:   85,
				Jitter:     0,
				PacketSize: 12,
				Route:      "[n4,n5,n2]",
			},
			{
				ID:         "t10",
				Priority:   10,
				Period:     110,
				Deadline:   90,
				Jitter:     0,
				PacketSize: 10,
				Route:      "[n8,n7,n6,n3,n0]",
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
				Priority:   1,
				Period:     100,
				Deadline:   80,
				Jitter:     10,
				PacketSize: 8,
				Route:      "[n0,n1,n2]",
			},
			{
				ID:         "t2",
				Priority:   2,
				Period:     120,
				Deadline:   115,
				Jitter:     3,
				PacketSize: 10,
				Route:      "[n1,n2,n5,n8]",
			},
			{
				ID:         "t3",
				Priority:   3,
				Period:     150,
				Deadline:   135,
				Jitter:     6,
				PacketSize: 16,
				Route:      "[n3,n4,n7]",
			},
			{
				ID:         "t4",
				Priority:   4,
				Period:     100,
				Deadline:   90,
				Jitter:     5,
				PacketSize: 8,
				Route:      "[n6,n7,n8,n5]",
			},
			{
				ID:         "t5",
				Priority:   5,
				Period:     120,
				Deadline:   100,
				Jitter:     1,
				PacketSize: 2,
				Route:      "[n2,n1,n0,n3]",
			},
			{
				ID:         "t6",
				Priority:   6,
				Period:     200,
				Deadline:   180,
				Jitter:     8,
				PacketSize: 24,
				Route:      "[n4,n3]",
			},
			{
				ID:         "t7",
				Priority:   7,
				Period:     70,
				Deadline:   50,
				Jitter:     2,
				PacketSize: 4,
				Route:      "[n6,n3,n0]",
			},
			{
				ID:         "t8",
				Priority:   8,
				Period:     50,
				Deadline:   45,
				Jitter:     5,
				PacketSize: 8,
				Route:      "[n7,n8,n5,n2]",
			},
			{
				ID:         "t9",
				Priority:   9,
				Period:     100,
				Deadline:   70,
				Jitter:     9,
				PacketSize: 12,
				Route:      "[n4,n5,n2]",
			},
			{
				ID:         "t10",
				Priority:   10,
				Period:     120,
				Deadline:   90,
				Jitter:     14,
				PacketSize: 10,
				Route:      "[n8,n7,n6,n3,n0]",
			},
			{
				ID:         "t11",
				Priority:   11,
				Period:     100,
				Deadline:   85,
				Jitter:     7,
				PacketSize: 16,
				Route:      "[n5,n4,n1]",
			},
			{
				ID:         "t12",
				Priority:   12,
				Period:     120,
				Deadline:   115,
				Jitter:     4,
				PacketSize: 8,
				Route:      "[n7,n4]",
			},
			{
				ID:         "t13",
				Priority:   13,
				Period:     100,
				Deadline:   95,
				Jitter:     10,
				PacketSize: 12,
				Route:      "[n8,n7,n6]",
			},
			{
				ID:         "t14",
				Priority:   14,
				Period:     75,
				Deadline:   70,
				Jitter:     13,
				PacketSize: 10,
				Route:      "[n5,n4,n7]",
			},
			{
				ID:         "t15",
				Priority:   15,
				Period:     100,
				Deadline:   85,
				Jitter:     12,
				PacketSize: 16,
				Route:      "[n3,n4,n5]",
			},
			{
				ID:         "t16",
				Priority:   16,
				Period:     120,
				Deadline:   100,
				Jitter:     11,
				PacketSize: 8,
				Route:      "[n1,n0,n3,n6]",
			},
			{
				ID:         "t17",
				Priority:   17,
				Period:     100,
				Deadline:   100,
				Jitter:     15,
				PacketSize: 12,
				Route:      "[n0,n1]",
			},
			{
				ID:         "t18",
				Priority:   18,
				Period:     115,
				Deadline:   100,
				Jitter:     16,
				PacketSize: 10,
				Route:      "[n2,n5,n8]",
			},
			{
				ID:         "t19",
				Priority:   19,
				Period:     110,
				Deadline:   100,
				Jitter:     17,
				PacketSize: 16,
				Route:      "[n4,n5]",
			},
			{
				ID:         "t20",
				Priority:   20,
				Period:     110,
				Deadline:   95,
				Jitter:     18,
				PacketSize: 8,
				Route:      "[n7,n6,n3]",
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
			network, err := network.NewNetwork(testCase.topologyFunc(b), testCase.networkConf, zerolog.New(io.Discard).With().Logger())
			require.NoError(b, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i], testCase.networkConf)
				require.NoError(b, err)

				trafficFlows[i] = tf
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := newSimulator(network, trafficFlows, testCase.cycles, zerolog.New(io.Discard))
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
	}

	// Restrict cycle limit for expected packets
	for name, testCase := range testCases {
		testCase.cycles = 20
		testCases[name] = testCase
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			if !testCase.run || !testCase.templateRun {
				t.Skip()
			}

			t.Parallel()

			// Setup & Run Simulation
			network, err := network.NewNetwork(testCase.topologyFunc(t), testCase.networkConf, zerolog.New(io.Discard).With().Logger())
			require.NoError(t, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i], testCase.networkConf)
				require.NoError(t, err)

				trafficFlows[i] = tf
			}

			simulator, err := newSimulator(network, trafficFlows, testCase.cycles, zerolog.New(io.Discard))
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
			network, err := network.NewNetwork(testCase.topologyFunc(b), testCase.networkConf, zerolog.New(io.Discard).With().Logger())
			require.NoError(b, err)

			var trafficFlows []traffic.TrafficFlow = make([]traffic.TrafficFlow, len(testCase.traffic))
			for i := 0; i < len(testCase.traffic); i++ {
				tf, err := traffic.NewTrafficFlow(testCase.traffic[i], testCase.networkConf)
				require.NoError(b, err)

				trafficFlows[i] = tf
			}

			simulator, err := newSimulator(network, trafficFlows, testCase.cycles, zerolog.New(io.Discard))
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
