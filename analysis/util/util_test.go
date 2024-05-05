package util

import (
	"testing"

	"main/domain"
	"main/topology"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstructTrafficFlowAndRoutes(t *testing.T) {
	t.Skip("Skipping to prevent replication of TestNewTrafficFlowAndRoute")
}

func TestNewTrafficFlowAndRoute(t *testing.T) {
	t.Parallel()

	top := topology.ThreeByThreeMesh(t)

	tf1 := domain.TrafficFlowConfig{
		ID:  "tf1",
		Src: "n3",
		Dst: "n5",
	}

	expectedRoute := domain.Route{
		top.Nodes()["n3"].NodeID(),
		top.Nodes()["n4"].NodeID(),
		top.Nodes()["n5"].NodeID(),
	}

	tfr, err := NewTrafficFlowAndRoute(top, tf1)
	require.NoError(t, err)

	assert.Equal(t, tf1, tfr.TrafficFlowConfig)
	assert.Equal(t, expectedRoute, tfr.Route)
}
