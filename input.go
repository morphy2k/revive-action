package main

import (
	"os"
	"strings"
)

const (
	defaultPath = "./..."
)

type input struct {
	exclude []string
	config  string
	path    string
}

func parseInput() *input {
	input := &input{
		exclude: make([]string, 0),
		path:    defaultPath,
	}

	if v, ok := os.LookupEnv("INPUT_EXCLUDE"); ok {
		switch {
		case strings.Contains(v, ";"):
			input.exclude = strings.Split(v, ";")
		default:
			input.exclude = strings.Split(v, "\n")
		}
	}

	if v, ok := os.LookupEnv("INPUT_CONFIG"); ok {
		input.config = v
	}

	if v, ok := os.LookupEnv("INPUT_PATH"); ok {
		input.path = v
	}

	return input
}
