package listing

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*model.ListingModel, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	UpdateListing(data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	AddListingPolicies(lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error)
	AddListingUnits(lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error)
	DeleteListingPolicies(lid uuid.UUID, ids []int64) error
	DeleteListingUnits(lid uuid.UUID, ids []uuid.UUID) error
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

func (s *service) GetListingByID(id uuid.UUID) (*model.ListingModel, error) {
	return s.repo.GetListingByID(context.Background(), id)
}

func (s *service) UpdateListing(data *dto.UpdateListing) error {
	return s.repo.UpdateListing(context.Background(), data)
}

func (s *service) DeleteListing(id uuid.UUID) error {
	return s.repo.DeleteListing(context.Background(), id)
}

func (s *service) AddListingPolicies(lid uuid.UUID, items []dto.CreateListingPolicy) ([]model.ListingPolicyModel, error) {
	return s.repo.AddListingPolicies(context.Background(), lid, items)
}

func (s *service) AddListingUnits(lid uuid.UUID, items []dto.CreateListingUnit) ([]model.ListingUnitModel, error) {
	return s.repo.AddListingUnits(context.Background(), lid, items)
}

func (s *service) DeleteListingPolicies(lid uuid.UUID, ids []int64) error {
	return s.repo.DeleteListingPolicies(context.Background(), lid, ids)
}

func (s *service) DeleteListingUnits(lid uuid.UUID, ids []uuid.UUID) error {
	return s.repo.DeleteListingUnits(context.Background(), lid, ids)
}

func (s *service) CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.repo.CheckListingOwnership(context.Background(), lid, uid)
}

func (s *service) CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.repo.CheckValidUnitForListing(context.Background(), lid, uid)
}
