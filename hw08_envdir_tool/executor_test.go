package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Set env without data", func(t *testing.T) {
		cmd := []string{"whoami"}

		code := RunCmd(cmd, Environment{})

		require.Equal(t, success, code)
	})

	t.Run("Set env value", func(t *testing.T) {
		cmd := []string{"whoami", "--help"}

		code := RunCmd(cmd, Environment{"val": EnvValue{Value: "value"}})

		require.Equal(t, success, code)
		require.Contains(t, os.Environ(), "val=value")
	})

	t.Run("Set env more values", func(t *testing.T) {
		cmd := []string{"whoami", "--version"}
		env := Environment{
			"val":  EnvValue{"value", true},
			"val1": EnvValue{"value1", false},
			"val2": EnvValue{"value2", false},
			"val3": EnvValue{"value3", false},
		}

		code := RunCmd(cmd, env)

		require.Equal(t, success, code)
		require.NotContains(t, os.Environ(), "val=value")
		require.Contains(t, os.Environ(), "val1=value1")
		require.Contains(t, os.Environ(), "val2=value2")
		require.Contains(t, os.Environ(), "val3=value3")
	})

	t.Run("Error: failure command", func(t *testing.T) {
		cmd := []string{"command"}

		code := RunCmd(cmd, Environment{})

		require.Equal(t, failure, code)
	})

	t.Run("Error: failure argument", func(t *testing.T) {
		cmd := []string{"whoami", "--wrong"}

		code := RunCmd(cmd, Environment{})

		require.Equal(t, failure, code)
	})
}
