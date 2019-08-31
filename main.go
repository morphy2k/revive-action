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

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

const (
	envRepo   = "GITHUB_REPOSITORY"
	envAction = "GITHUB_ACTION"
	envSHA    = "GITHUB_SHA"
	envToken  = "GITHUB_TOKEN"

	chunkLimit = 50
)

var (
	ghToken    string
	repoOwner  string
	repoName   string
	headSHA    string
	actionName string
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

	if env := os.Getenv(envAction); len(env) > 0 {
		actionName = env
	} else {
		fmt.Fprintln(os.Stderr, "Missing environment variable:", envAction)
		os.Exit(2)
	}

	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	))

	client = github.NewClient(tc)
}

func createCheck() *github.CheckRun {
	opts := github.CreateCheckRunOptions{
		Name:    actionName,
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

func completeCheck(check *github.CheckRun, concl conclusion) {
	opts := github.UpdateCheckRunOptions{
		Name:       actionName,
		HeadSHA:    github.String(headSHA),
		Conclusion: github.String(concl.String()),
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
				fmt.Sprintf("%s (%s)", strings.Title(f.RuleName), f.RuleName),
			),
			Message: github.String(f.Failure),
		}
	}

	return ann
}

func pushFailures(check *github.CheckRun, failures []*failure, total int, wg *sync.WaitGroup) {
	opts := github.UpdateCheckRunOptions{
		Name:    actionName,
		HeadSHA: github.String(headSHA),
		Output: &github.CheckRunOutput{
			Title:       github.String("Result"),
			Summary:     github.String(fmt.Sprintf("%d failures", total)),
			Annotations: createAnnotations(failures),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, _, err := client.Checks.UpdateCheckRun(ctx, repoOwner, repoName, check.GetID(), opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error while updating check-run:", err)
		os.Exit(1)
	}

	wg.Done()
}

func main() {
	var exitCode int
	var totalCount int
	var concl conclusion

	check := createCheck()

	failures := make([]*failure, 0)

	ch := make(chan *failure)
	go getFailures(ch)

	wg := &sync.WaitGroup{}

	warnCount, errCount := 0, 0
	chunks := 1
	for f := range ch {
		failures = append(failures, f)

		totalCount++

		switch f.Severity {
		case "warning":
			warnCount++
		case "error":
			errCount++
		}

		if c := chunks * chunkLimit; totalCount > c {
			wg.Add(1)
			go pushFailures(check, failures[c-chunkLimit:c], totalCount, wg)
			chunks++
		}
	}

	if totalCount > 0 {
		wg.Add(1)
		if chunks == 1 {
			go pushFailures(check, failures, totalCount, wg)
		} else {
			c := chunks * chunkLimit
			go pushFailures(check, failures[c-chunkLimit:], totalCount, wg)
		}
		wg.Wait()
		exitCode, concl = 1, conclFailure
	}

	completeCheck(check, concl)

	fmt.Fprintf(os.Stderr,
		"Successful run with %d failures (%d warnings, %d errors)\n",
		totalCount, warnCount, errCount)
	os.Exit(exitCode)
}
