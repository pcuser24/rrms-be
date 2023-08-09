package property

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
)

type Service interface {
	CheckOwnership(id uuid.UUID, userId uuid.UUID) (bool, error)
	CreateProperty(data *dto.CreateProperty) (*model.PropertyModel, error)
	GetPropertyById(id uuid.UUID) (*model.PropertyModel, error)
	UpdateProperty(data *dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
	AddPropertyMedium(id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error)
	AddPropertyAmenities(id uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error)
	AddPropertyFeatures(id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error)
	AddPropertyTags(id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error)
	GetAllAmenities() ([]model.PAmenity, error)
	GetAllFeatures() ([]model.PFeature, error)
	DeletePropertyAmenities(puid uuid.UUID, aid []int64) error
	DeletePropertyFeatures(puid uuid.UUID, fid []int64) error
	DeletePropertyMedium(puid uuid.UUID, mid []int64) error
	DeletePropertyTags(puid uuid.UUID, tid []int64) error
	ReplacePropertyAmenities(puid uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error)
	ReplacePropertyFeatures(puid uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error)
	ReplacePropertyTags(puid uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error)
	ReplacePropertyMedium(puid uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateProperty(data *dto.CreateProperty) (*model.PropertyModel, error) {
	return s.repo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*model.PropertyModel, error) {
	return s.repo.GetPropertyById(context.Background(), id)
}

func (s *service) UpdateProperty(data *dto.UpdateProperty) error {
	return s.repo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckOwnership(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return s.repo.CheckOwnership(context.Background(), id, userId)
}

func (s *service) DeleteProperty(id uuid.UUID) error {
	return s.repo.DeleteProperty(context.Background(), id)
}

func (s *service) AddPropertyMedium(id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error) {
	return s.repo.AddPropertyMedium(context.Background(), id, items)
}

func (s *service) AddPropertyAmenities(id uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error) {
	return s.repo.AddPropertyAmenities(context.Background(), id, items)
}

func (s *service) AddPropertyFeatures(id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error) {
	return s.repo.AddPropertyFeatures(context.Background(), id, items)
}

func (s *service) AddPropertyTags(id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error) {
	return s.repo.AddPropertyTag(context.Background(), id, items)
}

func (s *service) GetAllAmenities() ([]model.PAmenity, error) {
	return s.repo.GetAllAmenities(context.Background())
}

func (s *service) GetAllFeatures() ([]model.PFeature, error) {
	return s.repo.GetAllFeatures(context.Background())
}

func (s *service) DeletePropertyAmenities(puid uuid.UUID, aid []int64) error {
	return s.repo.DeletePropertyAmenities(context.Background(), puid, aid)
}

func (s *service) DeletePropertyFeatures(puid uuid.UUID, fid []int64) error {
	return s.repo.DeletePropertyFeatures(context.Background(), puid, fid)
}

func (s *service) DeletePropertyMedium(puid uuid.UUID, mid []int64) error {
	return s.repo.DeletePropertyMedium(context.Background(), puid, mid)
}

func (s *service) DeletePropertyTags(puid uuid.UUID, tid []int64) error {
	return s.repo.DeletePropertyTags(context.Background(), puid, tid)
}

func (s *service) ReplacePropertyAmenities(puid uuid.UUID, items []dto.CreatePropertyAmenity) ([]model.PropertyAmenityModel, error) {
	err := s.repo.DeleteAllPropertyAmenities(context.Background(), puid)
	if err != nil {
		return nil, err
	}

	return s.repo.AddPropertyAmenities(context.Background(), puid, items)
}

func (s *service) ReplacePropertyFeatures(puid uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error) {
	err := s.repo.DeleteAllPropertyFeatures(context.Background(), puid)
	if err != nil {
		return nil, err
	}

	return s.repo.AddPropertyFeatures(context.Background(), puid, items)
}

func (s *service) ReplacePropertyTags(puid uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error) {
	err := s.repo.DeleteAllPropertyTags(context.Background(), puid)
	if err != nil {
		return nil, err
	}

	return s.repo.AddPropertyTag(context.Background(), puid, items)
}

func (s *service) ReplacePropertyMedium(puid uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error) {
	// strategy:
	// delete all media
	// save new media
	err := s.repo.DeleteAllPropertyMedium(context.Background(), puid)
	if err != nil {
		return nil, err
	}

	return s.repo.AddPropertyMedium(context.Background(), puid, items)
}
