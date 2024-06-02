package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	contactTypes  = []string{"OWNER", "MANAGER", "BROKER"}
	postDurations = []int{7, 15, 30}
)
var (
	testingUser     *auth_model.UserModel
	testingProperty *property_model.PropertyModel
	testingUnits    []*unit_model.UnitModel
)

func PrepareRandomListing(
	t *testing.T,
	aRepo auth_repo.Repo, pRepo property_repo.Repo, uRrepo unit_repo.Repo,
	propertyId uuid.UUID, unitIds []uuid.UUID,
) dto.CreateListing {
	if aRepo != nil && pRepo != nil && uRrepo != nil {
		testingUser = auth_repo.NewRandomUserDB(t, aRepo)

		createProperty := property_repo.PrepareRandomProperty(t, aRepo, testingUser.ID)
		testingProperty = property_repo.NewRandomPropertyDBFromArg(t, pRepo, &createProperty)
		propertyId = testingProperty.ID

		for i := 0; i < 3; i++ {
			createUnit := unit_repo.PrepareRandomUnit(t, aRepo, pRepo, testingProperty.ID)
			newUnit := unit_repo.NewRandomUnitDBFromArg(t, uRrepo, &createUnit)
			testingUnits = append(testingUnits, newUnit)
			unitIds = append(unitIds, newUnit.ID)
		}
	}

	return dto.CreateListing{
		CreatorID:  testingUser.ID,
		PropertyID: propertyId,
		Units: []dto.CreateListingUnit{
			{
				UnitID: unitIds[0],
				Price:  random.RandomInt64(1000000, 10000000),
			},
			{
				UnitID: unitIds[1],
				Price:  random.RandomInt64(1000000, 10000000),
			},
			{
				UnitID: unitIds[2],
				Price:  random.RandomInt64(1000000, 10000000),
			},
		},
		Title:             random.RandomAlphanumericStr(100),
		Description:       random.RandomAlphanumericStr(1000),
		FullName:          testingUser.FirstName + " " + testingUser.LastName,
		Email:             testingUser.Email,
		Phone:             random.RandomNumericStr(10),
		ContactType:       contactTypes[random.RandomInt32(0, 2)],
		Price:             random.RandomFloat32(1000000, 10000000),
		PriceNegotiable:   random.RandomInt32(0, 10)%2 == 0,
		SecurityDeposit:   types.Ptr[float32](random.RandomFloat32(1000000, 10000000)),
		LeaseTerm:         types.Ptr[int32](random.RandomInt32(12, 48)),
		PetsAllowed:       types.Ptr[bool](true),
		NumberOfResidents: types.Ptr[int32](random.RandomInt32(1, 10)),
		Priority:          random.RandomInt32(1, 5),
		PostDuration:      postDurations[random.RandomInt32(0, 2)],
		Policies: []dto.CreateListingPolicy{
			{
				PolicyID: 1,
				Note:     types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				PolicyID: 2,
				Note:     types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
	}
}

func compareListingAndCreateDto(t *testing.T, l *model.ListingModel, c *dto.CreateListing) {
	require.NotEmpty(t, l)
	require.NotEmpty(t, c)
	require.Equal(t, c.CreatorID, l.CreatorID)
	require.Equal(t, c.PropertyID, l.PropertyID)
	require.Equal(t, c.Title, l.Title)
	require.Equal(t, c.Description, l.Description)
	require.Equal(t, c.FullName, l.FullName)
	require.Equal(t, c.Email, l.Email)
	require.Equal(t, c.Phone, l.Phone)
	require.Equal(t, c.ContactType, l.ContactType)
	require.Equal(t, c.Price, l.Price)
	require.Equal(t, c.PriceNegotiable, l.PriceNegotiable)
	require.Equal(t, c.SecurityDeposit, l.SecurityDeposit)
	require.Equal(t, c.LeaseTerm, l.LeaseTerm)
	require.Equal(t, c.PetsAllowed, l.PetsAllowed)
	require.Equal(t, c.NumberOfResidents, l.NumberOfResidents)
	require.Equal(t, c.Priority, l.Priority)
	require.WithinDuration(t, l.CreatedAt, time.Now(), time.Second)
	require.WithinDuration(t, l.UpdatedAt, time.Now(), time.Second)
	require.WithinDuration(t, time.Now().AddDate(0, 0, c.PostDuration), l.ExpiredAt, time.Second)
	require.Equal(t, len(c.Policies), len(l.Policies))
	require.Equal(t, len(c.Units), len(l.Units))
}

func NewRandomListingDB(
	t *testing.T,
	aRepo auth_repo.Repo, pRepo property_repo.Repo, uRrepo unit_repo.Repo, lRepo Repo,
) *model.ListingModel {
	arg := PrepareRandomListing(t, aRepo, pRepo, uRrepo, uuid.Nil, []uuid.UUID{})

	l, err := lRepo.CreateListing(context.Background(), &arg)
	require.NoError(t, err)
	compareListingAndCreateDto(t, l, &arg)

	return l
}

func NewRandomListingModel(
	t *testing.T,
	userId, propertyId uuid.UUID, unitIds []uuid.UUID,
) *model.ListingModel {
	id := uuid.MustParse("978fe220-663f-464c-9afd-fe05de7be44b")
	units := make([]model.ListingUnitModel, len(unitIds))
	for i, unitId := range unitIds {
		units[i] = model.ListingUnitModel{
			ListingID: id,
			UnitID:    unitId,
			Price:     random.RandomInt64(1000000, 10000000),
		}
	}

	return &model.ListingModel{
		ID:                id,
		CreatorID:         userId,
		PropertyID:        propertyId,
		Units:             units,
		Title:             random.RandomAlphanumericStr(100),
		Description:       random.RandomAlphanumericStr(1000),
		FullName:          random.RandomAlphabetStr(20),
		Email:             random.RandomEmail(),
		Phone:             random.RandomNumericStr(10),
		ContactType:       contactTypes[random.RandomInt32(0, 2)],
		Price:             random.RandomFloat32(1000000, 10000000),
		PriceNegotiable:   random.RandomInt32(0, 10)%2 == 0,
		SecurityDeposit:   types.Ptr[float32](random.RandomFloat32(1000000, 10000000)),
		LeaseTerm:         types.Ptr[int32](random.RandomInt32(12, 48)),
		PetsAllowed:       types.Ptr[bool](true),
		NumberOfResidents: types.Ptr[int32](random.RandomInt32(1, 10)),
		Priority:          random.RandomInt32(1, 5),
		ExpiredAt:         time.Now().AddDate(0, 0, 15),
		Policies: []model.ListingPolicyModel{
			{
				ListingID: id,
				PolicyID:  1,
				Note:      types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				ListingID: id,
				PolicyID:  6,
				Note:      types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				ListingID: id,
				PolicyID:  7,
				Note:      types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
	}
}

func NewRandomListingDBFromArg(
	t *testing.T,
	testListingRepo Repo,
	arg *dto.CreateListing,
) *model.ListingModel {

	p, err := testListingRepo.CreateListing(context.Background(), arg)
	require.NoError(t, err)
	compareListingAndCreateDto(t, p, arg)

	return p
}
