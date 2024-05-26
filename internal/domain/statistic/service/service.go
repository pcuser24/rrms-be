package service

import (
	"github.com/google/uuid"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	statistic_repo "github.com/user2410/rrms-backend/internal/domain/statistic/repo"
)

type Service interface {
	GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error)
	GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error)
	GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error)
	GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error)
	GetRentalPaymentIncomes(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentIncomeItem, err error)
	GetPaymentsStatistic(userId uuid.UUID, query dto.PaymentsStatisticQuery) (res []dto.PaymentsStatisticItem, err error)

	// Landing
}

type service struct {
	authRepo      auth_repo.Repo
	listingRepo   listing_repo.Repo
	statisticRepo statistic_repo.Repo
	propertyRepo  property_repo.Repo
}

func NewService(authRepo auth_repo.Repo, listingRepo listing_repo.Repo, statisticRepo statistic_repo.Repo, propertyRepo property_repo.Repo) Service {
	return &service{
		authRepo:      authRepo,
		listingRepo:   listingRepo,
		statisticRepo: statisticRepo,
		propertyRepo:  propertyRepo,
	}
}
