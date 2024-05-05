package config

import (
	"fmt"
	"os"
	"testing"

	"main/domain"

	"github.com/stretchr/testify/require"
)

var testConfig domain.SimConfig = domain.SimConfig{
	CycleLimit:       10000,
	RoutingAlgorithm: "XY",
	BufferSize:       2,
	FlitSize:         4,
	ProcessingDelay:  6,
}

func TestYaml(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		testYaml := fmt.Sprintf(`
cycle_limit: %d
buffer_size: %d
flit_size: %d
routing_algorithm: %s
processing_delay: %d`,
			testConfig.CycleLimit,
			testConfig.BufferSize,
			testConfig.FlitSize,
			testConfig.RoutingAlgorithm,
			testConfig.ProcessingDelay,
		)

		tmpFile := t.TempDir() + "/config.yaml"
		err := os.WriteFile(tmpFile, []byte(testYaml), 0o644)
		require.NoError(t, err)

		conf, err := readYaml(tmpFile)
		require.NoError(t, err)
		require.Equal(t, testConfig, conf)
	})
}

func TestJson(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		testJson := fmt.Sprintf(`
{
	"cycle_limit": %d,
	"buffer_size": %d,
	"flit_size": %d,
	"routing_algorithm": "%s",
	"processing_delay": %d
}`,
			testConfig.CycleLimit,
			testConfig.BufferSize,
			testConfig.FlitSize,
			testConfig.RoutingAlgorithm,
			testConfig.ProcessingDelay,
		)

		tmpFile := t.TempDir() + "/config.json"
		err := os.WriteFile(tmpFile, []byte(testJson), 0o644)
		require.NoError(t, err)

		conf, err := readJson(tmpFile)
		require.NoError(t, err)
		require.Equal(t, testConfig, conf)
	})
}
