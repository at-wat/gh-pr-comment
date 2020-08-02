package cienv

import (
	"testing"
)

func TestSlug(t *testing.T) {
	var s *Slug
	var err error

	if s, err = NewSlug("aaa"); err == nil {
		t.Fatalf("Invalid slug must not be unmarshalled")
	}

	if s, err = NewSlug("aaa/bbb"); err != nil {
		t.Fatalf("Valid slug must be unmarshalled")
	}

	if s.Owner != "aaa" {
		t.Errorf("Owner field is wrong: %s", s.Owner)
	}
	if s.Repo != "bbb" {
		t.Errorf("Repo field is wrong: %s", s.Repo)
	}

	if str := s.String(); str != "aaa/bbb" {
		t.Errorf("Stringer returns wrong string: %s", str)
	}
}
