package property

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	unitModel "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateProperty(data *dto.CreateProperty, creatorID uuid.UUID) (*model.PropertyModel, error)
	CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error)
	CheckManageability(id uuid.UUID, userId uuid.UUID) (bool, error)
	GetPropertyById(id uuid.UUID) (*model.PropertyModel, error)
	GetPropertiesByIds(ids []uuid.UUID, fields []string) ([]model.PropertyModel, error)
	GetUnitsOfProperty(id uuid.UUID) ([]unitModel.UnitModel, error)
	GetPropertiesOfUser(userId uuid.UUID, fields []string) ([]getPropertiesOfUserResponse, error)
	SearchListingCombination(data *dto.SearchPropertyCombinationQuery) (*dto.SearchPropertyCombinationResponse, error)
	UpdateProperty(data *dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
	GetAllFeatures() ([]model.PFeature, error)
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

func (s *service) GetPropertiesByIds(ids []uuid.UUID, fields []string) ([]model.PropertyModel, error) {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = id.String()
	}
	return s.repo.GetPropertiesByIds(context.Background(), idsStr, fields)
}

func (s *service) GetUnitsOfProperty(id uuid.UUID) ([]unitModel.UnitModel, error) {
	return s.repo.GetUnitsOfProperty(context.Background(), id)
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

func (s *service) GetAllFeatures() ([]model.PFeature, error) {
	return s.repo.GetAllFeatures(context.Background())
}

type getPropertiesOfUserResponse struct {
	Role     string              `json:"role"`
	Property model.PropertyModel `json:"property"`
}

func (s *service) GetPropertiesOfUser(userId uuid.UUID, fields []string) ([]getPropertiesOfUserResponse, error) {
	managedProps, err := s.repo.GetManagedProperties(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	var pids []string
	for _, p := range managedProps {
		pid := p.PropertyID.String()
		pids = append(pids, pid)
	}

	ps, err := s.repo.GetPropertiesByIds(context.Background(), pids, fields)
	if err != nil {
		return nil, err
	}

	var res []getPropertiesOfUserResponse
	for _, p := range managedProps {
		r := getPropertiesOfUserResponse{Role: p.Role}
		for i, pp := range ps {
			if pp.ID == p.PropertyID {
				r.Property = ps[i]
			}
		}
		res = append(res, r)
	}

	return res, nil
}

func (s *service) SearchListingCombination(data *dto.SearchPropertyCombinationQuery) (*dto.SearchPropertyCombinationResponse, error) {
	data.SortBy = types.Ptr(utils.PtrDerefence[string](data.SortBy, "created_at"))
	data.Order = types.Ptr(utils.PtrDerefence[string](data.Order, "desc"))
	data.Limit = types.Ptr(utils.PtrDerefence[int32](data.Limit, 1000))
	data.Offset = types.Ptr(utils.PtrDerefence[int32](data.Offset, 0))
	return s.repo.SearchPropertyCombination(context.Background(), data)
}
