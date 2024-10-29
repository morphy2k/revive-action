package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *input
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			want: &input{
				exclude: []string{},
				config:  "",
				path:    defaultPath,
			},
		},
		{
			name: "exclude with single value",
			envVars: map[string]string{
				"INPUT_EXCLUDE": "test",
			},
			want: &input{
				exclude: []string{"test"},
				config:  "",
				path:    defaultPath,
			},
		},
		{
			name: "exclude with multiple values",
			envVars: map[string]string{
				"INPUT_EXCLUDE": "test1;test2;test3",
			},
			want: &input{
				exclude: []string{"test1", "test2", "test3"},
				config:  "",
				path:    defaultPath,
			},
		},
		{
			name: "config with value",
			envVars: map[string]string{
				"INPUT_CONFIG": "config.yaml",
			},
			want: &input{
				exclude: []string{},
				config:  "config.yaml",
				path:    defaultPath,
			},
		},
		{
			name: "custom path",
			envVars: map[string]string{
				"INPUT_PATH": "./custom",
			},
			want: &input{
				exclude: []string{},
				config:  "",
				path:    "./custom",
			},
		},
		{
			name: "all values set",
			envVars: map[string]string{
				"INPUT_EXCLUDE": "test1;test2",
				"INPUT_CONFIG":  "config.yaml",
				"INPUT_PATH":    "./custom",
			},
			want: &input{
				exclude: []string{"test1", "test2"},
				config:  "config.yaml",
				path:    "./custom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer os.Clearenv()

			got := parseInput()
			assert.Equal(t, tt.want, got)
		})
	}
}
