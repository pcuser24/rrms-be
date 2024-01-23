package unit

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
)

type Service interface {
	CreateUnit(data *dto.CreateUnit) (*model.UnitModel, error)
	GetUnitById(id uuid.UUID) (*model.UnitModel, error)
	GetUnitsByIds(ids []uuid.UUID, fields []string) ([]model.UnitModel, error)
	SearchUnit(query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error)
	UpdateUnit(data *dto.UpdateUnit) error
	DeleteUnit(id uuid.UUID) error
	CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error)
	CheckUnitManageability(id uuid.UUID, userId uuid.UUID) (bool, error)
	CheckUnitOfProperty(pid, uid uuid.UUID) (bool, error)
	GetAllAmenities() ([]model.UAmenity, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUnit(data *dto.CreateUnit) (*model.UnitModel, error) {
	return s.repo.CreateUnit(context.Background(), data)
}

func (s *service) GetUnitById(id uuid.UUID) (*model.UnitModel, error) {
	return s.repo.GetUnitById(context.Background(), id)
}

func (s *service) GetUnitsByIds(ids []uuid.UUID, fields []string) ([]model.UnitModel, error) {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = id.String()
	}
	return s.repo.GetUnitsByIds(context.Background(), idsStr, fields)
}

func (s *service) UpdateUnit(data *dto.UpdateUnit) error {
	return s.repo.UpdateUnit(context.Background(), data)
}

func (s *service) DeleteUnit(id uuid.UUID) error {
	return s.repo.DeleteUnit(context.Background(), id)
}

func (s *service) GetAllAmenities() ([]model.UAmenity, error) {
	return s.repo.GetAllAmenities(context.Background())
}

func (s *service) CheckUnitManageability(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return s.repo.CheckUnitManageability(context.Background(), id, userId)
}

func (s *service) CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error) {
	isPublic, err := s.repo.IsPublic(context.Background(), id)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}
	return s.CheckUnitManageability(id, uid)
}

func (s *service) CheckUnitOfProperty(pid, uid uuid.UUID) (bool, error) {
	return s.repo.CheckUnitOfProperty(context.Background(), pid, uid)
}

func (s *service) SearchUnit(query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error) {
	return s.repo.SearchUnitCombination(context.Background(), query)
}
