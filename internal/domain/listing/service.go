package listing

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/listing/repo"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	GetListingsByIds(ids []uuid.UUID, fields []string) ([]model.ListingModel, error)
	GetListingsOfUser(userId uuid.UUID, fields []string) ([]model.ListingModel, error)
	UpdateListing(data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error)
}

type service struct {
	lRepo repo.Repo
	pRepo property_repo.Repo
}

func NewService(lRepo repo.Repo, pRepo property_repo.Repo) Service {
	return &service{
		lRepo: lRepo,
		pRepo: pRepo,
	}
}

func (s *service) CreateListing(data *dto.CreateListing) (*model.ListingModel, error) {
	listing, err := s.lRepo.CreateListing(context.Background(), data)
	if err != nil {
		return nil, err
	}
	err = s.pRepo.UpdateProperty(context.Background(), &property_dto.UpdateProperty{
		ID:       listing.PropertyID,
		IsPublic: types.Ptr[bool](true),
	})
	if err != nil {
		return listing, err
	}
	return listing, nil
}

func (s *service) SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	data.SortBy = types.Ptr(utils.PtrDerefence[string](data.SortBy, "created_at"))
	data.Order = types.Ptr(utils.PtrDerefence[string](data.Order, "desc"))
	data.Limit = types.Ptr(utils.PtrDerefence[int32](data.Limit, 1000))
	data.Offset = types.Ptr(utils.PtrDerefence[int32](data.Offset, 0))
	return s.lRepo.SearchListingCombination(context.Background(), data)
}

func (s *service) GetListingByID(id uuid.UUID) (*model.ListingModel, error) {
	return s.lRepo.GetListingByID(context.Background(), id)
}

func (s *service) GetListingsByIds(ids []uuid.UUID, fields []string) ([]model.ListingModel, error) {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = id.String()
	}
	return s.lRepo.GetListingsByIds(context.Background(), idsStr, fields)
}

func (s *service) UpdateListing(data *dto.UpdateListing) error {
	return s.lRepo.UpdateListing(context.Background(), data)
}

func (s *service) DeleteListing(id uuid.UUID) error {
	return s.lRepo.DeleteListing(context.Background(), id)
}

func (s *service) CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.lRepo.CheckListingOwnership(context.Background(), lid, uid)
}

func (s *service) CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.lRepo.CheckValidUnitForListing(context.Background(), lid, uid)
}

func (s *service) GetListingsOfUser(userId uuid.UUID, fields []string) ([]model.ListingModel, error) {
	myListings, err := s.lRepo.SearchListingCombination(context.Background(), &dto.SearchListingCombinationQuery{
		SearchListingQuery: dto.SearchListingQuery{
			LCreatorID: types.Ptr[string](userId.String()),
		},
		SearchSortPaginationQuery: requests.SearchSortPaginationQuery{
			Limit:  types.Ptr[int32](1000),
			Offset: types.Ptr[int32](0),
			SortBy: types.Ptr[string]("created_at"),
			Order:  types.Ptr[string]("desc"),
		},
	})
	if err != nil {
		return nil, err
	}

	var lids []string
	for _, listing := range myListings.Items {
		lids = append(lids, listing.LId.String())
	}

	return s.lRepo.GetListingsByIds(context.Background(), lids, fields)

}
