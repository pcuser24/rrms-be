package repo

import (
	"testing"
)

func TestInsertUser(t *testing.T) {
	NewRandomUserDB(t, testRepo)
}
