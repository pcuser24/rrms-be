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
	GetUnitsOfProperty(id uuid.UUID) ([]model.UnitModel, error)
	UpdateUnit(data *dto.UpdateUnit) error
	DeleteUnit(id uuid.UUID) error
	CheckUnitOwnership(id uuid.UUID, userId uuid.UUID) (bool, error)
	AddUnitAmenities(uid uuid.UUID, items []dto.CreateUnitAmenity) ([]model.UnitAmenityModel, error)
	AddUnitMedium(uid uuid.UUID, items []dto.CreateUnitMedia) ([]model.UnitMediaModel, error)
	GetAllAmenities() ([]model.UAmenity, error)
	DeleteUnitAmenities(uid uuid.UUID, ids []int64) error
	DeleteUnitMedium(uid uuid.UUID, ids []int64) error
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

func (s *service) GetUnitsOfProperty(id uuid.UUID) ([]model.UnitModel, error) {
	return s.repo.GetUnitsOfProperty(context.Background(), id)
}

func (s *service) UpdateUnit(data *dto.UpdateUnit) error {
	return s.repo.UpdateUnit(context.Background(), data)
}

func (s *service) DeleteUnit(id uuid.UUID) error {
	return s.repo.DeleteUnit(context.Background(), id)
}

func (s *service) AddUnitAmenities(uid uuid.UUID, items []dto.CreateUnitAmenity) ([]model.UnitAmenityModel, error) {
	return s.repo.AddUnitAmenities(context.Background(), uid, items)
}

func (s *service) AddUnitMedium(uid uuid.UUID, items []dto.CreateUnitMedia) ([]model.UnitMediaModel, error) {
	return s.repo.AddUnitMedium(context.Background(), uid, items)
}

func (s *service) GetAllAmenities() ([]model.UAmenity, error) {
	return s.repo.GetAllAmenities(context.Background())
}

func (s *service) DeleteUnitAmenities(uid uuid.UUID, ids []int64) error {
	return s.repo.DeleteUnitAmenities(context.Background(), uid, ids)
}

func (s *service) DeleteUnitMedium(uid uuid.UUID, ids []int64) error {
	return s.repo.DeleteUnitMedium(context.Background(), uid, ids)
}

func (s *service) CheckUnitOwnership(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return s.repo.CheckUnitOwnership(context.Background(), id, userId)
}
