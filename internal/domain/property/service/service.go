package service

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"time"

	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/ds/set"

	"github.com/google/uuid"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

const (
	MAX_IMAGE_SIZE      = 10 * 1024 * 1024 // 10MB
	UPLOAD_URL_LIFETIME = 5                // 5 minutes
)

type Service interface {
	PreCreateProperty(data *property_dto.PreCreateProperty, creatorID uuid.UUID) error
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
	PreUpdateProperty(data *property_dto.PreUpdateProperty, creatorID uuid.UUID) error
	UpdateProperty(data *property_dto.UpdateProperty) error
	DeleteProperty(id uuid.UUID) error

	CreatePropertyManagerRequest(data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error)
	GetRentalsOfProperty(id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]rental_model.RentalModel, error)
	GetNewPropertyManagerRequestsToUser(uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error)
	UpdatePropertyManagerRequest(pid, uid uuid.UUID, requestId int64, approved bool) error

	PreCreatePropertyVerificationRequest(data *property_dto.PreCreatePropertyVerificationRequest, creatorID uuid.UUID) error
	CreatePropertyVerificationRequest(data *property_dto.CreatePropertyVerificationRequest) (property_model.PropertyVerificationRequest, error)
	GetPropertyVerificationRequest(id int64) (property_model.PropertyVerificationRequest, error)
	GetPropertyVerificationRequests(filter *dto.GetPropertyVerificationRequestsQuery) (*property_dto.GetPropertyVerificationRequestsResponse, error)
	GetPropertyVerificationRequestsOfProperty(pid uuid.UUID, limit, offset int32) ([]property_model.PropertyVerificationRequest, error)
	UpdatePropertyVerificationRequestStatus(id int64, data *property_dto.UpdatePropertyVerificationRequestStatus) error
}

type service struct {
	domainRepo repos.DomainRepo

	s3Client        s3.S3Client
	imageBucketName string
}

func NewService(
	domainRepo repos.DomainRepo,
	s3Client s3.S3Client, imageBucketName string,
) Service {
	return &service{
		domainRepo: domainRepo,

		s3Client:        s3Client,
		imageBucketName: imageBucketName,
	}
}

func (s *service) PreCreateProperty(data *property_dto.PreCreateProperty, creatorID uuid.UUID) error {
	for i := range data.Media {
		m := &data.Media[i]
		// split file name and extension
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/properties/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

		url, err := s.s3Client.GetPutObjectPresignedURL(
			s.imageBucketName, objKey, m.Type, m.Size, UPLOAD_URL_LIFETIME*time.Minute,
		)
		if err != nil {
			return err
		}
		m.Url = url.URL
	}
	return nil
}

func (s *service) CreateProperty(data *property_dto.CreateProperty, creatorID uuid.UUID) (*property_model.PropertyModel, error) {
	data.CreatorID = creatorID
	data.Managers = []property_dto.CreatePropertyManager{
		{
			ManagerID: creatorID,
			Role:      "OWNER",
		},
	}
	return s.domainRepo.PropertyRepo.CreateProperty(context.Background(), data)
}

func (s *service) GetPropertyById(id uuid.UUID) (*property_model.PropertyModel, error) {
	return s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), id)
}

func (s *service) GetPropertiesByIds(ids []uuid.UUID, fields []string, uid uuid.UUID) ([]property_model.PropertyModel, error) {
	visibleIDS, err := s.FilterVisibleProperties(ids, uid)
	if err != nil {
		return nil, err
	}
	return s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), visibleIDS, fields)
}

var (
	ErrMissingPrimaryImage = fmt.Errorf("missing primary image")
	ErrMissingImage        = fmt.Errorf("empty media")
)

func (s *service) PreUpdateProperty(data *property_dto.PreUpdateProperty, creatorID uuid.UUID) error {
	for i := range data.Media {
		m := &data.Media[i]
		// split file name and extension
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/properties/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

		url, err := s.s3Client.GetPutObjectPresignedURL(
			s.imageBucketName, objKey, m.Type, m.Size, UPLOAD_URL_LIFETIME*time.Minute,
		)
		if err != nil {
			return err
		}
		m.Url = url.URL
	}
	return nil
}

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

	return s.domainRepo.PropertyRepo.UpdateProperty(context.Background(), data)
}

func (s *service) CheckManageability(pid uuid.UUID, userId uuid.UUID) (bool, error) {
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), pid)
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
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), pid)
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

func (s *service) FilterVisibleProperties(pids []uuid.UUID, uid uuid.UUID) ([]uuid.UUID, error) {
	lidSet := set.NewSet[uuid.UUID]()
	lidSet.AddAll(pids...)
	return s.domainRepo.PropertyRepo.FilterVisibleProperties(context.Background(), lidSet.ToSlice(), uid)
}

func (s *service) CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.domainRepo.PropertyRepo.IsPropertyVisible(context.Background(), uid, id)
}

func (s *service) DeleteProperty(id uuid.UUID) error {
	return s.domainRepo.PropertyRepo.DeleteProperty(context.Background(), id)
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
	managedProps, err := s.domainRepo.PropertyRepo.GetManagedProperties(context.Background(), userId, &_query)
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
	pids := make([]uuid.UUID, 0, actualLength)
	for _, p := range managedProps[0:actualLength] {
		pids = append(pids, p.PropertyID)
	}

	ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), pids, query.Fields)
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
		r.Listings, err = s.domainRepo.PropertyRepo.GetListingsOfProperty(
			context.Background(), p.PropertyID, &listing_dto.GetListingsOfPropertyQuery{
				Expired: false,
				Offset:  types.Ptr(int32(0)),
				Limit:   types.Ptr(int32(1000)),
			})
		if err != nil {
			return total, nil, err
		}
		// get active rentals
		r.Rentals, err = s.domainRepo.PropertyRepo.GetRentalsOfProperty(context.Background(), p.PropertyID, &rental_dto.GetRentalsOfPropertyQuery{
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
	return s.domainRepo.PropertyRepo.SearchPropertyCombination(context.Background(), q)
}
