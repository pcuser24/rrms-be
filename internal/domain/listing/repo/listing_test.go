package repo

import "testing"

func TestCreateRepo(t *testing.T) {
	NewRandomListingDB(t, testAuthRepo, testPropertyRepo, testUnitRepo, testListingRepo)
}
