package repos

import (
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	misc_repo "github.com/user2410/rrms-backend/internal/domain/misc/repo"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	reminder_repo "github.com/user2410/rrms-backend/internal/domain/reminder/repo"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	statistic_repo "github.com/user2410/rrms-backend/internal/domain/statistic/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/redisd"
	"go.uber.org/mock/gomock"
)

type DomainRepo struct {
	AuthRepo        auth_repo.Repo
	PropertyRepo    property_repo.Repo
	UnitRepo        unit_repo.Repo
	ListingRepo     listing_repo.Repo
	RentalRepo      rental_repo.Repo
	ApplicationRepo application_repo.Repo
	PaymentRepo     payment_repo.Repo
	ChatRepo        chat_repo.Repo
	ReminderRepo    reminder_repo.Repo
	StatisticRepo   statistic_repo.Repo
	MiscRepo        misc_repo.Repo
}

func NewDomainRepo(dao database.DAO, redisClient redisd.RedisClient) DomainRepo {
	return DomainRepo{
		AuthRepo:        auth_repo.NewRepo(dao),
		PropertyRepo:    property_repo.NewRepo(dao, redisClient),
		UnitRepo:        unit_repo.NewRepo(dao, redisClient),
		ListingRepo:     listing_repo.NewRepo(dao, redisClient),
		RentalRepo:      rental_repo.NewRepo(dao),
		ApplicationRepo: application_repo.NewRepo(dao),
		PaymentRepo:     payment_repo.NewRepo(dao),
		ChatRepo:        chat_repo.NewRepo(dao),
		ReminderRepo:    reminder_repo.NewRepo(dao, redisClient),
		StatisticRepo:   statistic_repo.NewRepo(dao),
		MiscRepo:        misc_repo.NewRepo(dao),
	}
}

func NewDomainRepoFromMockCtrl(ctrl *gomock.Controller) DomainRepo {
	return DomainRepo{
		AuthRepo:        auth_repo.NewMockRepo(ctrl),
		PropertyRepo:    property_repo.NewMockRepo(ctrl),
		UnitRepo:        unit_repo.NewMockRepo(ctrl),
		ListingRepo:     listing_repo.NewMockRepo(ctrl),
		RentalRepo:      rental_repo.NewMockRepo(ctrl),
		ApplicationRepo: application_repo.NewMockRepo(ctrl),
		PaymentRepo:     payment_repo.NewMockRepo(ctrl),
		ChatRepo:        chat_repo.NewMockRepo(ctrl),
		ReminderRepo:    reminder_repo.NewMockRepo(ctrl),
		StatisticRepo:   statistic_repo.NewMockRepo(ctrl),
		MiscRepo:        misc_repo.NewMockRepo(ctrl),
	}
}
