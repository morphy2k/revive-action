package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

const name = "Revive Action"

const (
	envRepo  = "GITHUB_REPOSITORY"
	envSHA   = "GITHUB_SHA"
	envToken = "GITHUB_TOKEN"

	chunkLimit = 50
)

var (
	ghToken   string
	repoOwner string
	repoName  string
	headSHA   string
)

var client *github.Client

func init() {
	if env := os.Getenv(envToken); len(env) > 0 {
		ghToken = env
	} else {
		fmt.Fprintln(os.Stderr, "Missing environment variable:", envToken)
		os.Exit(2)
	}

	if env := os.Getenv(envRepo); len(env) > 0 {
		s := strings.SplitN(env, "/", 2)
		repoOwner, repoName = s[0], s[1]
	} else {
		fmt.Fprintln(os.Stderr, "Missing environment variable:", envRepo)
		os.Exit(2)
	}

	if env := os.Getenv(envSHA); len(env) > 0 {
		headSHA = env
	} else {
		fmt.Fprintln(os.Stderr, "Missing environment variable:", envSHA)
		os.Exit(2)
	}

	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	))

	client = github.NewClient(tc)
}

func createCheck() *github.CheckRun {
	opts := github.CreateCheckRunOptions{
		Name:    name,
		HeadSHA: headSHA,
		Status:  github.String("in_progress"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	check, _, err := client.Checks.CreateCheckRun(ctx, repoOwner, repoName, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while creating check-run:", err)
		os.Exit(1)
	}

	return check
}

type conclusion int

const (
	conclSuccess conclusion = iota
	conclFailure
)

func (c conclusion) String() string {
	return [...]string{"success", "failure"}[c]
}

func completeCheck(check *github.CheckRun, concl conclusion, stats *failureStats) {
	opts := github.UpdateCheckRunOptions{
		Name:       name,
		HeadSHA:    github.String(headSHA),
		Conclusion: github.String(concl.String()),
		Output: &github.CheckRunOutput{
			Title:   github.String("Result"),
			Summary: github.String(stats.String()),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, _, err := client.Checks.UpdateCheckRun(
		ctx, repoOwner, repoName, check.GetID(), opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error while completing check-run:", err)
		os.Exit(1)
	}
}

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

func createAnnotations(failures []*failure) []*github.CheckRunAnnotation {
	ann := make([]*github.CheckRunAnnotation, len(failures))

	for i, f := range failures {
		var level string
		switch f.Severity {
		case "warning":
			level = "warning"
		case "error":
			level = "failure"
		}

		ann[i] = &github.CheckRunAnnotation{
			Path:            github.String(f.Position.Start.Filename),
			StartLine:       github.Int(f.Position.Start.Line),
			EndLine:         github.Int(f.Position.End.Line),
			AnnotationLevel: github.String(level),
			Title: github.String(
				fmt.Sprintf("%s (%s)", strings.Title(f.Category), f.RuleName),
			),
			Message: github.String(f.Failure),
		}
	}

	return ann
}

func pushFailures(check *github.CheckRun, failures []*failure, stats *failureStats, wg *sync.WaitGroup) {
	opts := github.UpdateCheckRunOptions{
		Name:    name,
		HeadSHA: github.String(headSHA),
		Output: &github.CheckRunOutput{
			Title:       github.String("Result"),
			Summary:     github.String(stats.String()),
			Annotations: createAnnotations(failures),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, _, err := client.Checks.UpdateCheckRun(
		ctx, repoOwner, repoName, check.GetID(), opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error while updating check-run:", err)
		os.Exit(1)
	}

	wg.Done()
}

func main() {
	var concl conclusion

	check := createCheck()

	failures := make([]*failure, 0)
	stats := &failureStats{}

	ch := make(chan *failure)
	go getFailures(ch)

	wg := &sync.WaitGroup{}

	chunks := 1
	for f := range ch {
		failures = append(failures, f)

		stats.Total++

		switch f.Severity {
		case "warning":
			stats.Warnings++
		case "error":
			stats.Errors++
		}

		if c := chunks * chunkLimit; stats.Total > c {
			wg.Add(1)
			go pushFailures(check, failures[c-chunkLimit:c], stats, wg)
			chunks++
		}
	}

	if stats.Total > 0 {
		wg.Add(1)
		if chunks == 1 {
			go pushFailures(check, failures, stats, wg)
		} else {
			c := chunks * chunkLimit
			go pushFailures(check, failures[c-chunkLimit:], stats, wg)
		}
		wg.Wait()
		concl = conclFailure
	}

	completeCheck(check, concl, stats)

	fmt.Println("Successful run with", stats.String())
}
