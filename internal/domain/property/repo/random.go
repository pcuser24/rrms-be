package repo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	propertyTypes = []string{"APARTMENT", "PRIVATE", "TOWNHOUSE", "SHOPHOUSE", "VILLA", "ROOM", "STORE", "OFFICE", "BLOCK", "COMPLEX"}
	orientations  = []string{"se", "sw", "ne", "nw", "e", "w", "n", "s"}
	roles         = []string{"OWNER", "MANAGER"}
)

var testingUsers []*auth_model.UserModel = make([]*auth_model.UserModel, 0, 3)

func prepareRandomProperty(
	t *testing.T,
	testAuthRepo auth_repo.Repo,
) dto.CreateProperty {
	for len(testingUsers) < 3 {
		user := auth_repo.NewRandomUser(t, testAuthRepo)
		testingUsers = append(testingUsers, user)
	}

	return dto.CreateProperty{
		CreatorID:      testingUsers[0].ID,
		Name:           random.RandomAlphabetStr(10),
		Building:       types.Ptr[string](random.RandomAlphabetStr(10)),
		Project:        types.Ptr[string](random.RandomAlphabetStr(10)),
		Area:           random.RandomFloat32(10, 50),
		NumberOfFloors: types.Ptr[int32](random.RandomInt32(1, 10)),
		YearBuilt:      types.Ptr[int32](random.RandomInt32(1990, 2020)),
		Orientation:    types.Ptr[string](orientations[random.RandomInt32(0, int32(len(orientations)-1))]),
		EntranceWidth:  types.Ptr[float32](random.RandomFloat32(1, 10)),
		Facade:         types.Ptr[float32](random.RandomFloat32(1, 10)),
		FullAddress:    random.RandomAddress(),
		District:       random.RandomDistrict(),
		City:           random.RandomCity(),
		Ward:           types.Ptr[string](random.RandomWard()),
		PlaceUrl:       fmt.Sprintf("https://maps.app.goo.gl/%s", random.RandomAlphanumericStr(17)),
		Lat:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Lng:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Description:    types.Ptr[string](random.RandomAlphanumericStr(100)),
		Type:           database.PROPERTYTYPE(propertyTypes[random.RandomInt32(0, int32(len(propertyTypes)-1))]),
		Managers: []dto.CreatePropertyManager{
			{
				ManagerID: testingUsers[0].ID,
				Role:      roles[random.RandomInt32(0, 1)],
			},
			{
				ManagerID: testingUsers[1].ID,
				Role:      roles[random.RandomInt32(0, 1)],
			},
			{
				ManagerID: testingUsers[2].ID,
				Role:      roles[random.RandomInt32(0, 1)],
			},
		},
		Media: []dto.CreatePropertyMedia{
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
		Features: []dto.CreatePropertyFeature{
			{
				FeatureID:   7,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				FeatureID:   9,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				FeatureID:   8,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
		Tags: []dto.CreatePropertyTag{
			{Tag: random.RandomAlphanumericStr(10)},
			{Tag: random.RandomAlphanumericStr(10)},
			{Tag: random.RandomAlphanumericStr(10)},
		},
	}
}

func sameProperties(t *testing.T, p1, p2 *model.PropertyModel) {
	require.NotEmpty(t, p1)
	require.NotEmpty(t, p2)
	require.Equal(t, p2.ID, p1.ID)
	require.Equal(t, p1.CreatorID, p2.CreatorID)
	require.Equal(t, p1.Name, p2.Name)
	require.Equal(t, *p1.Building, *p2.Building)
	require.Equal(t, *p1.Project, *p2.Project)
	require.Equal(t, p1.Area, p2.Area)
	require.Equal(t, *p1.NumberOfFloors, *p2.NumberOfFloors)
	require.Equal(t, *p1.YearBuilt, *p2.YearBuilt)
	require.Equal(t, *p1.Orientation, *p2.Orientation)
	require.Equal(t, *p1.EntranceWidth, *p2.EntranceWidth)
	require.Equal(t, *p1.Facade, *p2.Facade)
	require.Equal(t, p1.FullAddress, p2.FullAddress)
	require.Equal(t, p1.District, p2.District)
	require.Equal(t, p1.City, p2.City)
	require.Equal(t, *p1.Ward, *p2.Ward)
	require.Equal(t, p1.PlaceUrl, p2.PlaceUrl)
	require.Equal(t, *p1.Lat, *p2.Lat)
	require.Equal(t, *p1.Lng, *p2.Lng)
	require.Equal(t, *p1.Description, *p2.Description)
	require.Equal(t, p1.Type, p2.Type)
	require.Equal(t, len(p1.Managers), len(p2.Managers))
	require.Equal(t, len(p1.Media), len(p2.Media))
	require.Equal(t, len(p1.Features), len(p2.Features))
	require.Equal(t, len(p1.Tags), len(p2.Tags))
	for i := 0; i < len(p1.Managers); i++ {
		require.Equal(t, p1.Managers[i].PropertyID, p2.ID)
		require.Equal(t, p1.Managers[i].ManagerID, p2.Managers[i].ManagerID)
		require.Equal(t, p1.Managers[i].Role, p2.Managers[i].Role)
	}
	for i := 0; i < len(p1.Media); i++ {
		require.Equal(t, p1.Media[i].PropertyID, p2.ID)
		require.Equal(t, p1.Media[i].Url, p2.Media[i].Url)
		require.Equal(t, p1.Media[i].Type, p2.Media[i].Type)
		require.Equal(t, *p1.Media[i].Description, *p2.Media[i].Description)
	}
	for i := 0; i < len(p1.Features); i++ {
		require.Equal(t, p1.Features[i].PropertyID, p2.ID)
		require.Equal(t, p1.Features[i].FeatureID, p2.Features[i].FeatureID)
		require.Equal(t, *p1.Features[i].Description, *p2.Features[i].Description)
	}
	for i := 0; i < len(p1.Tags); i++ {
		require.Equal(t, p1.Tags[i].PropertyID, p2.ID)
		require.Equal(t, p1.Tags[i].Tag, p2.Tags[i].Tag)
	}
}

func comparePropertyAndCreateDto(t *testing.T, p *model.PropertyModel, arg *dto.CreateProperty) {
	require.NotEmpty(t, p)
	require.Equal(t, arg.CreatorID, p.CreatorID)
	require.Equal(t, arg.Name, p.Name)
	require.Equal(t, *arg.Building, *p.Building)
	require.Equal(t, *arg.Project, *p.Project)
	require.Equal(t, arg.Area, p.Area)
	require.Equal(t, *arg.NumberOfFloors, *p.NumberOfFloors)
	require.Equal(t, *arg.YearBuilt, *p.YearBuilt)
	require.Equal(t, *arg.Orientation, *p.Orientation)
	require.Equal(t, *arg.EntranceWidth, *p.EntranceWidth)
	require.Equal(t, *arg.Facade, *p.Facade)
	require.Equal(t, arg.FullAddress, p.FullAddress)
	require.Equal(t, arg.District, p.District)
	require.Equal(t, arg.City, p.City)
	require.Equal(t, *arg.Ward, *p.Ward)
	require.Equal(t, arg.PlaceUrl, p.PlaceUrl)
	require.Equal(t, *arg.Lat, *p.Lat)
	require.Equal(t, *arg.Lng, *p.Lng)
	require.Equal(t, *arg.Description, *p.Description)
	require.Equal(t, arg.Type, p.Type)
	require.False(t, p.IsPublic)
	require.Equal(t, len(arg.Managers), len(p.Managers))
	require.Equal(t, len(arg.Media), len(p.Media))
	require.Equal(t, len(arg.Features), len(p.Features))
	require.Equal(t, len(arg.Tags), len(p.Tags))
	for i := 0; i < len(arg.Managers); i++ {
		require.Equal(t, p.Managers[i].PropertyID, p.ID)
		require.Equal(t, arg.Managers[i].ManagerID, p.Managers[i].ManagerID)
		require.Equal(t, arg.Managers[i].Role, p.Managers[i].Role)
	}
	for i := 0; i < len(arg.Media); i++ {
		require.Equal(t, p.Media[i].PropertyID, p.ID)
		require.Equal(t, arg.Media[i].Url, p.Media[i].Url)
		require.Equal(t, arg.Media[i].Type, p.Media[i].Type)
		require.Equal(t, *arg.Media[i].Description, *p.Media[i].Description)
	}
	for i := 0; i < len(arg.Features); i++ {
		require.Equal(t, p.Features[i].PropertyID, p.ID)
		require.Equal(t, arg.Features[i].FeatureID, p.Features[i].FeatureID)
		require.Equal(t, *arg.Features[i].Description, *p.Features[i].Description)
	}
	for i := 0; i < len(arg.Tags); i++ {
		require.Equal(t, p.Tags[i].PropertyID, p.ID)
		require.Equal(t, arg.Tags[i].Tag, p.Tags[i].Tag)
	}
}

func NewRandomProperty(
	t *testing.T,
	testPropertyRepo Repo,
	testAuthRepo auth_repo.Repo,
) *model.PropertyModel {
	arg := prepareRandomProperty(t, testAuthRepo)

	p, err := testPropertyRepo.CreateProperty(context.Background(), &arg)
	require.NoError(t, err)
	comparePropertyAndCreateDto(t, p, &arg)

	return p
}

func newRandomPropertyFromArg(
	t *testing.T,
	testPropertyRepo Repo,
	arg *dto.CreateProperty,
) *model.PropertyModel {

	p, err := testPropertyRepo.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	comparePropertyAndCreateDto(t, p, arg)

	return p
}
