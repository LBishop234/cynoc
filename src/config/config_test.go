package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"main/src/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const testResourcesDir = "test_resources"

func TestReadConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		baseFile  string
		enabled   bool
		err       error
		overrides map[string]any
		expected  domain.SimConfig
	}

	testCases := []testCase{
		{
			name:      "valid_basic",
			baseFile:  "valid_basic.yaml",
			enabled:   true,
			overrides: nil,
			err:       nil,
			expected: domain.SimConfig{
				CycleLimit:      1000,
				MaxPriority:     6,
				BufferSize:      24,
				LinkBandwidth:   2,
				ProcessingDelay: 6,
			},
		},
		{
			name:     "invalid_cycle_limit_zero",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidCycleLimit,
			overrides: map[string]any{
				"cycle_limit": 0,
			},
		},
		{
			name:     "invalid_max_priority_zero",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidMaxPriority,
			overrides: map[string]any{
				"max_priority": 0,
			},
		},
		{
			name:     "invalid_buffer_size_zero",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidBufferSize,
			overrides: map[string]any{
				"buffer_size": 0,
			},
		},
		{
			name:     "invalid_buffer_size_not_multiple_max_priority",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidBufferSize,
			overrides: map[string]any{
				"max_priority": 3,
				"buffer_size":  8,
			},
		},
		{
			name:     "invalid_processing_delay_zero",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidProcessingDelay,
			overrides: map[string]any{
				"processing_delay": 0,
			},
		},
		{
			name:     "invalid_link_bandwidth_zero",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidLinkBandwidth,
			overrides: map[string]any{
				"link_bandwidth": 0,
			},
		},
		{
			name:     "invalid_link_bandwidth_greater_than_half_virtual_channel_size",
			baseFile: "valid_basic.yaml",
			enabled:  true,
			err:      ErrInvalidLinkBandwidth,
			overrides: map[string]any{
				"buffer_size":    12,
				"max_priority":   6,
				"link_bandwidth": 2,
			},
		},
	}

	tmpDir := t.TempDir()

	for _, tc := range testCases {
		if tc.enabled {
			t.Run(tc.name, func(t *testing.T) {
				// Copy the base file to a temporary directory.
				basePath := path.Join(testResourcesDir, tc.baseFile)
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

				err = os.WriteFile(yPath, yBytes, 0o644)
				require.NoError(t, err)

				t.Run("YAML", func(t *testing.T) {
					conf, err := ReadConfig(yPath)
					require.ErrorIs(t, err, tc.err)
					if tc.err != nil {
						require.ErrorIs(t, err, ErrInvalidConfig)
					}

					if err == nil {
						assert.Equal(t, tc.expected, conf)
					}
				})

				t.Run("JSON", func(t *testing.T) {
					jPath := yamlFileToJsonFile(t, tmpDir, yPath)

					conf, err := ReadConfig(jPath)
					require.ErrorIs(t, err, tc.err)
					if tc.err != nil {
						require.ErrorIs(t, err, ErrInvalidConfig)
					}

					if err == nil {
						assert.Equal(t, tc.expected, conf)
					}
				})
			})
		}
	}
}

// Custom test cases for edge cases.
func TestReadConfigEdgeCases(t *testing.T) {
	t.Parallel()

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

	t.Run("InvalidFileExtension", func(t *testing.T) {
		_, err := ReadConfig(path.Join(testResourcesDir, "invalid_file_extension.txt"))
		require.Error(t, err)
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		_, err := ReadConfig(path.Join(testResourcesDir, "invalid_yaml.yaml"))
		require.Error(t, err)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		_, err := ReadConfig(path.Join(testResourcesDir, "invalid_json.json"))
		require.Error(t, err)
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

	err = os.WriteFile(jPath, jBytes, 0o644)
	require.NoError(tb, err)

	return jPath
}
