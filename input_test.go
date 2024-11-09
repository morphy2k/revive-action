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
		wantErr  bool
	}{
		{
			name: "default values",
			env:  map[string]string{},
			expected: &input{
				exclude:   make([]string, 0),
				path:      defaultPath,
				config:    "",
				failOnAny: false,
			},
		},
		{
			name: "with semicolon separated excludes",
			env: map[string]string{
				envExclude: "dir1;./dir2/...;dir3/dir4/...",
			},
			expected: &input{
				exclude:   []string{"dir1", "./dir2/...", "dir3/dir4/..."},
				path:      defaultPath,
				config:    "",
				failOnAny: false,
			},
		},
		{
			name: "with newline separated excludes",
			env: map[string]string{
				envExclude: "dir1\n./dir2/...\ndir3/dir4/...",
			},
			expected: &input{
				exclude:   []string{"dir1", "./dir2/...", "dir3/dir4/..."},
				path:      defaultPath,
				config:    "",
				failOnAny: false,
			},
		},
		{
			name: "with custom config",
			env: map[string]string{
				envConfig: "custom.toml",
			},
			expected: &input{
				exclude:   make([]string, 0),
				path:      defaultPath,
				config:    "custom.toml",
				failOnAny: false,
			},
		},
		{
			name: "with custom path",
			env: map[string]string{
				envPath: "./src/...",
			},
			expected: &input{
				exclude:   make([]string, 0),
				path:      "./src/...",
				config:    "",
				failOnAny: false,
			},
		},
		{
			name: "with fail on any true",
			env: map[string]string{
				envFailOnAny: "true",
			},
			expected: &input{
				exclude:   make([]string, 0),
				path:      defaultPath,
				config:    "",
				failOnAny: true,
			},
		},
		{
			name: "with invalid fail on any",
			env: map[string]string{
				envFailOnAny: "invalid",
			},
			wantErr: true,
		},
		{
			name: "with all options",
			env: map[string]string{
				envExclude:   "dir1;dir2",
				envConfig:    "custom.toml",
				envPath:      "./src/...",
				envFailOnAny: "true",
			},
			expected: &input{
				exclude:   []string{"dir1", "dir2"},
				path:      "./src/...",
				config:    "custom.toml",
				failOnAny: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			got, err := parseInput()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
