package service

import (
	"github.com/google/uuid"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	statistic_dto "github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
)

type Service interface {
	GetPropertiesStatistic(userId uuid.UUID, query statistic_dto.PropertiesStatisticQuery) (res statistic_dto.PropertiesStatisticResponse, err error)
	GetApplicationStatistic(userId uuid.UUID) (res statistic_dto.ApplicationStatisticResponse, err error)
	GetRentalStatistic(userId uuid.UUID) (res statistic_dto.RentalStatisticResponse, err error)
	GetRentalPaymentArrears(userId uuid.UUID, query *statistic_dto.RentalPaymentStatisticQuery) (res []statistic_dto.RentalPaymentArrearsItem, err error)
	GetRentalPaymentIncomes(userId uuid.UUID, query *statistic_dto.RentalPaymentStatisticQuery) (res []statistic_dto.RentalPaymentIncomeItem, err error)
	GetPaymentsStatistic(userId uuid.UUID, query statistic_dto.PaymentsStatisticQuery) (res []statistic_dto.PaymentsStatisticItem, err error)
	GetTenantRentalStatistic(userId uuid.UUID) (res statistic_dto.TenantRentalStatisticResponse, err error)
	GetTenantMaintenanceStatistic(userId uuid.UUID) (res statistic_dto.TenantMaintenanceStatisticResponse, err error)
	GetTenantExpenditureStatistic(userId uuid.UUID, query *statistic_dto.RentalPaymentStatisticQuery) ([]statistic_dto.TenantExpenditureStatisticItem, error)
	GetTenantArrearsStatistic(userId uuid.UUID, query *statistic_dto.RentalPaymentStatisticQuery) (statistic_dto.TenantArrearsStatistic, error)
	GetTotalTenantsManagedByUserStatistic(userId uuid.UUID, query *statistic_dto.RentalStatisticQuery) (int32, error)
	GetTotalTenantsOfUnitStatistic(unitId uuid.UUID) (int32, error)
	// Landing
	GetRecentListings(limit int32, fields []string) ([]listing_model.ListingModel, error)
	GetSimilarListingsToListing(id uuid.UUID, limit int) (statistic_dto.ListingsSuggestionResult, error)
}

type service struct {
	// Repositories
	domainRepo repos.DomainRepo
	// ElasticSearch
	esClient *es.ElasticSearchClient
}

func NewService(
	domainRepo repos.DomainRepo,
	esClient *es.ElasticSearchClient,
) Service {
	return &service{
		domainRepo: domainRepo,

		esClient: esClient,
	}
}
