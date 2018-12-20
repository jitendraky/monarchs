package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	overriddenConfig = &Config{
		ListenAddress: "255.255.255.255",
		ListenPort:    9999,
		LogFormat:     "json",
		LogLevel:      "error",
		LogOutput:     "stdout",
	}
)

func TestParseFlags(t *testing.T) {
	for _, ti := range []struct {
		title    string
		args     []string
		envVars  map[string]string
		expected *Config
	}{
		{
			title:    "default config with minimal flags",
			args:     []string{},
			envVars:  map[string]string{},
			expected: defaultConfig,
		},
		{
			title: "override everything with flags",
			args: []string{
				"--addr=255.255.255.255",
				"--port=9999",
				"--log-format=json",
				"--log-level=error",
				"--log-output=stdout",
			},
			envVars:  map[string]string{},
			expected: overriddenConfig,
		},
		{
			title: "override everything with env vars",
			args:  []string{},
			envVars: map[string]string{
				"MONARCHS_ADDR":       "255.255.255.255",
				"MONARCHS_PORT":       "9999",
				"MONARCHS_LOG_FORMAT": "json",
				"MONARCHS_LOG_LEVEL":  "error",
				"MONARCHS_LOG_OUTPUT": "stdout",
			},
			expected: overriddenConfig,
		},
	} {
		t.Run(ti.title, func(t *testing.T) {
			env0 := setEnv(t, ti.envVars)
			defer func() { restoreEnv(t, env0) }()

			cfg := NewConfig()
			require.NoError(t, cfg.ParseFlags(ti.args))
			assert.Equal(t, ti.expected, cfg)
		})
	}
}

func setEnv(t *testing.T, env map[string]string) map[string]string {
	env2 := map[string]string{}

	for k, v := range env {
		env2[k] = os.Getenv(k)
		require.NoError(t, os.Setenv(k, v))
	}

	return env2
}

func restoreEnv(t *testing.T, env map[string]string) {
	for k, v := range env {
		require.NoError(t, os.Setenv(k, v))
	}
}