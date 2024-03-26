package property

import (
	"context"
	"fmt"

	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	"github.com/google/uuid"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateProperty(data *property_dto.CreateProperty, creatorID uuid.UUID) (*property_model.PropertyModel, error)
	CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error)
	CheckManageability(id uuid.UUID, userId uuid.UUID) (bool, error)
	CheckOwnership(id uuid.UUID, userId uuid.UUID) (bool, error)
	GetPropertyById(id uuid.UUID) (*property_model.PropertyModel, error)
	GetPropertiesByIds(ids []uuid.UUID, fields []string, userId uuid.UUID) ([]property_model.PropertyModel, error)
	GetUnitsOfProperty(id uuid.UUID) ([]unit_model.UnitModel, error)
	GetListingsOfProperty(id uuid.UUID, query *listing_dto.GetListingsOfPropertyQuery) ([]listing_model.ListingModel, error)
	GetApplicationsOfProperty(id uuid.UUID, query *application_dto.GetApplicationsOfPropertyQuery) ([]application_model.ApplicationModel, error)
	GetManagedProperties(userId uuid.UUID, fields []string) ([]GetManagedPropertiesItem, error)
	SearchListingCombination(data *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error)
	UpdateProperty(data *property_dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error
}

type service struct {
	pRepo property_repo.Repo
	uRepo unit_repo.Repo
	lRepo listing_repo.Repo
	aRepo application_repo.Repo
}

func NewService(pRepo property_repo.Repo, uRepo unit_repo.Repo, lRepo listing_repo.Repo, aRepo application_repo.Repo) Service {
	return &service{
		pRepo: pRepo,
		uRepo: uRepo,
		lRepo: lRepo,
		aRepo: aRepo,
	}
}

func (s *service) CreateProperty(data *property_dto.CreateProperty, creatorID uuid.UUID) (*property_model.PropertyModel, error) {
	data.CreatorID = creatorID
	foundCreator := false
	for _, m := range data.Managers {
		if m.ManagerID == creatorID {
			foundCreator = true
			// m.Role = "OWNER"
			break
		}
	}
	if !foundCreator {
		data.Managers = append(data.Managers, property_dto.CreatePropertyManager{
			ManagerID: creatorID,
			Role:      "OWNER", // TODO: add role to user
		})
	}
	return s.pRepo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*property_model.PropertyModel, error) {
	return s.pRepo.GetPropertyById(context.Background(), id)
}

func (s *service) GetPropertiesByIds(ids []uuid.UUID, fields []string, userId uuid.UUID) ([]property_model.PropertyModel, error) {
	var _ids []string
	for _, id := range ids {
		isVisible, err := s.CheckVisibility(id, userId)
		if err != nil {
			return nil, err
		}
		if isVisible {
			_ids = append(_ids, id.String())
		}
	}
	return s.pRepo.GetPropertiesByIds(context.Background(), _ids, fields)
}

func (s *service) GetUnitsOfProperty(id uuid.UUID) ([]unit_model.UnitModel, error) {
	return s.uRepo.GetUnitsOfProperty(context.Background(), id)
}

func (s *service) GetListingsOfProperty(id uuid.UUID, query *listing_dto.GetListingsOfPropertyQuery) ([]listing_model.ListingModel, error) {
	ids, err := s.pRepo.GetListingsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	idsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idsStr = append(idsStr, id.String())
	}
	return s.lRepo.GetListingsByIds(context.Background(), idsStr, query.Fields)
}

func (s *service) GetApplicationsOfProperty(id uuid.UUID, query *application_dto.GetApplicationsOfPropertyQuery) ([]application_model.ApplicationModel, error) {
	ids, err := s.pRepo.GetApplicationsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	return s.aRepo.GetApplicationsByIds(context.Background(), ids, query.Fields)
}

var (
	ErrMissingPrimaryImage = fmt.Errorf("missing primary image")
	ErrMissingImage        = fmt.Errorf("empty media")
)

func (s *service) UpdateProperty(data *property_dto.UpdateProperty) error {
	if data.Media != nil {
		// primaryImageUrl must exists
		if data.PrimaryImageUrl == nil {
			return ErrMissingPrimaryImage
		}
		// make sure that there is at least one image
		images := []property_model.PropertyMediaModel{}
		data.PrimaryImage = nil
		for i, m := range data.Media {
			if m.Type == database.MEDIATYPEIMAGE {
				images = append(images, m)
				if m.Url == *data.PrimaryImageUrl {
					data.PrimaryImage = types.Ptr(int64(i))
				}
			}
		}
		if len(images) == 0 {
			return ErrMissingImage
		}
	}

	return s.pRepo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckManageability(pid uuid.UUID, userId uuid.UUID) (bool, error) {
	managers, err := s.pRepo.GetPropertyManagers(context.Background(), pid)
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

func (s *service) CheckOwnership(pid uuid.UUID, userId uuid.UUID) (bool, error) {
	managers, err := s.pRepo.GetPropertyManagers(context.Background(), pid)
	if err != nil {
		return false, err
	}
	for _, manager := range managers {
		if manager.ManagerID == userId && manager.Role == "OWNER" {
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

type GetManagedPropertiesItem struct {
	Role     string                       `json:"role"`
	Property property_model.PropertyModel `json:"property"`
}

func (s *service) GetManagedProperties(userId uuid.UUID, fields []string) ([]GetManagedPropertiesItem, error) {
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

	var res []GetManagedPropertiesItem
	for _, p := range managedProps {
		r := GetManagedPropertiesItem{Role: p.Role}
		for _, pp := range ps {
			if pp.ID == p.PropertyID {
				r.Property = pp
			}
		}
		res = append(res, r)
	}

	return res, nil
}

func (s *service) SearchListingCombination(data *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error) {
	data.SortBy = types.Ptr(utils.PtrDerefence[string](data.SortBy, "created_at"))
	data.Order = types.Ptr(utils.PtrDerefence[string](data.Order, "desc"))
	data.Limit = types.Ptr(utils.PtrDerefence[int32](data.Limit, 1000))
	data.Offset = types.Ptr(utils.PtrDerefence[int32](data.Offset, 0))
	return s.pRepo.SearchPropertyCombination(context.Background(), data)
}
