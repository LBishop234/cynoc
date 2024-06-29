package traffic

import (
	"io"
	"testing"

	"main/src/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTrafficFlow(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
			Priority:   1,
			Period:     75,
			Deadline:   50,
			Jitter:     10,
			PacketSize: 32,
			Route:      "[n1,n2,n3]",
		})
		require.NoError(t, err)
		assert.Implements(t, (*TrafficFlow)(nil), trafficFlow)
	})

	t.Run("Valid", func(t *testing.T) {
		config := domain.TrafficFlowConfig{
			ID:         "t1",
			Priority:   5,
			Period:     4,
			Deadline:   3,
			Jitter:     2,
			PacketSize: 1,
			Route:      "[n1,n2,n3]",
		}

		trafficFlow, err := NewTrafficFlow(config)
		require.NoError(t, err)
		assert.Equal(t, config.ID, trafficFlow.id)
		assert.Equal(t, config.Priority, trafficFlow.priority)
		assert.Equal(t, config.Period, trafficFlow.releasePeriod)
		assert.Equal(t, config.Deadline, trafficFlow.deadline)
		assert.Equal(t, config.Jitter, trafficFlow.jitter)
		assert.Equal(t, config.PacketSize, trafficFlow.packetSize)
	})
}

func TestTrafficFlowID(t *testing.T) {
	t.Parallel()

	id := "t1"

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		ID:         id,
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, id, trafficFlow.ID())
}

func TestTrafficFlowPriority(t *testing.T) {
	t.Parallel()

	priority := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   priority,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, priority, trafficFlow.Priority())
}

func TestTrafficFlowReleasePeriod(t *testing.T) {
	t.Parallel()

	releasePeriod := 100

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   1,
		Period:     releasePeriod,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, releasePeriod, trafficFlow.ReleasePeriod())
}

func TestTrafficFlowDeadline(t *testing.T) {
	t.Parallel()

	deadline := 50

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   1,
		Period:     75,
		Deadline:   deadline,
		Jitter:     10,
		PacketSize: 32,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, deadline, trafficFlow.Deadline())
}

func TestTrafficFlowJitter(t *testing.T) {
	t.Parallel()

	jitter := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     jitter,
		PacketSize: 32,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, jitter, trafficFlow.Jitter())
}

func TestTrafficFlowPacketSize(t *testing.T) {
	t.Parallel()

	packetSize := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: packetSize,
		Route:      "[n1,n2,n3]",
	})
	require.NoError(t, err)
	assert.Equal(t, packetSize, trafficFlow.PacketSize())
}

func TestTrafficFlowRoute(t *testing.T) {
	t.Parallel()

	routeStr := "[n1,n2,n3]"
	route := []string{"n1", "n2", "n3"}

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
		Route:      routeStr,
	})
	require.NoError(t, err)
	assert.Equal(t, route, trafficFlow.Route())
}

func TestValidateAgainstConfig(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		tfConf := domain.TrafficFlowConfig{
			ID:         "t1",
			Priority:   5,
			Period:     4,
			Deadline:   3,
			Jitter:     2,
			PacketSize: 1,
			Route:      "[n1,n2,n3]",
		}
		trafficFlow, err := NewTrafficFlow(tfConf)
		require.NoError(t, err)

		config := domain.SimConfig{
			MaxPriority: 5,
		}

		err = trafficFlow.ValidateAgainstConfig(config)
		require.NoError(t, err)
	})

	t.Run("InvalidPriority", func(t *testing.T) {
		var priority int = 3

		tfConf := domain.TrafficFlowConfig{
			Priority:   priority,
			Period:     75,
			Deadline:   50,
			Jitter:     10,
			PacketSize: 32,
			Route:      "[n1,n2,n3]",
		}
		trafficFlow, err := NewTrafficFlow(tfConf)
		require.NoError(t, err)

		config := domain.SimConfig{
			MaxPriority: 2,
		}

		err = trafficFlow.ValidateAgainstConfig(config)
		require.ErrorIs(t, err, domain.ErrInvalidConfig)
	})
}

func TestTrafficFlowRleasePacket(t *testing.T) {
	t.Parallel()

	t.Run("ValidNoJitter", func(t *testing.T) {
		const releasePeriod = 75
		const jitter = 0

		trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
			Priority:   1,
			Period:     releasePeriod,
			Deadline:   50,
			Jitter:     jitter,
			PacketSize: 32,
			Route:      "[n1,n2,n3]",
		})
		require.NoError(t, err)

		var cycle int = 0

		released, _, periodCycle := trafficFlow.ReleasePacket(cycle, trafficFlow, domain.Route{}, zerolog.New(io.Discard))
		assert.True(t, released)
		assert.Equal(t, cycle, periodCycle)

		cycle = 1
		released, _, _ = trafficFlow.ReleasePacket(cycle, trafficFlow, domain.Route{}, zerolog.New(io.Discard))
		assert.False(t, released)

		cycle = 75
		released, _, periodCycle = trafficFlow.ReleasePacket(cycle, trafficFlow, domain.Route{}, zerolog.New(io.Discard))
		assert.True(t, released)
		assert.Equal(t, cycle, periodCycle)
	})
}
