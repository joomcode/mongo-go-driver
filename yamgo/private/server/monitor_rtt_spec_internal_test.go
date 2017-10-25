package server

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"fmt"

	"github.com/10gen/mongo-go-driver/yamgo/internal/testutil/helpers"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	AvgRttMs  json.Number `json:"avg_rtt_ms"`
	NewRttMs  float64     `json:"new_rtt_ms"`
	NewAvgRtt float64     `json:"new_avg_rtt"`
}

const testsDir string = "../../../data/server-selection/rtt"

func runTest(t *testing.T, filename string) {
	filepath := path.Join(testsDir, filename)
	content, err := ioutil.ReadFile(filepath)
	require.NoError(t, err)

	// Remove ".json" from filename.
	testName := filename[:len(filename)-5]

	t.Run(testName, func(t *testing.T) {
		var test testCase
		require.NoError(t, json.Unmarshal(content, &test))
		fmt.Println(string(content))
		fmt.Println(test.AvgRttMs)

		var monitor Monitor

		if test.AvgRttMs != "NULL" {
			avg, err := test.AvgRttMs.Float64()
			require.NoError(t, err)

			monitor.averageRTT = time.Duration(avg * float64(time.Millisecond))
			monitor.averageRTTSet = true
		}

		monitor.updateAverageRTT(time.Duration(test.NewRttMs * float64(time.Millisecond)))
		require.Equal(t, monitor.averageRTT, time.Duration(test.NewAvgRtt*float64(time.Millisecond)))
	})
}

// Test case for all server selection rtt spec tests.
func TestServerSelectionRTTSpec(t *testing.T) {
	for _, file := range testhelpers.FindJSONFilesInDir(t, testsDir) {
		runTest(t, file)
	}
}
