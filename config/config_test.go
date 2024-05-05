package config

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"main/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	t.Parallel()

	const testResourcesPath = "test_resources"

	type testCase struct {
		filename string
		err      error
		conf     domain.SimConfig
	}

	testCases := []testCase{
		{
			filename: "valid_basic.yaml",
			err:      nil,
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
	}

	tmpDir := t.TempDir()

	for _, tc := range testCases {
		t.Run(strings.TrimSuffix(tc.filename, filepath.Ext(tc.filename)), func(t *testing.T) {
			yPath := path.Join(testResourcesPath, tc.filename)

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
