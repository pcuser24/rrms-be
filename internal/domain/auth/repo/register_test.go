package repo

import (
	"testing"
)

func TestInsertUser(t *testing.T) {
	NewRandomUser(t, testRepo)
}
