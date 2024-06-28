package repo

import (
	"log"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/redisd"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/config"
)

var (
	basePath         = utils.GetBasePath()
	testAuthRepo     auth_repo.Repo
	testPropertyRepo property_repo.Repo
	testUnitRepo     unit_repo.Repo
	testListingRepo  Repo
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})
	redisClient := redisd.NewRedisClient(rdb)

	testListingRepo = NewRepo(dao, redisClient)
	testAuthRepo = auth_repo.NewRepo(dao)
	testPropertyRepo = property_repo.NewRepo(dao, redisClient)
	testUnitRepo = unit_repo.NewRepo(dao, redisClient)

	os.Exit(m.Run())
}
