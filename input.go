package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// Environmental variable constants
	envExclude   = "INPUT_EXCLUDE"
	envConfig    = "INPUT_CONFIG"
	envPath      = "INPUT_PATH"
	envFailOnAny = "INPUT_FAIL_ON_ANY"

	// Default path
	defaultPath = "./..."
)

type input struct {
	exclude   []string
	config    string
	path      string
	failOnAny bool
}

func parseInput() (*input, error) {
	input := &input{
		exclude: make([]string, 0),
		path:    defaultPath,
		config:  "",
	}

	if err := input.parseExclude(); err != nil {
		return nil, fmt.Errorf("parsing exclude: %w", err)
	}

	if err := input.parseConfig(); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := input.parsePath(); err != nil {
		return nil, fmt.Errorf("parsing path: %w", err)
	}

	if err := input.parseFailOnAny(); err != nil {
		return nil, fmt.Errorf("parsing fail-on-any: %w", err)
	}

	return input, nil
}

func (i *input) parseExclude() error {
	if v, ok := os.LookupEnv(envExclude); ok && v != "" {
		switch {
		case strings.Contains(v, ";"):
			i.exclude = strings.Split(v, ";")
		default:
			i.exclude = strings.Split(v, "\n")
		}
	}

	return nil
}

func (i *input) parseConfig() error {
	if v, ok := os.LookupEnv(envConfig); ok {
		i.config = v
	}

	return nil
}

func (i *input) parsePath() error {
	if v, ok := os.LookupEnv(envPath); ok {
		i.path = v
	}

	return nil
}

func (i *input) parseFailOnAny() error {
	if v, ok := os.LookupEnv(envFailOnAny); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("invalid value: %w", err)
		}
		i.failOnAny = b
	}

	return nil
}
