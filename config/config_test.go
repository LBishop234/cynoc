package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"main/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadConfig(t *testing.T) {
	t.Parallel()

	const testResourcesPath = "test_resources"

	type testCase struct {
		name      string
		baseFile  string
		overrides map[string]any
		err       error
		conf      domain.SimConfig
	}

	testCases := []testCase{
		{
			name:      "valid_basic",
			baseFile:  "valid_basic.yaml",
			overrides: nil,
			err:       nil,
			conf: domain.SimConfig{
				CycleLimit:       1000,
				RoutingAlgorithm: "XY",
				MaxPriority:      6,
				BufferSize:       12,
				FlitSize:         2,
				LinkBandwidth:    6,
				ProcessingDelay:  6,
			},
		},
		{
			name:     "invalid_cycle_limit_zero",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"cycle_limit": 0,
			},
		},
		{
			name:     "invalid_routing_algorithm",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"routing_algorithm": "UNKNOWN",
			},
		},
		{
			name:     "invalid_max_priority_zero",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"max_priority": 0,
			},
		},
		{
			name:     "invalid_flit_size_zero",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"flit_size": 0,
			},
		},
		{
			name:     "invalid_buffer_size_zero",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"buffer_size": 0,
			},
		},
		{
			name:     "invalid_buffer_size_not_multiple_max_priority",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"max_priority": 3,
				"flit_size":    2,
				"buffer_size":  8,
			},
		},
		{
			name:     "invalid_buffer_size_not_multiple_flit_size",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"max_priority": 5,
				"flit_size":    4,
				"buffer_size":  10,
			},
		},
		{
			name:     "invalid_processing_delay_zero",
			err:      domain.ErrInvalidConfig,
			baseFile: "valid_basic.yaml",
			overrides: map[string]any{
				"processing_delay": 0,
			},
		},
		// {
		// 	name:     "invalid_link_bandwidth_zero",
		// 	err:      domain.ErrInvalidConfig,
		// 	baseFile: "valid_basic.yaml",
		// 	overrides: map[string]any{
		// 		"link_bandwidth": 0,
		// 	},
		// },
		// {
		// 	name:     "invalid_link_bandwidth_not_multiple_flit_size",
		// 	err:      domain.ErrInvalidConfig,
		// 	baseFile: "valid_basic.yaml",
		// 	overrides: map[string]any{
		// 		"flit_size":     4,
		// 		"link_bandwidth": 5,
		// 	},
		// },
	}

	tmpDir := t.TempDir()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Copy the base file to a temporary directory.
			basePath := path.Join(testResourcesPath, tc.baseFile)
			yPath := path.Join(tmpDir, fmt.Sprint(tc.name, ".yaml"))

			baseBytes, err := os.ReadFile(basePath)
			require.NoError(t, err)

			var data map[string]any
			err = yaml.Unmarshal(baseBytes, &data)
			require.NoError(t, err)

			if tc.overrides != nil {
				for k, v := range tc.overrides {
					data[k] = v
				}
			}

			yBytes, err := yaml.Marshal(data)
			require.NoError(t, err)

			err = os.WriteFile(yPath, yBytes, 0644)
			require.NoError(t, err)

			t.Run("YAML", func(t *testing.T) {
				conf, err := ReadConfig(yPath)
				require.ErrorIs(t, tc.err, err)

				if err == nil {
					assert.Equal(t, tc.conf, conf)
				}
			})

			t.Run("JSON", func(t *testing.T) {
				jPath := yamlFileToJsonFile(t, tmpDir, yPath)

				conf, err := ReadConfig(jPath)
				require.ErrorIs(t, tc.err, err)

				if err == nil {
					assert.Equal(t, tc.conf, conf)
				}

			})
		})
	}

	//
	// Custom test cases for edge cases.
	//
	t.Run("NoFile", func(t *testing.T) {
		t.Run("YAML", func(t *testing.T) {
			_, err := ReadConfig("no_file.yaml")
			require.Error(t, err)
		})

		t.Run("JSON", func(t *testing.T) {
			_, err := ReadConfig("no_file.json")
			require.Error(t, err)
		})
	})
}

// Creates an equivalent JSON file from a YAML file.
// Returns the path to the temporary JSON file.
func yamlFileToJsonFile(tb testing.TB, tmpDir string, yPath string) string {
	conf, err := readYaml(yPath)
	require.NoError(tb, err)

	jBytes, err := json.Marshal(conf)
	require.NoError(tb, err)

	jFilename := strings.TrimSuffix(filepath.Base(yPath), filepath.Ext(yPath)) + ".json"
	jPath := path.Join(tmpDir, jFilename)

	err = os.WriteFile(jPath, jBytes, 0644)
	require.NoError(tb, err)

	return jPath
}
