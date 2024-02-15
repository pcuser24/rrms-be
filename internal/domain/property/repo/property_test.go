package repo

import (
	"cmp"
	"context"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestCreateProperty(t *testing.T) {
	// case 0: all fields are valid
	NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)
	// case 2: an error causes rollback
}

func TestGetPropertyManagers(t *testing.T) {
	p := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)

	managers, err := testPropertyRepo.GetPropertyManagers(context.Background(), p.ID)
	require.NoError(t, err)
	require.Equal(t, len(p.Managers), len(managers))
	pmaCmp := func(a, b model.PropertyManagerModel) int {
		return cmp.Compare[string](a.ManagerID.String(), b.ManagerID.String())
	}
	slices.SortFunc(p.Managers, pmaCmp)
	slices.SortFunc(managers, pmaCmp)
	require.Equal(t, p.Managers, managers)

	require.Equal(t, len(managers), len(p.Managers))
	for i := 0; i < len(managers); i++ {
		require.Equal(t, managers[i].PropertyID, p.ID)
		require.Equal(t, managers[i].ManagerID, p.Managers[i].ManagerID)
		require.Equal(t, managers[i].Role, p.Managers[i].Role)
	}
}

func TestGetPropertyById(t *testing.T) {
	p := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)

	p_1, err := testPropertyRepo.GetPropertyById(context.Background(), p.ID)
	require.NoError(t, err)
	sameProperties(t, p, p_1)
}

func TestGetPropertyByIds(t *testing.T) {
	selectedFields := random.RandomlyPickNFromSlice[string](dto.GetRetrievableFields(), 5)
	p1 := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)
	p2 := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)

	ps, err := testPropertyRepo.GetPropertiesByIds(
		context.Background(),
		[]string{p1.ID.String(), p2.ID.String()},
		selectedFields,
	)
	require.NoError(t, err)
	require.Equal(t, len(ps), 2)

	compareFn := func(p_1, p_2 *model.PropertyModel) {
		for _, f := range selectedFields {
			switch f {
			case "name":
				require.Equal(t, p_1.Name, p_2.Name)
			case "area":
				require.Equal(t, p_1.Area, p_2.Area)
			case "number_of_floors":
				require.Equal(t, *p_1.NumberOfFloors, *p_2.NumberOfFloors)
			case "full_address":
				require.Equal(t, p_1.FullAddress, p_2.FullAddress)
			case "city":
				require.Equal(t, p_1.City, p_2.City)
			case "district":
				require.Equal(t, p_1.District, p_2.District)
			case "primary_image":
				require.Equal(t, p_1.PrimaryImage, p_2.PrimaryImage)
			case "type":
				require.Equal(t, p_1.Type, p_2.Type)
			case "is_public":
				require.Equal(t, p_1.IsPublic, p_2.IsPublic)
			case "created_at":
				require.WithinDuration(t, p_1.CreatedAt, p_2.CreatedAt, time.Second)
			case "updated_at":
				require.WithinDuration(t, p_1.UpdatedAt, p_2.UpdatedAt, time.Second)
			case "building":
				require.Equal(t, *p_1.Building, *p_2.Building)
			case "project":
				require.Equal(t, *p_1.Project, *p_2.Project)
			case "year_built":
				require.Equal(t, *p_1.YearBuilt, *p_2.YearBuilt)
			case "orientation":
				require.Equal(t, *p_1.Orientation, *p_2.Orientation)
			case "entrance_width":
				require.Equal(t, *p_1.EntranceWidth, *p_2.EntranceWidth)
			case "facade":
				require.Equal(t, *p_1.Facade, *p_2.Facade)
			case "ward":
				require.Equal(t, *p_1.Ward, *p_2.Ward)
			case "lat":
				require.Equal(t, *p_1.Lat, *p_2.Lat)
			case "lng":
				require.Equal(t, *p_1.Lng, *p_2.Lng)
			case "description":
				require.Equal(t, *p_1.Description, *p_2.Description)
			case "features":
				require.Equal(t, len(p_1.Features), len(p_2.Features))
			case "tags":
				require.Equal(t, len(p_1.Tags), len(p_2.Tags))
			case "media":
				require.Equal(t, len(p_1.Media), len(p_2.Media))
			}
		}
	}
	for _, p := range ps {
		if p.ID == p1.ID {
			compareFn(&p, p1)
		} else if p.ID == p2.ID {
			compareFn(&p, p2)
		} else {
			t.Error("unexpected property", p)
		}
	}
}

