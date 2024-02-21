package mock

import (
	"fmt"

	"go.uber.org/mock/gomock"
)

type inMatcher struct {
	s []any
}

func InMatcher(s []any) inMatcher {
	return inMatcher{
		s: s,
	}
}

// Matches returns whether x is a match.
func (m inMatcher) Matches(x any) bool {
	for _, s := range m.s {
		if gomock.Eq(x).Matches(s) {
			return true
		}
	}
	return false
}

// String describes what the matcher matches.
func (m inMatcher) String() string {
	return fmt.Sprintf("is in %v", m.s)
}
