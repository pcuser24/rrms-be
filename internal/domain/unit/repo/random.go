package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	unitTypes = []string{"APARTMENT", "ROOM", "STUDIO"}
)

var (
	testingUser     *auth_model.UserModel
	testingProperty *property_model.PropertyModel
)

func PrepareRandomUnit(
	t *testing.T,
	aRepo auth_repo.Repo, pRepo property_repo.Repo,
	propertyId uuid.UUID,
) dto.CreateUnit {
	if aRepo != nil && pRepo != nil {
		testingUser = auth_repo.NewRandomUserDB(t, aRepo)
		createProperty := property_repo.PrepareRandomProperty(t, aRepo, testingUser.ID)
		testingProperty = property_repo.NewRandomPropertyDBFromArg(t, pRepo, &createProperty)
		propertyId = testingProperty.ID
	}

	return dto.CreateUnit{
		PropertyID:          propertyId,
		Name:                types.Ptr[string](random.RandomAlphabetStr(10)),
		Area:                random.RandomFloat32(100, 200),
		Floor:               types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfLivingRooms: types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfBedrooms:    types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfBathrooms:   types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfToilets:     types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfKitchens:    types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfBalconies:   types.Ptr[int32](random.RandomInt32(1, 10)),
		Type:                database.UNITTYPE(unitTypes[random.RandomInt32(0, 1)]),
		Amenities: []dto.CreateUnitAmenity{
			{
				AmenityID:   7,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				AmenityID:   6,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
		Media: []dto.CreateUnitMedia{
			{
				Url:         random.RandomURL(),
				Type:        "IMAGE",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				Url:         random.RandomURL(),
				Type:        "VIDEO",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				Url:         random.RandomURL(),
				Type:        "VIDEO",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				Url:         random.RandomURL(),
				Type:        "IMAGE",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
	}
}

func NewRandomUnitModel(t *testing.T, propertyID uuid.UUID) *model.UnitModel {
	id := uuid.MustParse("1b3a7930-401f-4c6d-8c7f-c5283f5430ad")

	return &model.UnitModel{
		ID:                  id,
		PropertyID:          propertyID,
		Name:                random.RandomAlphabetStr(10),
		Area:                random.RandomFloat32(100, 200),
		Floor:               types.Ptr(random.RandomInt32(1, 10)),
		NumberOfLivingRooms: types.Ptr(random.RandomInt32(1, 10)),
		NumberOfBedrooms:    types.Ptr(random.RandomInt32(1, 10)),
		NumberOfBathrooms:   types.Ptr(random.RandomInt32(1, 10)),
		NumberOfToilets:     types.Ptr(random.RandomInt32(1, 10)),
		NumberOfKitchens:    types.Ptr(random.RandomInt32(1, 10)),
		NumberOfBalconies:   types.Ptr(random.RandomInt32(1, 10)),
		Type:                database.UNITTYPE(unitTypes[random.RandomInt32(0, 2)]),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		Amenities: []model.UnitAmenityModel{
			{
				UnitID:      id,
				AmenityID:   7,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				UnitID:      id,
				AmenityID:   9,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				UnitID:      id,
				AmenityID:   8,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
		Media: []model.UnitMediaModel{
			{
				UnitID:      id,
				Url:         random.RandomURL(),
				Type:        "IMAGE",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				UnitID:      id,
				Url:         random.RandomURL(),
				Type:        "VIDEO",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
	}
}

func SameUnits(t *testing.T, u1, u2 *model.UnitModel) {
	require.NotEmpty(t, u1)
	require.NotEmpty(t, u2)
	require.Equal(t, u1.ID, u2.ID)
	require.Equal(t, u1.PropertyID, u2.PropertyID)
	require.Equal(t, u1.Name, u2.Name)
	require.Equal(t, u1.Area, u2.Area)
	require.Equal(t, u1.Floor, u2.Floor)
	require.Equal(t, u1.NumberOfLivingRooms, u2.NumberOfLivingRooms)
	require.Equal(t, u1.NumberOfBedrooms, u2.NumberOfBedrooms)
	require.Equal(t, u1.NumberOfBathrooms, u2.NumberOfBathrooms)
	require.Equal(t, u1.NumberOfToilets, u2.NumberOfToilets)
	require.Equal(t, u1.NumberOfKitchens, u2.NumberOfKitchens)
	require.Equal(t, u1.NumberOfBalconies, u2.NumberOfBalconies)
	require.Equal(t, u1.Type, u2.Type)
	require.WithinDuration(t, u1.CreatedAt, u2.CreatedAt, time.Second)
	require.WithinDuration(t, u1.UpdatedAt, u2.UpdatedAt, time.Second)
	require.ElementsMatch(t, u1.Amenities, u2.Amenities)
	require.ElementsMatch(t, u1.Media, u2.Media)
}

func compareUnitAndCreateDto(t *testing.T, u *model.UnitModel, arg *dto.CreateUnit) {
	require.NotEmpty(t, u)
	require.Equal(t, arg.PropertyID, u.PropertyID)
	require.Equal(t, *arg.Name, u.Name)
	require.Equal(t, arg.Area, u.Area)
	require.Equal(t, *arg.Floor, *u.Floor)
	require.Equal(t, *arg.NumberOfLivingRooms, *u.NumberOfLivingRooms)
	require.Equal(t, *arg.NumberOfBedrooms, *u.NumberOfBedrooms)
	require.Equal(t, *arg.NumberOfBathrooms, *u.NumberOfBathrooms)
	require.Equal(t, *arg.NumberOfToilets, *u.NumberOfToilets)
	require.Equal(t, *arg.NumberOfKitchens, *u.NumberOfKitchens)
	require.Equal(t, *arg.NumberOfBalconies, *u.NumberOfBalconies)
	require.Equal(t, arg.Type, u.Type)
	require.NotEmpty(t, u.CreatedAt)
	require.NotEmpty(t, u.UpdatedAt)
	var uAmenities []dto.CreateUnitAmenity
	for i := range u.Amenities {
		uAmenities = append(uAmenities, dto.CreateUnitAmenity{
			AmenityID:   u.Amenities[i].AmenityID,
			Description: u.Amenities[i].Description,
		})
	}
	require.ElementsMatch(t, arg.Amenities, uAmenities)
	var uMedia []dto.CreateUnitMedia
	for i := range u.Media {
		uMedia = append(uMedia, dto.CreateUnitMedia{
			Url:         u.Media[i].Url,
			Type:        u.Media[i].Type,
			Description: u.Media[i].Description,
		})
	}
	require.ElementsMatch(t, arg.Media, uMedia)
}

func NewRandomUnitDB(
	t *testing.T,
	uRepo Repo, pRepo property_repo.Repo, aRepo auth_repo.Repo,
) *model.UnitModel {
	arg := PrepareRandomUnit(t, aRepo, pRepo, uuid.Nil)

	u, err := uRepo.CreateUnit(context.Background(), &arg)
	require.NoError(t, err)
	compareUnitAndCreateDto(t, u, &arg)

	return u
}

func NewRandomUnitDBFromArg(
	t *testing.T,
	uRepo Repo,
	arg *dto.CreateUnit,
) *model.UnitModel {

	u, err := uRepo.CreateUnit(context.Background(), arg)
	require.NoError(t, err)
	compareUnitAndCreateDto(t, u, arg)

	return u
}
