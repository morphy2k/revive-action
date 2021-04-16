package main

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"sync"
)

type failure struct {
	Failure    string
	RuleName   string
	Category   string
	Position   failurePosition
	Confidence float64
	Severity   string
}

type failurePosition struct {
	Start token.Position
	End   token.Position
}

type failureStats struct {
	Total, Warnings, Errors int
}

func (f failureStats) String() string {
	return fmt.Sprintf("%d failures (%d warnings, %d errors)",
		f.Total, f.Warnings, f.Errors)
}

func getFailures(ch chan *failure) {
	dec := json.NewDecoder(os.Stdin)

	for dec.More() {
		f := &failure{}
		if err := dec.Decode(f); err != nil {
			fmt.Fprintln(os.Stderr, "Error while decoding stdin:", err)
			os.Exit(1)
		}
		ch <- f
	}

	close(ch)
}

func printCommand(f *failure, wg *sync.WaitGroup) {
	s := fmt.Sprintf("file=%s,line=%d,col=%d::%s\n",
		f.Position.Start.Filename, f.Position.Start.Line, f.Position.Start.Column, f.Failure)

	if f.Severity == "warning" {
		fmt.Printf("::warning %s", s)
	} else {
		fmt.Printf("::error %s", s)
	}

	wg.Done()
}

func main() {
	stats := &failureStats{}

	ch := make(chan *failure)
	go getFailures(ch)

	wg := &sync.WaitGroup{}

	for f := range ch {
		wg.Add(1)

		stats.Total++

		switch f.Severity {
		case "warning":
			stats.Warnings++
		case "error":
			stats.Errors++
		}

		go printCommand(f, wg)
	}

	wg.Wait()

	fmt.Println("Successful run with", stats.String())
}
