package travis

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/at-wat/gh-pr-comment/pkg/cienv"
)

func init() {
	cienv.Register(&Travis{})
}

type Travis struct {
}

func (t *Travis) Name() string {
	return "Travis-CI"
}

func (t *Travis) Detect() (*cienv.CIEnv, error) {
	env := &cienv.CIEnv{}

	if isTravis, ok := os.LookupEnv("TRAVIS"); !ok || isTravis != "true" {
		return nil, cienv.ErrNotDetected
	}

	pr, ok := os.LookupEnv("TRAVIS_PULL_REQUEST")
	if !ok {
		return nil, errors.New("TRAVIS_PULL_REQUEST is not set")
	}
	env.IsPullRequest = pr != "false"
	if env.IsPullRequest {
		var err error
		env.PullRequest, err = strconv.Atoi(pr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TRAVIS_PULL_REQUEST: %w", err)
		}
	}

	if env.IsPullRequest {
		if prSlug, ok := os.LookupEnv("TRAVIS_PULL_REQUEST_SLUG"); ok {
			var err error
			env.PullRequestSlug, err = cienv.NewSlug(prSlug)
			if err != nil {
				return nil, fmt.Errorf("failed to parse TRAVIS_PULL_REQUEST_SLUG: %w", err)
			}
		} else {
			return nil, errors.New("TRAVIS_PULL_REQUEST_SLUG is not set")
		}
	}

	if slug, ok := os.LookupEnv("TRAVIS_REPO_SLUG"); ok {
		var err error
		env.RepoSlug, err = cienv.NewSlug(slug)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TRAVIS_REPO_SLUG: %w", err)
		}
	}

	return env, nil
}
