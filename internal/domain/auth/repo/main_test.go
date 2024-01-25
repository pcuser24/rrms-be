package repo

import (
	"log"
	"os"
	"testing"

	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/config"
)

var (
	basePath = utils.GetBasePath()
	testRepo Repo
)

func TestMain(m *testing.M) {
	// init connection to database
	conf, err := config.NewTestRepoConfig(basePath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	dao, err := database.NewDAO(conf.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testRepo = NewRepo(dao)

	os.Exit(m.Run())
}
