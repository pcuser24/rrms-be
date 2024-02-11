package repo

import (
	"log"
	"os"
	"testing"

	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/config"
)

var (
	basePath         = utils.GetBasePath()
	testAuthRepo     auth_repo.Repo
	testPropertyRepo property_repo.Repo
	testUnitRepo     Repo
)

func TestMain(m *testing.M) {
	// init connection to database
	conf, err := config.NewTestRepoConfig(basePath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	dao, err := database.NewPostgresDAO(conf.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testUnitRepo = NewRepo(dao)
	testPropertyRepo = property_repo.NewRepo(dao)
	testAuthRepo = auth_repo.NewRepo(dao)

	os.Exit(m.Run())
}
