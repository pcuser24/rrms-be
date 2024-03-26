package repo

import (
	"cmp"
	"context"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	rentalIntentions = []string{"RESIDE", "BUSINESS"}
	employmentStatus = []string{"EMPLOYED", "UNEMPLOYED", "STUDENT", "RETIRED"}
	identityTypes    = []string{"ID", "PASSPORT", "DRIVER_LICENSE"}
)

var (
	testingUser     *auth_model.UserModel
	testingProperty *property_model.PropertyModel
	testingListing  *listing_model.ListingModel
	testingUnits    []*unit_model.UnitModel
)

func PrepareRandomApplication(
	t *testing.T,
	aRepo auth_repo.Repo, lRepo listing_repo.Repo, pRepo property_repo.Repo, uRrepo unit_repo.Repo,
	userId, listingId, propertyId uuid.UUID, unitIds []uuid.UUID,
) dto.CreateApplication {
	if aRepo != nil && pRepo != nil && uRrepo != nil {
		testingUser = auth_repo.NewRandomUserDB(t, aRepo)
		userId = testingUser.ID

		createProperty := property_repo.PrepareRandomProperty(t, aRepo, testingUser.ID)
		testingProperty = property_repo.NewRandomPropertyDBFromArg(t, pRepo, &createProperty)
		propertyId = testingProperty.ID

		for i := 0; i < 3; i++ {
			createUnit := unit_repo.PrepareRandomUnit(t, aRepo, pRepo, testingProperty.ID)
			newUnit := unit_repo.NewRandomUnitDBFromArg(t, uRrepo, &createUnit)
			testingUnits = append(testingUnits, newUnit)
			unitIds = append(unitIds, newUnit.ID)
		}

		createListing := listing_repo.PrepareRandomListing(t, aRepo, pRepo, uRrepo, testingProperty.ID, unitIds)
		testingListing = listing_repo.NewRandomListingDBFromArg(t, lRepo, &createListing)
		listingId = testingListing.ID
	}

	ret := dto.CreateApplication{
		ListingID:               listingId,
		PropertyID:              propertyId,
		CreatorID:               userId,
		FullName:                testingUser.FirstName + " " + testingUser.LastName,
		Email:                   testingUser.Email,
		Phone:                   random.RandomNumericStr(10),
		Dob:                     random.RandomDate(),
		ProfileImage:            random.RandomURL(),
		MoveinDate:              random.RandomDate(),
		PreferredTerm:           random.RandomInt32(12, 36),
		RentalIntention:         rentalIntentions[random.RandomInt32(0, int32(len(rentalIntentions)-1))],
		RhAddress:               types.Ptr(random.RandomAlphanumericStr(100)),
		RhCity:                  types.Ptr(random.RandomCity()),
		RhDistrict:              types.Ptr(random.RandomDistrict()),
		RhWard:                  types.Ptr(random.RandomWard()),
		RhRentalDuration:        types.Ptr(random.RandomInt32(12, 36)),
		RhMonthlyPayment:        types.Ptr(random.RandomInt64(1000000, 10000000)),
		RhReasonForLeaving:      types.Ptr(random.RandomAlphabetStr(100)),
		EmploymentStatus:        employmentStatus[random.RandomInt32(0, int32(len(employmentStatus)-1))],
		EmploymentCompanyName:   types.Ptr(random.RandomAlphabetStr(50)),
		EmploymentPosition:      types.Ptr(random.RandomAlphabetStr(50)),
		EmploymentMonthlyIncome: types.Ptr(random.RandomInt64(1000000, 10000000)),
		EmploymentComment:       types.Ptr(random.RandomAlphabetStr(100)),
		Minors: []dto.CreateApplicationMinor{
			{
				FullName:    random.RandomAlphabetStr(50),
				Dob:         random.RandomDate(),
				Email:       types.Ptr(random.RandomEmail()),
				Phone:       types.Ptr(random.RandomNumericStr(10)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				FullName:    random.RandomAlphabetStr(50),
				Dob:         random.RandomDate(),
				Email:       types.Ptr(random.RandomEmail()),
				Phone:       types.Ptr(random.RandomNumericStr(10)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				FullName:    random.RandomAlphabetStr(50),
				Dob:         random.RandomDate(),
				Email:       types.Ptr(random.RandomEmail()),
				Phone:       types.Ptr(random.RandomNumericStr(10)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
		Coaps: []dto.CreateApplicationCoapModel{
			{
				FullName:    random.RandomAlphabetStr(50),
				Dob:         random.RandomDate(),
				Email:       types.Ptr(random.RandomEmail()),
				Phone:       types.Ptr(random.RandomNumericStr(10)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
				Job:         random.RandomAlphabetStr(10),
				Income:      random.RandomInt32(1000000, 10000000),
			},
			{
				FullName:    random.RandomAlphabetStr(50),
				Dob:         random.RandomDate(),
				Email:       types.Ptr(random.RandomEmail()),
				Phone:       types.Ptr(random.RandomNumericStr(10)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
				Job:         random.RandomAlphabetStr(10),
				Income:      random.RandomInt32(1000000, 10000000),
			},
		},
		Pets: []dto.CreateApplicationPet{
			{
				Type:        random.RandomAlphabetStr(10),
				Weight:      types.Ptr(random.RandomFloat32(1, 100)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				Type:        random.RandomAlphabetStr(10),
				Weight:      types.Ptr(random.RandomFloat32(1, 100)),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
		Vehicles: []dto.CreateApplicationVehicle{
			{
				Type:        random.RandomAlphabetStr(10),
				Model:       types.Ptr(random.RandomAlphabetStr(10)),
				Code:        random.RandomAlphabetStr(10),
				Description: types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
	}

	return ret
}

func sameApplications(t *testing.T, a1, a2 *model.ApplicationModel) {
	require.NotEmpty(t, a1)
	require.NotEmpty(t, a2)
	require.Equal(t, a1.ListingID, a2.ListingID)
	require.Equal(t, a1.PropertyID, a2.PropertyID)
	require.Equal(t, a1.CreatorID, a2.CreatorID)
	require.Equal(t, a1.FullName, a2.FullName)
	require.Equal(t, a1.Email, a2.Email)
	require.Equal(t, a1.Phone, a2.Phone)
	require.Equal(t, a1.Dob, a2.Dob)
	require.Equal(t, a1.ProfileImage, a2.ProfileImage)
	require.Equal(t, a1.MoveinDate, a2.MoveinDate)
	require.Equal(t, a1.PreferredTerm, a2.PreferredTerm)
	require.Equal(t, a1.RentalIntention, a2.RentalIntention)
	require.Equal(t, a1.RhAddress, a2.RhAddress)
	require.Equal(t, a1.RhCity, a2.RhCity)
	require.Equal(t, a1.RhDistrict, a2.RhDistrict)
	require.Equal(t, a1.RhWard, a2.RhWard)
	require.Equal(t, a1.RhRentalDuration, a2.RhRentalDuration)
	require.Equal(t, a1.RhMonthlyPayment, a2.RhMonthlyPayment)
	require.Equal(t, a1.RhReasonForLeaving, a2.RhReasonForLeaving)
	require.Equal(t, a1.EmploymentStatus, a2.EmploymentStatus)
	require.Equal(t, a1.EmploymentCompanyName, a2.EmploymentCompanyName)
	require.Equal(t, a1.EmploymentPosition, a2.EmploymentPosition)
	require.Equal(t, a1.EmploymentMonthlyIncome, a2.EmploymentMonthlyIncome)
	require.Equal(t, a1.EmploymentComment, a2.EmploymentComment)

	require.Equal(t, len(a1.Minors), len(a2.Minors))
	pmCmp := func(m1, m2 model.ApplicationMinorModel) int {
		return cmp.Compare(m1.FullName, m2.FullName)
	}
	slices.SortFunc(a1.Minors, pmCmp)
	slices.SortFunc(a2.Minors, pmCmp)
	require.Equal(t, a1.Minors, a2.Minors)

	require.Equal(t, len(a1.Coaps), len(a2.Coaps))
	pcCmp := func(m1, m2 model.ApplicationCoapModel) int {
		return cmp.Compare(m1.FullName, m2.FullName)
	}
	slices.SortFunc(a1.Coaps, pcCmp)
	slices.SortFunc(a2.Coaps, pcCmp)
	require.Equal(t, a1.Coaps, a2.Coaps)

	require.Equal(t, len(a1.Pets), len(a2.Pets))
	ppCmp := func(m1, m2 model.ApplicationPetModel) int {
		return cmp.Compare(m1.Type, m2.Type)
	}
	slices.SortFunc(a1.Pets, ppCmp)
	slices.SortFunc(a2.Pets, ppCmp)

	require.Equal(t, len(a1.Vehicles), len(a2.Vehicles))
	pvCmp := func(m1, m2 model.ApplicationVehicle) int {
		return cmp.Compare(m1.Type, m2.Type)
	}
	slices.SortFunc(a1.Vehicles, pvCmp)
	slices.SortFunc(a2.Vehicles, pvCmp)
}

func compareApplicationAndCreateDto(t *testing.T, a *model.ApplicationModel, arg *dto.CreateApplication) {

}

func NewRandomApplicationDB(
	t *testing.T,
	aRepo auth_repo.Repo, lRepo listing_repo.Repo, pRepo property_repo.Repo, uRrepo unit_repo.Repo,
	appRepo Repo,
) *model.ApplicationModel {
	arg := PrepareRandomApplication(t, aRepo, lRepo, pRepo, uRrepo, uuid.Nil, uuid.Nil, uuid.Nil, []uuid.UUID{})

	a, err := appRepo.CreateApplication(context.Background(), &arg)
	require.NoError(t, err)
	compareApplicationAndCreateDto(t, a, &arg)

	return a
}

func NewRandomApplicationModel(
	t *testing.T,
	userId, listingId, propertyId uuid.UUID, unitIds []uuid.UUID,
) *model.ApplicationModel {
	id := random.RandomInt64(1, 100)

	ret := &model.ApplicationModel{
		ID:                      id,
		ListingID:               listingId,
		PropertyID:              propertyId,
		CreatorID:               userId,
		FullName:                testingUser.FirstName + " " + testingUser.LastName,
		Email:                   testingUser.Email,
		Phone:                   random.RandomNumericStr(10),
		Dob:                     random.RandomDate(),
		ProfileImage:            random.RandomURL(),
		MoveinDate:              random.RandomDate(),
		PreferredTerm:           random.RandomInt32(12, 36),
		RentalIntention:         rentalIntentions[random.RandomInt32(0, int32(len(rentalIntentions)-1))],
		RhAddress:               types.Ptr(random.RandomAlphanumericStr(100)),
		RhCity:                  types.Ptr(random.RandomCity()),
		RhDistrict:              types.Ptr(random.RandomDistrict()),
		RhWard:                  types.Ptr(random.RandomWard()),
		RhRentalDuration:        types.Ptr(random.RandomInt32(12, 36)),
		RhMonthlyPayment:        types.Ptr(random.RandomInt64(1000000, 10000000)),
		RhReasonForLeaving:      types.Ptr(random.RandomAlphabetStr(100)),
		EmploymentStatus:        employmentStatus[random.RandomInt32(0, int32(len(employmentStatus)-1))],
		EmploymentCompanyName:   types.Ptr(random.RandomAlphabetStr(50)),
		EmploymentPosition:      types.Ptr(random.RandomAlphabetStr(50)),
		EmploymentMonthlyIncome: types.Ptr(random.RandomInt64(1000000, 10000000)),
		EmploymentComment:       types.Ptr(random.RandomAlphabetStr(100)),
		Minors: []model.ApplicationMinorModel{
			{
				ApplicationID: id,
				FullName:      random.RandomAlphabetStr(50),
				Dob:           random.RandomDate(),
				Email:         types.Ptr(random.RandomEmail()),
				Phone:         types.Ptr(random.RandomNumericStr(10)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				ApplicationID: id,
				FullName:      random.RandomAlphabetStr(50),
				Dob:           random.RandomDate(),
				Email:         types.Ptr(random.RandomEmail()),
				Phone:         types.Ptr(random.RandomNumericStr(10)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				ApplicationID: id,
				FullName:      random.RandomAlphabetStr(50),
				Dob:           random.RandomDate(),
				Email:         types.Ptr(random.RandomEmail()),
				Phone:         types.Ptr(random.RandomNumericStr(10)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
		Coaps: []model.ApplicationCoapModel{
			{
				ApplicationID: id,
				FullName:      random.RandomAlphabetStr(50),
				Dob:           random.RandomDate(),
				Email:         types.Ptr(random.RandomEmail()),
				Phone:         types.Ptr(random.RandomNumericStr(10)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
				Job:           random.RandomAlphabetStr(10),
				Income:        random.RandomInt32(1000000, 10000000),
			},
			{
				ApplicationID: id,
				FullName:      random.RandomAlphabetStr(50),
				Dob:           random.RandomDate(),
				Email:         types.Ptr(random.RandomEmail()),
				Phone:         types.Ptr(random.RandomNumericStr(10)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
				Job:           random.RandomAlphabetStr(10),
				Income:        random.RandomInt32(1000000, 10000000),
			},
		},
		Pets: []model.ApplicationPetModel{
			{
				ApplicationID: id,
				Type:          random.RandomAlphabetStr(9),
				Weight:        types.Ptr(random.RandomFloat32(1, 100)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
			{
				ApplicationID: id,
				Type:          random.RandomAlphabetStr(9),
				Weight:        types.Ptr(random.RandomFloat32(1, 100)),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
		Vehicles: []model.ApplicationVehicle{
			{
				ApplicationID: id,
				Type:          random.RandomAlphabetStr(9),
				Model:         types.Ptr(random.RandomAlphabetStr(10)),
				Code:          random.RandomAlphabetStr(10),
				Description:   types.Ptr(random.RandomAlphabetStr(100)),
			},
		},
	}

	return ret
}

func NewRandomApplicationFromArg(
	t *testing.T,
	aRepo Repo,
	arg *dto.CreateApplication,
) *model.ApplicationModel {
	a, err := aRepo.CreateApplication(context.Background(), arg)
	require.NoError(t, err)
	compareApplicationAndCreateDto(t, a, arg)

	return a
}
