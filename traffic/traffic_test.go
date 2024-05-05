package traffic

import (
	"testing"

	"main/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTrafficFlow(t *testing.T) {
	t.Parallel()

	t.Run("ImplementsInterface", func(t *testing.T) {
		trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
			Src:        "n1",
			Dst:        "n2",
			Priority:   1,
			Period:     75,
			Deadline:   50,
			Jitter:     10,
			PacketSize: 32,
		})
		require.NoError(t, err)
		assert.Implements(t, (*TrafficFlow)(nil), trafficFlow)
	})

	t.Run("Valid", func(t *testing.T) {
		config := domain.TrafficFlowConfig{
			ID:         "t1",
			Src:        "n1",
			Dst:        "n2",
			Priority:   5,
			Period:     4,
			Deadline:   3,
			Jitter:     2,
			PacketSize: 1,
		}

		trafficFlow, err := NewTrafficFlow(config)
		require.NoError(t, err)
		assert.Equal(t, config.ID, trafficFlow.id)
		assert.Equal(t, config.Src, trafficFlow.src)
		assert.Equal(t, config.Dst, trafficFlow.dst)
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
		Src:        "n1",
		Dst:        "n2",
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, id, trafficFlow.ID())
}

func TestTrafficFlowSrc(t *testing.T) {
	t.Parallel()

	src := "n1"

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        src,
		Dst:        "n2",
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, src, trafficFlow.Src())
}

func TestTrafficFlowDst(t *testing.T) {
	t.Parallel()

	dst := "n2"

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        dst,
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, dst, trafficFlow.Dst())
}

func TestTrafficFlowPriority(t *testing.T) {
	t.Parallel()

	priority := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        "n2",
		Priority:   priority,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, priority, trafficFlow.Priority())
}

func TestTrafficFlowReleasePeriod(t *testing.T) {
	t.Parallel()

	releasePeriod := 100

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        "n2",
		Priority:   1,
		Period:     releasePeriod,
		Deadline:   50,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, releasePeriod, trafficFlow.ReleasePeriod())
}

func TestTrafficFlowDeadline(t *testing.T) {
	t.Parallel()

	deadline := 50

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        "n2",
		Priority:   1,
		Period:     75,
		Deadline:   deadline,
		Jitter:     10,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, deadline, trafficFlow.Deadline())
}

func TestTrafficFlowJitter(t *testing.T) {
	t.Parallel()

	jitter := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        "n2",
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     jitter,
		PacketSize: 32,
	})
	require.NoError(t, err)
	assert.Equal(t, jitter, trafficFlow.Jitter())
}

func TestTrafficFlowPacketSize(t *testing.T) {
	t.Parallel()

	packetSize := 1

	trafficFlow, err := NewTrafficFlow(domain.TrafficFlowConfig{
		Src:        "n1",
		Dst:        "n2",
		Priority:   1,
		Period:     75,
		Deadline:   50,
		Jitter:     10,
		PacketSize: packetSize,
	})
	require.NoError(t, err)
	assert.Equal(t, packetSize, trafficFlow.PacketSize())
}

func TestValidateAgainstConfig(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		tfConf := domain.TrafficFlowConfig{
			ID:         "t1",
			Src:        "n1",
			Dst:        "n2",
			Priority:   5,
			Period:     4,
			Deadline:   3,
			Jitter:     2,
			PacketSize: 1,
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
			Src:        "n1",
			Dst:        "n2",
			Priority:   priority,
			Period:     75,
			Deadline:   50,
			Jitter:     10,
			PacketSize: 32,
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
			Src:        "n1",
			Dst:        "n2",
			Priority:   1,
			Period:     releasePeriod,
			Deadline:   50,
			Jitter:     jitter,
			PacketSize: 32,
		})
		require.NoError(t, err)

		var cycle int = 0

		released, periodCycle := trafficFlow.ReleasePacket(cycle)
		assert.True(t, released)
		assert.Equal(t, cycle, periodCycle)

		cycle = 1
		released, _ = trafficFlow.ReleasePacket(cycle)
		assert.False(t, released)

		cycle = 75
		released, periodCycle = trafficFlow.ReleasePacket(cycle)
		assert.True(t, released)
		assert.Equal(t, cycle, periodCycle)
	})
}
