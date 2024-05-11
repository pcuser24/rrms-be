package repo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

var (
	dao           database.DAO
	statisticRepo Repo
)

func TestMain(m *testing.M) {
	_dao, err := database.NewPostgresDAO("postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable")
	if err != nil {
		panic(err)
	}

	dao = _dao
	statisticRepo = NewRepo(dao)

	os.Exit(m.Run())
}

func TestGetNewApplications(t *testing.T) {
	res, err := statisticRepo.GetNewApplications(context.Background(), uuid.MustParse("e0a8d123-c55b-4230-91e8-bd1b7b762366"), time.Now().AddDate(0, -2, 0))
	require.NoError(t, err)

	t.Log(res)
}
