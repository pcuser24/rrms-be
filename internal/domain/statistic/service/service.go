package service

import (
	"github.com/google/uuid"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/domain/statistic/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
)

type Service interface {
	GetPropertiesStatistic(userId uuid.UUID, query dto.PropertiesStatisticQuery) (res dto.PropertiesStatisticResponse, err error)
	GetApplicationStatistic(userId uuid.UUID) (res dto.ApplicationStatisticResponse, err error)
	GetRentalStatistic(userId uuid.UUID) (res dto.RentalStatisticResponse, err error)
	GetRentalPaymentArrears(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentArrearsItem, err error)
	GetRentalPaymentIncomes(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (res []dto.RentalPaymentIncomeItem, err error)
	GetPaymentsStatistic(userId uuid.UUID, query dto.PaymentsStatisticQuery) (res []dto.PaymentsStatisticItem, err error)
	GetTenantRentalStatistic(userId uuid.UUID) (res dto.TenantRentalStatisticResponse, err error)
	GetTenantMaintenanceStatistic(userId uuid.UUID) (res dto.TenantMaintenanceStatisticResponse, err error)
	GetTenantExpenditureStatistic(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) ([]dto.TenantExpenditureStatisticItem, error)
	GetTenantArrearsStatistic(userId uuid.UUID, query *dto.RentalPaymentStatisticQuery) (dto.TenantArrearsStatistic, error)
	// Landing
	GetRecentListings(limit int32, fields []string) ([]listing_model.ListingModel, error)
	GetListingSuggestion(id uuid.UUID, limit int) (dto.ListingsSuggestionResult, error)
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
