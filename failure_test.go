package main

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		failure  failure
		expected string
	}{
		{
			name: "error severity",
			failure: failure{
				Position: struct {
					Start token.Position
					End   token.Position
				}{
					Start: token.Position{Filename: "main.go", Line: 10, Column: 5},
					End:   token.Position{Line: 10, Column: 10},
				},
				Failure:  "some error",
				Severity: "error",
			},
			expected: "::error file=main.go,line=10,endLine=10,col=5,endColumn=10::some error",
		},
		{
			name: "warning severity",
			failure: failure{
				Position: struct {
					Start token.Position
					End   token.Position
				}{
					Start: token.Position{Filename: "main.go", Line: 20, Column: 15},
					End:   token.Position{Line: 20, Column: 20},
				},
				Failure:  "some warning",
				Severity: "warning",
			},
			expected: "::warning file=main.go,line=20,endLine=20,col=15,endColumn=20::some warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.failure.Format()
			assert.Equal(t, tt.expected, result)
		})
	}
}
