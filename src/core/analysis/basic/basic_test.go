package basic

import (
	"strconv"
	"testing"

	"main/src/core/analysis/util"
	"main/src/domain"
	"main/src/topology"

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
				ProcessingDelay: 6,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				PacketSize: 1,
				Route:      "[n0,n1,n2]",
			},
			expected: 19,
		},
		{
			conf: domain.SimConfig{
				ProcessingDelay: 3,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				PacketSize: 2,
				Route:      "[n0,n1,n2]",
			},
			expected: 11,
		},
		{
			conf: domain.SimConfig{
				ProcessingDelay: 6,
			},
			top: topology.ThreeByThreeMesh,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				PacketSize: 8,
				Route:      "[n0,n1,n2]",
			},
			expected: 26,
		},
		{
			conf: domain.SimConfig{
				ProcessingDelay: 6,
			},
			top: topology.ThreeByThreeMesh,
			tf: domain.TrafficFlowConfig{
				ID:         "t5",
				PacketSize: 13,
				Route:      "[n3,n4,n1]",
			},
			expected: 31,
		},
		{
			conf: domain.SimConfig{
				ProcessingDelay: 6,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				PacketSize: 4,
				Route:      "[n0,n1,n2]",
			},
			expected: 22,
		},
		{
			conf: domain.SimConfig{
				ProcessingDelay: 3,
			},
			top: topology.ThreeHorozontalLine,
			tf: domain.TrafficFlowConfig{
				ID:         "t1",
				PacketSize: 12,
				Route:      "[n0,n1,n2]",
			},
			expected: 21,
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
