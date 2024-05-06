package basic

import (
	"strconv"
	"testing"

	"main/core/analysis/util"
	"main/domain"
	"main/topology"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicLatency(t *testing.T) {
	t.Parallel()

	type testCase struct {
		conf     domain.SimConfig
		top      func(t testing.TB) *topology.Topology
		tf       domain.TrafficFlowConfig
		expected int
	}

	testCases := []testCase{
		{
			conf: domain.SimConfig{
				RoutingAlgorithm: "XY",
				FlitSize:         4,
				ProcessingDelay:  6,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				PacketSize: 4,
			},
			expected: 21,
		},
		{
			conf: domain.SimConfig{
				RoutingAlgorithm: "XY",
				FlitSize:         1,
				ProcessingDelay:  3,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				PacketSize: 2,
			},
			expected: 13,
		},
		{
			conf: domain.SimConfig{
				RoutingAlgorithm: "XY",
				FlitSize:         4,
				ProcessingDelay:  6,
			},
			top: topology.ThreeByThreeMesh,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				Src:        "n0",
				Dst:        "n2",
				PacketSize: 32,
			},
			expected: 28,
		},
		{
			conf: domain.SimConfig{
				RoutingAlgorithm: "XY",
				FlitSize:         4,
				ProcessingDelay:  6,
			},
			top: topology.ThreeByThreeMesh,
			tf: domain.TrafficFlowConfig{
				ID:         "t5",
				Src:        "n3",
				Dst:        "n1",
				PacketSize: 50,
			},
			expected: 33,
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			top := tc.top(t)

			tfr, err := util.NewTrafficFlowAndRoute(top, tc.tf)
			require.NoError(t, err)

			lat := BasicLatency(tc.conf, tfr)

			assert.Equal(t, tc.expected, lat)
		})
	}
}