func TestPublicity(t *testing.T) {
	p := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)
	require.False(t, p.IsPublic)

	err := testPropertyRepo.UpdateProperty(context.Background(), &dto.UpdateProperty{
		ID:       p.ID,
		IsPublic: types.Ptr[bool](true),
	})
	require.NoError(t, err)

	isPublic, err := testPropertyRepo.IsPublic(context.Background(), p.ID)
	require.NoError(t, err)
	require.True(t, isPublic)
}

func TestUpdateProperty(t *testing.T) {
	p := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)

	// case 0: all fields are valid
	arg := dto.UpdateProperty{
		ID:             p.ID,
		Name:           types.Ptr[string](random.RandomAlphanumericStr(100)),
		Building:       types.Ptr[string](random.RandomAlphanumericStr(100)),
		Project:        types.Ptr[string](random.RandomAlphanumericStr(100)),
		Area:           types.Ptr[float32](random.RandomFloat32(10, 50)),
		NumberOfFloors: types.Ptr[int32](random.RandomInt32(1, 10)),
		YearBuilt:      types.Ptr[int32](random.RandomInt32(1990, 2020)),
		Orientation:    types.Ptr[string](orientations[random.RandomInt32(0, int32(len(orientations)-1))]),
		EntranceWidth:  types.Ptr[float32](random.RandomFloat32(1, 10)),
		Facade:         types.Ptr[float32](random.RandomFloat32(1, 10)),
		FullAddress:    types.Ptr[string](random.RandomAlphanumericStr(100)),
		District:       types.Ptr[string](random.RandomAlphanumericStr(100)),
		City:           types.Ptr[string](random.RandomAlphanumericStr(100)),
		Ward:           types.Ptr[string](random.RandomAlphanumericStr(100)),
		Lat:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Lng:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		PrimaryImage:   types.Ptr[int64](random.RandomInt64(1, 10)),
		Description:    types.Ptr[string](random.RandomAlphanumericStr(100)),
		IsPublic:       utils.Ternary(p.IsPublic, types.Ptr[bool](false), types.Ptr[bool](true)).(*bool),
	}
	err := testPropertyRepo.UpdateProperty(context.Background(), &arg)
	require.NoError(t, err)

	p1, err := testPropertyRepo.GetPropertyById(context.Background(), p.ID)
	require.NoError(t, err)
	require.Equal(t, *arg.Name, p1.Name)
	require.Equal(t, *arg.Building, *p1.Building)
	require.Equal(t, *arg.Project, *p1.Project)
	require.Equal(t, *arg.Area, p1.Area)
	require.Equal(t, *arg.NumberOfFloors, *p1.NumberOfFloors)
	require.Equal(t, *arg.YearBuilt, *p1.YearBuilt)
	require.Equal(t, *arg.Orientation, *p1.Orientation)
	require.Equal(t, *arg.EntranceWidth, *p1.EntranceWidth)
	require.Equal(t, *arg.Facade, *p1.Facade)
	require.Equal(t, *arg.FullAddress, p1.FullAddress)
	require.Equal(t, *arg.District, p1.District)
	require.Equal(t, *arg.City, p1.City)
	require.Equal(t, *arg.Ward, *p1.Ward)
	require.Equal(t, *arg.PrimaryImage, p1.PrimaryImage)
	require.Equal(t, *arg.Lat, *p1.Lat)
	require.Equal(t, *arg.Lng, *p1.Lng)
	require.Equal(t, *arg.Description, *p1.Description)
	require.Equal(t, *arg.IsPublic, p1.IsPublic)
}

func TestDeleteProperty(t *testing.T) {
	p := NewRandomPropertyDB(t, testPropertyRepo, testAuthRepo)

	err := testPropertyRepo.DeleteProperty(context.Background(), p.ID)
	require.NoError(t, err)

	p1, err := testPropertyRepo.GetPropertyById(context.Background(), p.ID)
	require.Error(t, err)
	require.Equal(t, err, database.ErrRecordNotFound)
	require.Empty(t, p1)
}
