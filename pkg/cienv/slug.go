package cienv

import (
	"errors"
	"strings"
)

type Slug struct {
	Owner string
	Repo  string
}

func (s *Slug) String() string {
	return s.Owner + "/" + s.Repo
}

func NewSlug(slug string) (*Slug, error) {
	ss := strings.Split(slug, "/")
	if len(ss) != 2 {
		return nil, errors.New("invalid repository slug")
	}
	return &Slug{
		ss[0],
		ss[1],
	}, nil
}
