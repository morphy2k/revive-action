// The Revive action binary
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	revivePath = "/revive"
	formatter  = "ndjson"
)

var version = ""

func runRevive(args []string) (*statistics, int, error) {
	cmd := exec.Command(revivePath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting stdout pipe: %w", err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return nil, 0, fmt.Errorf("error while running revive: %w", err)
	}

	dec := json.NewDecoder(stdout)

	stats := &statistics{}

	fmt.Println("::group::Failures")

	for dec.More() {
		f := &failure{}
		if err := dec.Decode(f); err != nil {
			fmt.Println("::endgroup::")
			return nil, 0, fmt.Errorf("error while decoding revive output: %w", err)
		}

		stats.Total++

		switch f.Severity {
		case "warning":
			stats.Warnings++
		case "error":
			stats.Errors++
		}

		fmt.Println(f.Format())
	}

	fmt.Println("::endgroup::")

	var exitErr *exec.ExitError
	if err := cmd.Wait(); err != nil && !errors.As(err, &exitErr) {
		return nil, 0, fmt.Errorf("error while waiting for revive: %w", err)
	}

	code := cmd.ProcessState.ExitCode()

	return stats, code, nil
}

func getReviveVersion() (string, error) {
	cmd := exec.Command(revivePath, "-version")

	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error while getting revive version: %w", err)
	}

	output := strings.TrimSpace(string(stdout))
	parts := strings.Fields(output)
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output format: %s", output)
	}

	version := parts[1]

	return version, nil
}

func buildArgs(input *input) []string {
	args := []string{"-formatter", formatter}

	if input.config != "" {
		args = append(args, "-config", input.config)
	}

	for _, path := range input.exclude {
		args = append(args, "-exclude", path)
	}

	if input.failOnAny {
		args = append(args, "-set_exit_status")
	}

	args = append(args, input.path)

	return args
}

func writeSummary(stats *statistics) error {
	path, ok := os.LookupEnv("GITHUB_STEP_SUMMARY")
	if !ok {
		return fmt.Errorf("GITHUB_STEP_SUMMARY not set")
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open summary file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(stats.ToMarkdownTable()); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	return nil
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "::error %s", err)
	os.Exit(1)
}

func main() {
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	input, err := parseInput()
	if err != nil {
		handleError(err)
	}

	args := buildArgs(input)

	reviveVersion, err := getReviveVersion()
	if err != nil {
		handleError(err)
	}

	fmt.Printf("ACTION: %s\nREVIVE: %s\n", version, reviveVersion)

	stats, code, err := runRevive(args)
	if err != nil {
		handleError(err)
	}

	fmt.Println("Successful run with", stats.String())

	if v := os.Getenv("GITHUB_ACTIONS"); v == "true" {
		if err := writeSummary(stats); err != nil {
			handleError(err)
		}
	} else {
		fmt.Println("Running outside of GitHub Actions, skipping summary")
	}

	os.Exit(code)
}
