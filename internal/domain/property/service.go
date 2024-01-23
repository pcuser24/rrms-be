package property

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	unitDto "github.com/user2410/rrms-backend/internal/domain/unit/dto"
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
	GetPropertiesOfUser(userId uuid.UUID, fields []string) ([]GetPropertiesOfUserItem, error)
	SearchListingCombination(data *dto.SearchPropertyCombinationQuery) (*dto.SearchPropertyCombinationResponse, error)
	UpdateProperty(data *dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
	GetAllFeatures() ([]model.PFeature, error)
}

// import cycle is not allowed
type unitRepo interface {
	GetUnitById(ctx context.Context, id uuid.UUID) (*unitModel.UnitModel, error)
	SearchUnitCombination(ctx context.Context, query *unitDto.SearchUnitCombinationQuery) (*unitDto.SearchUnitCombinationResponse, error)
}

type service struct {
	pRepo Repo
	uRepo unitRepo
}

func NewService(pRepo Repo, uRepo unitRepo) Service {
	return &service{
		pRepo: pRepo,
		uRepo: uRepo,
	}
}

func (s *service) CreateProperty(data *dto.CreateProperty, creatorID uuid.UUID) (*model.PropertyModel, error) {
	data.CreatorID = creatorID
	data.Managers = append(data.Managers, dto.CreatePropertyManager{
		ManagerID: creatorID,
		Role:      "OWNER", // TODO: add role to user
	})
	return s.pRepo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*model.PropertyModel, error) {
	return s.pRepo.GetPropertyById(context.Background(), id)
}

func (s *service) GetPropertiesByIds(ids []uuid.UUID, fields []string) ([]model.PropertyModel, error) {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = id.String()
	}
	return s.pRepo.GetPropertiesByIds(context.Background(), idsStr, fields)
}

func (s *service) GetUnitsOfProperty(id uuid.UUID) ([]unitModel.UnitModel, error) {
	ids, err := s.uRepo.SearchUnitCombination(
		context.Background(),
		&unitDto.SearchUnitCombinationQuery{
			SearchUnitQuery: unitDto.SearchUnitQuery{
				UPropertyID: types.Ptr[string](id.String()),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	res := make([]unitModel.UnitModel, 0, len(ids.Items))
	for _, id := range ids.Items {
		_res, err := s.uRepo.GetUnitById(context.Background(), id.UId)
		if err != nil {
			return nil, err
		}
		res = append(res, *_res)
	}
	return res, nil
}

func (s *service) UpdateProperty(data *dto.UpdateProperty) error {
	return s.pRepo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckManageability(id uuid.UUID, userId uuid.UUID) (bool, error) {
	managers, err := s.pRepo.GetPropertyManagers(context.Background(), id)
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
	isPublic, err := s.pRepo.IsPublic(context.Background(), id)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}
	managers, err := s.pRepo.GetPropertyManagers(context.Background(), id)
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
	return s.pRepo.DeleteProperty(context.Background(), id)
}

func (s *service) GetAllFeatures() ([]model.PFeature, error) {
	return s.pRepo.GetAllFeatures(context.Background())
}

type GetPropertiesOfUserItem struct {
	Role     string              `json:"role"`
	Property model.PropertyModel `json:"property"`
}

func (s *service) GetPropertiesOfUser(userId uuid.UUID, fields []string) ([]GetPropertiesOfUserItem, error) {
	managedProps, err := s.pRepo.GetManagedProperties(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	var pids []string
	for _, p := range managedProps {
		pid := p.PropertyID.String()
		pids = append(pids, pid)
	}

	ps, err := s.pRepo.GetPropertiesByIds(context.Background(), pids, fields)
	if err != nil {
		return nil, err
	}

	var res []GetPropertiesOfUserItem
	for _, p := range managedProps {
		r := GetPropertiesOfUserItem{Role: p.Role}
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
	return s.pRepo.SearchPropertyCombination(context.Background(), data)
}
