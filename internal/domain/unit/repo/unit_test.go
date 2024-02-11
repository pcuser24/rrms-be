package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_dto "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestCreateUnit(t *testing.T) {
	NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)
}

func TestGetUnitById(t *testing.T) {
	u := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)

	u_1, err := testUnitRepo.GetUnitById(context.Background(), u.ID)
	require.NoError(t, err)
	SameUnits(t, u, u_1)
}

func TestGetUnitByIds(t *testing.T) {
	selectedFields := random.RandomlyPickNFromSlice[string](unit_dto.GetRetrievableFields(), 5)
	u1 := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)
	u2 := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)

	us, err := testUnitRepo.GetUnitsByIds(
		context.Background(),
		[]string{u1.ID.String(), u2.ID.String()},
		selectedFields,
	)
	require.NoError(t, err)
	require.Equal(t, len(us), 2)

	compareFn := func(u_1, u_2 *unit_model.UnitModel) {
		for _, f := range selectedFields {
			switch f {
			case "name":
				require.Equal(t, u_1.Name, u_2.Name)
			case "area":
				require.Equal(t, u_1.Area, u_2.Area)
			case "floor":
				require.Equal(t, u_1.Floor, u_2.Floor)
			case "number_of_living_rooms":
				require.Equal(t, u_1.NumberOfLivingRooms, u_2.NumberOfLivingRooms)
			case "number_of_bedrooms":
				require.Equal(t, u_1.NumberOfBedrooms, u_2.NumberOfBedrooms)
			case "number_of_bathrooms":
				require.Equal(t, u_1.NumberOfBathrooms, u_2.NumberOfBathrooms)
			case "number_of_toilets":
				require.Equal(t, u_1.NumberOfToilets, u_2.NumberOfToilets)
			case "number_of_kitchens":
				require.Equal(t, u_1.NumberOfKitchens, u_2.NumberOfKitchens)
			case "number_of_balconies":
				require.Equal(t, u_1.NumberOfBalconies, u_2.NumberOfBalconies)
			}
		}
	}
	for _, u := range us {
		if u.ID == u1.ID {
			compareFn(&u, u1)
		} else if u.ID == u2.ID {
			compareFn(&u, u2)
		} else {
			t.Errorf("unexpected unit: %v", u)
		}
	}
}

func TestCheckUnitManageability(t *testing.T) {
	// prepare data
	user1 := auth_repo.NewRandomUserDB(t, testAuthRepo)
	user2 := auth_repo.NewRandomUserDB(t, testAuthRepo)

	prop_arg := property_repo.PrepareRandomProperty(t, nil, user1.ID)
	prop_arg.Managers = append(prop_arg.Managers, property_dto.CreatePropertyManager{
		ManagerID: user2.ID,
		Role:      "MANAGER",
	})
	prop := property_repo.NewRandomPropertyDBFromArg(t, testPropertyRepo, &prop_arg)

	unit_arg := PrepareRandomUnit(t, nil, nil, prop.ID)
	unit := NewRandomUnitDBFromArg(t, testUnitRepo, &unit_arg)

	// case 0: user is the creator of the property => should return true
	res, err := testUnitRepo.CheckUnitManageability(context.Background(), unit.ID, user1.ID)
	require.NoError(t, err)
	require.True(t, res)
	// case 1: user is a manager of the property => should return true
	res, err = testUnitRepo.CheckUnitManageability(context.Background(), unit.ID, user2.ID)
	require.NoError(t, err)
	require.True(t, res)
	// case 2: user is not the creator or a manager of the property => should return false
	res, err = testUnitRepo.CheckUnitManageability(context.Background(), unit.ID, uuid.Nil)
	require.NoError(t, err)
	require.False(t, res)
}

func CheckUnitOfProperty(t *testing.T) {
	// prepare data
	unit := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)
	// case 0: unit is of the property => should return true
	res, err := testUnitRepo.CheckUnitOfProperty(context.Background(), unit.PropertyID, unit.ID)
	require.NoError(t, err)
	require.True(t, res)
	// case 1: unit is not of the property => should return false
	res, err = testUnitRepo.CheckUnitOfProperty(context.Background(), uuid.Nil, unit.ID)
	require.NoError(t, err)
	require.False(t, res)
}

func TestPublicity(t *testing.T) {
	// prepare data
	unit := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)
	// case 0: the property is not public by default => should return false
	res, err := testUnitRepo.IsPublic(context.Background(), unit.ID)
	require.NoError(t, err)
	require.False(t, res)
	// case 1: update publicity of the property => should return true
	err = testPropertyRepo.UpdateProperty(context.Background(), &property_dto.UpdateProperty{
		ID:       unit.PropertyID,
		IsPublic: types.Ptr[bool](true),
	})
	require.NoError(t, err)
	res, err = testUnitRepo.IsPublic(context.Background(), unit.ID)
	require.NoError(t, err)
	require.True(t, res)
}

func TestUpdateUnit(t *testing.T) {
	// prepare data
	unit := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)

	// case 0: all fields are valid
	arg := unit_dto.UpdateUnit{
		ID:                  unit.ID,
		Name:                types.Ptr[string]("new name"),
		Area:                types.Ptr[float32](random.RandomFloat32(10, 100)),
		Floor:               types.Ptr[int32](random.RandomInt32(1, 10)),
		NumberOfLivingRooms: types.Ptr[int32](random.RandomInt32(1, 5)),
		NumberOfBedrooms:    types.Ptr[int32](random.RandomInt32(1, 5)),
		NumberOfBathrooms:   types.Ptr[int32](random.RandomInt32(1, 5)),
		NumberOfToilets:     types.Ptr[int32](random.RandomInt32(1, 5)),
		NumberOfKitchens:    types.Ptr[int32](random.RandomInt32(1, 5)),
		NumberOfBalconies:   types.Ptr[int32](random.RandomInt32(1, 5)),
	}
	err := testUnitRepo.UpdateUnit(context.Background(), &arg)
	require.NoError(t, err)

	u, err := testUnitRepo.GetUnitById(context.Background(), unit.ID)
	require.NoError(t, err)
	require.Equal(t, *arg.Name, u.Name)
	require.Equal(t, *arg.Area, u.Area)
	require.Equal(t, *arg.Floor, *u.Floor)
	require.Equal(t, *arg.NumberOfLivingRooms, *u.NumberOfLivingRooms)
	require.Equal(t, *arg.NumberOfBedrooms, *u.NumberOfBedrooms)
	require.Equal(t, *arg.NumberOfBathrooms, *u.NumberOfBathrooms)
	require.Equal(t, *arg.NumberOfToilets, *u.NumberOfToilets)
	require.Equal(t, *arg.NumberOfKitchens, *u.NumberOfKitchens)
	require.Equal(t, *arg.NumberOfBalconies, *u.NumberOfBalconies)

	// case 1: invalid value
	arg.NumberOfBathrooms = types.Ptr[int32](-1)
	err = testUnitRepo.UpdateUnit(context.Background(), &arg)
	require.Error(t, err)
}

func TestDeleteUnit(t *testing.T) {
	// prepare data
	unit := NewRandomUnitDB(t, testUnitRepo, testPropertyRepo, testAuthRepo)

	err := testUnitRepo.DeleteUnit(context.Background(), unit.ID)
	require.NoError(t, err)

	_, err = testUnitRepo.GetUnitById(context.Background(), unit.ID)
	require.Error(t, err)
	require.Equal(t, database.ErrRecordNotFound, err)
}
