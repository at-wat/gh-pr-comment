package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"

	"github.com/at-wat/gh-pr-comment/pkg/cienv"
	_ "github.com/at-wat/gh-pr-comment/pkg/cienv/githubactions"
	_ "github.com/at-wat/gh-pr-comment/pkg/cienv/travisci"
)

func main() {
	repl := flag.String("stdin", "", "replace this keyword in comment by text from stdin")
	pr := flag.Int("pr", 0, "override PR number")
	timeout := flag.Duration("timeout", time.Minute, "GitHub API timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [option] title comment\n", os.Args[0])
		fmt.Fprint(os.Stderr, `env:
  GITHUB_TOKEN (or TRAVIS_BOT_GITHUB_TOKEN):
    token with comment write permission

env for Travis-CI:
  TRAVIS:
    must be true
  TRAVIS_PULL_REQUEST_SLUG:
    owner/repos (TRAVIS_REPO_SLUG is used if not set)
  TRAVIS_PULL_REQUEST:
    pull request number

env for GitHub Actions:
  GITHUB_ACTIONS:
    must be true
	GITHUB_REPOSITORY:
    owner/repos (TRAVIS_REPO_SLUG is used if not set)
	GITHUB_EVENT_NAME:
		action event name
	GITHUB_EVENT_PATH:
		path to webhook payload

options:
`)

		flag.PrintDefaults()
		os.Exit(1)
	}
	title := flag.Args()[0]
	body := flag.Args()[1]

	env, err := cienv.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *pr > 0 {
		env.IsPullRequest = true
		env.PullRequest = *pr
	}

	if !env.IsPullRequest {
		fmt.Fprint(os.Stderr, "info: not on a pull request\n")
		os.Exit(0)
	}

	var ghToken string

	if t, ok := os.LookupEnv("TRAVIS_BOT_GITHUB_TOKEN"); ok {
		ghToken = t
	} else if t, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		ghToken = t
	}

	if ghToken == "" {
		fmt.Fprint(os.Stderr, "error: neither TRAVIS_BOT_GITHUB_TOKEN nor GITHUB_TOKEN is not set.\n")
		os.Exit(1)
	}

	if *repl != "" {
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, "error: failed to read from stdin.\n")
			os.Exit(1)
		}
		body = strings.Replace(body, *repl, string(in), 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	tc := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken}),
	)

	var gh *github.Client
	if ep, customEP := os.LookupEnv("GITHUB_API_URL_BASE"); customEP {
		var err error
		gh, err = github.NewEnterpriseClient(ep, ep, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to init GitHub client: %v\n", err)
			os.Exit(1)
		}
	} else {
		gh = github.NewClient(tc)
	}

	bodyStr := fmt.Sprintf("## %s\n\n%s", title, body)
	_, resp, err := gh.Issues.CreateComment(
		ctx,
		env.RepoSlug.Owner, env.RepoSlug.Repo,
		env.PullRequest,
		&github.IssueComment{Body: &bodyStr},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to call CreateComment API: %v\n", err)
		os.Exit(1)
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		fmt.Fprintf(os.Stderr, "error: failed to CreateComment: %s\n", resp.Status)
		os.Exit(1)
	}
}
