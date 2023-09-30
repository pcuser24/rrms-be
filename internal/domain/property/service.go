package property

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
)

type Service interface {
	CreateProperty(data *dto.CreateProperty, creatorID uuid.UUID) (*model.PropertyModel, error)
	CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error)
	CheckManageability(id uuid.UUID, userId uuid.UUID) (bool, error)
	GetPropertyById(id uuid.UUID) (*model.PropertyModel, error)
	UpdateProperty(data *dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
	AddPropertyMedia(id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error)
	AddPropertyFeatures(id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error)
	AddPropertyTags(id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error)
	AddPropertyManagers(id uuid.UUID, items []dto.CreatePropertyManager) ([]model.PropertyManagerModel, error)
	GetAllFeatures() ([]model.PFeature, error)
	DeletePropertyFeatures(puid uuid.UUID, fid []int64) error
	DeletePropertyMedia(puid uuid.UUID, mid []int64) error
	DeletePropertyTags(puid uuid.UUID, tid []int64) error
	DeletePropertyManager(puid uuid.UUID, mid uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateProperty(data *dto.CreateProperty, creatorID uuid.UUID) (*model.PropertyModel, error) {
	data.CreatorID = creatorID
	data.Managers = append(data.Managers, dto.CreatePropertyManager{
		ManagerID: creatorID,
		Role:      "OWNER",
	})
	return s.repo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*model.PropertyModel, error) {
	return s.repo.GetPropertyById(context.Background(), id)
}

func (s *service) UpdateProperty(data *dto.UpdateProperty) error {
	return s.repo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckManageability(id uuid.UUID, userId uuid.UUID) (bool, error) {
	managers, err := s.repo.GetPropertyManagers(context.Background(), id)
	if err != nil {
		return false, err
	}
	for _, manager := range managers {
		if manager.ManagerID == userId {
			return true, nil
		}
	}
	return false, nil
}
func (s *service) CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error) {
	isPublic, err := s.repo.IsPublic(context.Background(), id)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}
	managers, err := s.repo.GetPropertyManagers(context.Background(), id)
	if err != nil {
		return false, err
	}
	for _, manager := range managers {
		if manager.ManagerID == uid {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) DeleteProperty(id uuid.UUID) error {
	return s.repo.DeleteProperty(context.Background(), id)
}

func (s *service) AddPropertyMedia(id uuid.UUID, items []dto.CreatePropertyMedia) ([]model.PropertyMediaModel, error) {
	return s.repo.AddPropertyMedia(context.Background(), id, items)
}

func (s *service) AddPropertyFeatures(id uuid.UUID, items []dto.CreatePropertyFeature) ([]model.PropertyFeatureModel, error) {
	return s.repo.AddPropertyFeatures(context.Background(), id, items)
}

func (s *service) AddPropertyTags(id uuid.UUID, items []dto.CreatePropertyTag) ([]model.PropertyTagModel, error) {
	return s.repo.AddPropertyTag(context.Background(), id, items)
}

func (s *service) GetAllFeatures() ([]model.PFeature, error) {
	return s.repo.GetAllFeatures(context.Background())
}

func (s *service) DeletePropertyFeatures(puid uuid.UUID, fid []int64) error {
	return s.repo.DeletePropertyFeatures(context.Background(), puid, fid)
}

func (s *service) DeletePropertyMedia(puid uuid.UUID, mid []int64) error {
	return s.repo.DeletePropertyMedia(context.Background(), puid, mid)
}

func (s *service) DeletePropertyTags(puid uuid.UUID, tid []int64) error {
	return s.repo.DeletePropertyTags(context.Background(), puid, tid)
}

func (s *service) AddPropertyManagers(id uuid.UUID, items []dto.CreatePropertyManager) ([]model.PropertyManagerModel, error) {
	return s.repo.AddPropertyManagers(context.Background(), id, items)
}

func (s *service) DeletePropertyManager(puid uuid.UUID, mid uuid.UUID) error {
	return s.repo.DeletePropertyManager(context.Background(), puid, mid)
}
