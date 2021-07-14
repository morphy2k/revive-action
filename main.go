package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/token"
	"os"
	"sync"
)

var version = "unknown"

type failure struct {
	Failure    string
	RuleName   string
	Category   string
	Position   position
	Confidence float64
	Severity   string
}

type position struct {
	Start token.Position
	End   token.Position
}

type statistics struct {
	Total, Warnings, Errors int
}

func (f statistics) String() string {
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

func printFailure(f *failure, wg *sync.WaitGroup) {
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
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	stats := &statistics{}

	ch := make(chan *failure)
	go getFailures(ch)

	wg := &sync.WaitGroup{}

	fmt.Println("::group::Failures")

	for f := range ch {
		wg.Add(1)

		stats.Total++

		switch f.Severity {
		case "warning":
			stats.Warnings++
		case "error":
			stats.Errors++
		}

		go printFailure(f, wg)
	}

	wg.Wait()

	fmt.Println("::endgroup::")

	fmt.Println("Successful run with", stats.String())
}
