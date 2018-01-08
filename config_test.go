// Author hoenig

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	homeEnv = "HOME"
)

func Test_defaultConfigfileLocations(t *testing.T) {
	oldHOME := os.Getenv(homeEnv)
	defer func() {
		os.Setenv(homeEnv, oldHOME)
	}()

	os.Setenv(homeEnv, "/path/to/nowhere")

	locations := defaultConfigfileLocations()
	require.Equal(t, []string{
		"/path/to/nowhere/.config/marathonctl/config",
		"/etc/marathonctl",
	}, locations)
}
