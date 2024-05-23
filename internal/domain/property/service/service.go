package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	"github.com/google/uuid"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
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
	GetManagedProperties(userId uuid.UUID, query *property_dto.GetPropertiesQuery) (int, []GetManagedPropertiesItem, error)
	SearchListingCombination(data *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error)
	UpdateProperty(data *property_dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error

	CreatePropertyManagerRequest(data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error)
	GetRentalsOfProperty(id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]rental_model.RentalModel, error)
	GetNewPropertyManagerRequestsToUser(uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error)
	UpdatePropertyManagerRequest(pid, uid uuid.UUID, requestId int64, approved bool) error
}

type service struct {
	pRepo    property_repo.Repo
	uRepo    unit_repo.Repo
	lRepo    listing_repo.Repo
	aRepo    application_repo.Repo
	rRepo    rental_repo.Repo
	authRepo auth_repo.Repo
}

func NewService(pRepo property_repo.Repo, uRepo unit_repo.Repo, lRepo listing_repo.Repo, aRepo application_repo.Repo, rRepo rental_repo.Repo, authRepo auth_repo.Repo) Service {
	return &service{
		pRepo:    pRepo,
		uRepo:    uRepo,
		lRepo:    lRepo,
		aRepo:    aRepo,
		rRepo:    rRepo,
		authRepo: authRepo,
	}
}

func (s *service) CreateProperty(data *property_dto.CreateProperty, creatorID uuid.UUID) (*property_model.PropertyModel, error) {
	data.CreatorID = creatorID
	data.Managers = []property_dto.CreatePropertyManager{
		{
			ManagerID: creatorID,
			Role:      "OWNER",
		},
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

func (s *service) GetRentalsOfProperty(id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]rental_model.RentalModel, error) {
	ids, err := s.pRepo.GetRentalsOfProperty(context.Background(), id, query)
	if err != nil {
		return nil, err
	}

	return s.rRepo.GetRentalsByIds(context.Background(), ids, query.Fields)
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
	return s.pRepo.IsPropertyVisible(context.Background(), uid, id)
}

func (s *service) DeleteProperty(id uuid.UUID) error {
	return s.pRepo.DeleteProperty(context.Background(), id)
}

type GetManagedPropertiesItem struct {
	Role     string                       `json:"role"`
	Property property_model.PropertyModel `json:"property"`
	// Active listings of the property
	Listings []uuid.UUID `json:"listings"`
	// Active rentals of the property
	Rentals []int64 `json:"rentals"`
}

func (s *service) GetManagedProperties(userId uuid.UUID, query *property_dto.GetPropertiesQuery) (int, []GetManagedPropertiesItem, error) {
	var _query property_dto.GetPropertiesQuery = *query
	_query.Limit = types.Ptr[int32](math.MaxInt32)
	managedProps, err := s.pRepo.GetManagedProperties(context.Background(), userId, &_query)
	if err != nil {
		return 0, nil, err
	}

	total := len(managedProps)
	var actualLength int
	if query.Limit == nil {
		actualLength = total
	} else {
		actualLength = utils.Ternary(total > int(*query.Limit), int(*query.Limit), total)
	}
	pids := make([]string, 0, actualLength)
	for _, p := range managedProps[0:actualLength] {
		pid := p.PropertyID.String()
		pids = append(pids, pid)
	}

	ps, err := s.pRepo.GetPropertiesByIds(context.Background(), pids, query.Fields)
	if err != nil {
		return total, nil, err
	}

	res := make([]GetManagedPropertiesItem, 0, actualLength)
	for _, p := range managedProps[0:actualLength] {
		r := GetManagedPropertiesItem{Role: p.Role}
		for _, pp := range ps {
			if pp.ID == p.PropertyID {
				r.Property = pp
			}
		}
		// get active listings of the property
		r.Listings, err = s.pRepo.GetListingsOfProperty(
			context.Background(), p.PropertyID, &listing_dto.GetListingsOfPropertyQuery{
				Expired: false,
				Offset:  types.Ptr(int32(0)),
				Limit:   types.Ptr(int32(1000)),
			})
		if err != nil {
			return total, nil, err
		}
		// get active rentals
		r.Rentals, err = s.pRepo.GetRentalsOfProperty(context.Background(), p.PropertyID, &rental_dto.GetRentalsOfPropertyQuery{
			Expired: false,
			Offset:  types.Ptr(int32(0)),
			Limit:   types.Ptr(int32(1000)),
		})
		if err != nil {
			return total, nil, err
		}
		res = append(res, r)
	}

	return total, res, nil
}

func (s *service) SearchListingCombination(q *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error) {
	if len(q.SortBy) == 0 {
		q.SortBy = append(q.SortBy, "properties.created_at")
		q.Order = append(q.Order, "desc")
	}
	q.Limit = types.Ptr(utils.PtrDerefence(q.Limit, 1000))
	q.Offset = types.Ptr(utils.PtrDerefence(q.Offset, 0))
	return s.pRepo.SearchPropertyCombination(context.Background(), q)
}

var ErrUserIsAlreadyManager = errors.New("user is already a manager of the property")

func (s *service) CreatePropertyManagerRequest(data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error) {
	managers, err := s.pRepo.GetPropertyManagers(context.Background(), data.PropertyID)
	if err != nil {
		return property_model.NewPropertyManagerRequest{}, err
	}
	if exists, err := func() (bool, error) {
		for _, manager := range managers {
			user, err := s.authRepo.GetUserById(context.Background(), manager.ManagerID)
			if err != nil {
				return false, err
			}
			if user.Email == data.Email {
				return true, nil
			}
		}
		return false, nil
	}(); exists || err != nil {
		if err != nil {
			return property_model.NewPropertyManagerRequest{}, err
		}
		return property_model.NewPropertyManagerRequest{}, ErrUserIsAlreadyManager
	}

	user, err := s.authRepo.GetUserByEmail(context.Background(), data.Email)
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return property_model.NewPropertyManagerRequest{}, err
	} else {
		data.UserID = user.ID
	}

	return s.pRepo.CreatePropertyManagerRequest(context.Background(), data)

	// TODO: send email to user, push notification if user is already registered
}

func (s *service) GetNewPropertyManagerRequestsToUser(uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error) {
	return s.pRepo.GetNewPropertyManagerRequestsToUser(context.Background(), uid, limit, offset)
}

var ErrUpdateRequestInfoMismatch = errors.New("request update info mismatch")

func (s *service) UpdatePropertyManagerRequest(pid, uid uuid.UUID, requestId int64, approved bool) error {
	user, err := s.authRepo.GetUserById(context.Background(), uid)
	if err != nil {
		return err
	}
	request, err := s.pRepo.GetNewPropertyManagerRequest(context.Background(), requestId)
	if err != nil {
		return err
	}
	if (request.UserID != uuid.Nil && uid != request.UserID) ||
		request.PropertyID != pid ||
		user.Email != request.Email {
		return ErrUpdateRequestInfoMismatch
	}

	return s.pRepo.UpdatePropertyManagerRequest(context.Background(), requestId, user.ID, approved)
}