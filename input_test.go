package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *input
	}{
		{
			name: "default values",
			env:  map[string]string{},
			expected: &input{
				exclude: make([]string, 0),
				path:    defaultPath,
				config:  "",
			},
		},
		{
			name: "with semicolon separated excludes",
			env: map[string]string{
				"INPUT_EXCLUDE": "dir1;./dir2/...;dir3/dir4/...",
			},
			expected: &input{
				exclude: []string{"dir1", "./dir2/...", "dir3/dir4/..."},
				path:    defaultPath,
				config:  "",
			},
		},
		{
			name: "with newline separated excludes",
			env: map[string]string{
				"INPUT_EXCLUDE": "dir1\n./dir2/...\ndir3/dir4/...",
			},
			expected: &input{
				exclude: []string{"dir1", "./dir2/...", "dir3/dir4/..."},
				path:    defaultPath,
				config:  "",
			},
		},
		{
			name: "with custom config",
			env: map[string]string{
				"INPUT_CONFIG": "custom.toml",
			},
			expected: &input{
				exclude: make([]string, 0),
				path:    defaultPath,
				config:  "custom.toml",
			},
		},
		{
			name: "with custom path",
			env: map[string]string{
				"INPUT_PATH": "./src/...",
			},
			expected: &input{
				exclude: make([]string, 0),
				path:    "./src/...",
				config:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set up test environment
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			got := parseInput()
			assert.Equal(t, tt.expected, got)
		})
	}
}
