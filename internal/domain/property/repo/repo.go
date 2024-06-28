package repo

import (
	"context"

	"github.com/google/uuid"
	application_dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/redisd"
)

type Repo interface {
	CreateProperty(ctx context.Context, data *property_dto.CreateProperty) (*property_model.PropertyModel, error)
	GetPropertyManagers(ctx context.Context, id uuid.UUID) ([]property_model.PropertyManagerModel, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (*property_model.PropertyModel, error)
	GetPropertiesByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]property_model.PropertyModel, error) // Get properties with custom fields by ids
	GetManagedProperties(ctx context.Context, userId uuid.UUID, query *property_dto.GetPropertiesQuery) ([]GetManagedPropertiesRow, error)
	GetListingsOfProperty(ctx context.Context, id uuid.UUID, query *listing_dto.GetListingsOfPropertyQuery) ([]uuid.UUID, error)
	GetApplicationsOfProperty(ctx context.Context, id uuid.UUID, query *application_dto.GetApplicationsOfPropertyQuery) ([]int64, error)
	GetRentalsOfProperty(ctx context.Context, id uuid.UUID, query *rental_dto.GetRentalsOfPropertyQuery) ([]int64, error)
	SearchPropertyCombination(ctx context.Context, query *property_dto.SearchPropertyCombinationQuery) (*property_dto.SearchPropertyCombinationResponse, error)
	IsPropertyVisible(ctx context.Context, uid, pid uuid.UUID) (bool, error)
	FilterVisibleProperties(ctx context.Context, pids []uuid.UUID, uid uuid.UUID) ([]uuid.UUID, error)
	UpdateProperty(ctx context.Context, data *property_dto.UpdateProperty) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error

	CreatePropertyManagerRequest(ctx context.Context, data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error)
	GetNewPropertyManagerRequest(ctx context.Context, id int64) (property_model.NewPropertyManagerRequest, error)
	GetNewPropertyManagerRequestsToUser(ctx context.Context, uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error)
	UpdatePropertyManagerRequest(ctx context.Context, id int64, uid uuid.UUID, approved bool) error

	CreatePropertyVerificationRequest(ctx context.Context, data *property_dto.CreatePropertyVerificationRequest) (property_model.PropertyVerificationRequest, error)
	GetPropertyVerificationRequest(ctx context.Context, id int64) (property_model.PropertyVerificationRequest, error)
	GetPropertyVerificationRequests(ctx context.Context, filter *property_dto.GetPropertyVerificationRequestsQuery) (*property_dto.GetPropertyVerificationRequestsResponse, error)
	GetPropertiesVerificationStatus(ctx context.Context, ids []uuid.UUID) ([]property_dto.GetPropertyVerificationStatus, error)
	GetPropertyVerificationRequestsOfProperty(ctx context.Context, pid uuid.UUID, limit, offset int32) ([]property_model.PropertyVerificationRequest, error)
	UpdatePropertyVerificationRequestStatus(ctx context.Context, id int64, data *property_dto.UpdatePropertyVerificationRequestStatus) error
}

type repo struct {
	dao         database.DAO
	redisClient redisd.RedisClient
}

func NewRepo(d database.DAO, redisClient redisd.RedisClient) Repo {
	return &repo{
		dao:         d,
		redisClient: redisClient,
	}
}
