package githubactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/at-wat/gh-pr-comment/pkg/cienv"
)

func init() {
	cienv.Register(&Actions{})
}

type Actions struct {
}

func (t *Actions) Name() string {
	return "GitHub Actions"
}

func (t *Actions) Detect() (*cienv.CIEnv, error) {
	env := &cienv.CIEnv{}

	if isActions, ok := os.LookupEnv("GITHUB_ACTIONS"); !ok || isActions != "true" {
		return nil, cienv.ErrNotDetected
	}

	eventName, ok := os.LookupEnv("GITHUB_EVENT_NAME")
	if !ok {
		return nil, errors.New("GITHUB_EVENT_NAME is not set")
	}
	env.IsPullRequest = eventName == "pull_request"

	if env.IsPullRequest {
		if prEvent, ok := os.LookupEnv("GITHUB_EVENT_PATH"); ok {
			f, err := os.Open(prEvent)
			if err != nil {
				return nil, fmt.Errorf("failed to read event payload: %w", err)
			}
			defer f.Close()
			dec := json.NewDecoder(f)

			var event pullRequestEvent
			if err := dec.Decode(&event); err != nil {
				return nil, fmt.Errorf("failed to parse event payload: %w", err)
			}

			env.PullRequest = event.Number
			env.PullRequestSlug, err = cienv.NewSlug(event.PullRequest.Head.Repo.FullName)
			if err != nil {
				return nil, fmt.Errorf("failed to parse head.repo.full_name: %w", err)
			}
		} else {
			return nil, errors.New("GITHUB_EVENT_PATH is not set")
		}
	}

	if slug, ok := os.LookupEnv("GITHUB_REPOSITORY"); ok {
		var err error
		env.RepoSlug, err = cienv.NewSlug(slug)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GITHUB_REPOSITORY: %w", err)
		}
	}

	return env, nil
}

type pullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Head struct {
			Repo struct {
				FullName string `json:"full_name"`
			} `json:"repo"`
		} `json:"head"`
	} `json:"pull_request"`
}
