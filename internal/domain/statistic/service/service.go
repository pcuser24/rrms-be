package service

import (
	"github.com/google/uuid"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	statistic_repo "github.com/user2410/rrms-backend/internal/domain/statistic/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
)

type Service interface {
	GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error)
	GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error)
	GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error)
	GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error)
	GetRentalPaymentIncomes(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentIncomeItem, err error)
	GetPaymentsStatistic(userId uuid.UUID, query dto.PaymentsStatisticQuery) (res []dto.PaymentsStatisticItem, err error)

	// Landing
	GetRecentListings(limit int32, fields []string) ([]listing_model.ListingModel, error)
	GetListingSuggestion(id uuid.UUID, limit int) (dto.ListingsSuggestionResult, error)
}

type service struct {
	// Repositories
	authRepo      auth_repo.Repo
	listingRepo   listing_repo.Repo
	statisticRepo statistic_repo.Repo
	propertyRepo  property_repo.Repo
	unitRepo      unit_repo.Repo
	// ElasticSearch
	esClient *es.ElasticSearchClient
}

func NewService(
	authRepo auth_repo.Repo, listingRepo listing_repo.Repo, statisticRepo statistic_repo.Repo, propertyRepo property_repo.Repo, unitRepo unit_repo.Repo,
	esClient *es.ElasticSearchClient,
) Service {
	return &service{
		authRepo:      authRepo,
		listingRepo:   listingRepo,
		statisticRepo: statisticRepo,
		propertyRepo:  propertyRepo,
		unitRepo:      unitRepo,

		esClient: esClient,
	}
}
