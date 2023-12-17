package listing

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateListing(data *dto.CreateListing) (*model.ListingModel, error) {
	return s.repo.CreateListing(context.Background(), data)
}

func (s *service) SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	if data.Limit == nil {
		data.Limit = types.Ptr[int32](10)
	}
	if data.Offset == nil {
		data.Offset = types.Ptr[int32](0)
	}
	if data.SortBy == nil {
		data.SortBy = types.Ptr[string]("post_at")
	}
	if data.Order == nil {
		data.Order = types.Ptr[string]("desc")
	}
	return s.repo.SearchListingCombination(context.Background(), data)
}

func (s *service) GetListingByID(id uuid.UUID) (*model.ListingModel, error) {
	return s.repo.GetListingByID(context.Background(), id)
}

func (s *service) UpdateListing(data *dto.UpdateListing) error {
	return s.repo.UpdateListing(context.Background(), data)
}

func (s *service) DeleteListing(id uuid.UUID) error {
	return s.repo.DeleteListing(context.Background(), id)
}

func (s *service) CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.repo.CheckListingOwnership(context.Background(), lid, uid)
}

func (s *service) CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.repo.CheckValidUnitForListing(context.Background(), lid, uid)
}
