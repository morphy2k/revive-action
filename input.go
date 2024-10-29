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
		input.exclude = strings.Split(v, ";")
	}

	if v, ok := os.LookupEnv("INPUT_CONFIG"); ok {
		input.config = v
	}

	if v, ok := os.LookupEnv("INPUT_PATH"); ok {
		input.path = v
	}

	return input
}
