package cienv

import (
	"fmt"
)

// CIEnv stores CI enironment variables.
type CIEnv struct {
	PullRequestSlug *Slug
	RepoSlug        *Slug
	IsPullRequest   bool
	PullRequest     int
}

type Detector interface {
	Name() string
	Detect() (*CIEnv, error)
}

var dets []Detector

func Register(d Detector) {
	dets = append(dets, d)
}

func Detect() (*CIEnv, error) {
	err := ErrNotDetected
	for _, d := range dets {
		env, errDetect := d.Detect()
		switch errDetect {
		case nil:
			return env, nil
		case ErrNotDetected:
		default:
			err = fmt.Errorf("%s: %w", d.Name(), errDetect)
		}
	}
	return nil, err
}
