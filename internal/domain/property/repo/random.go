package repo

import (
	"cmp"
	"context"
	"slices"
	"testing"

	"github.com/google/uuid"
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
	propertyTypes = []string{"APARTMENT", "PRIVATE", "ROOM", "STORE", "OFFICE", "MINIAPARTMENT"}
	orientations  = []string{"se", "sw", "ne", "nw", "e", "w", "n", "s"}
	roles         = []string{"OWNER", "MANAGER"}
)

var testingUsers []*auth_model.UserModel = make([]*auth_model.UserModel, 0, 3)

func PrepareRandomProperty(
	t *testing.T,
	testAuthRepo auth_repo.Repo,
	creatorId uuid.UUID,
) dto.CreateProperty {
	if testAuthRepo != nil {
		for len(testingUsers) < 3 {
			user := auth_repo.NewRandomUserDB(t, testAuthRepo)
			testingUsers = append(testingUsers, user)
		}
		creatorId = testingUsers[0].ID
	}

	primaryImageUrl := random.RandomURL()
	ret := dto.CreateProperty{
		CreatorID:      creatorId,
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
		PrimaryImage:   primaryImageUrl,
		Lat:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Lng:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Description:    types.Ptr[string](random.RandomAlphanumericStr(100)),
		Type:           database.PROPERTYTYPE(propertyTypes[random.RandomInt32(0, int32(len(propertyTypes)-1))]),
		Media: []dto.CreatePropertyMedia{
			{
				Url:         primaryImageUrl,
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

	if testAuthRepo != nil {
		ret.Managers = []dto.CreatePropertyManager{
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
		}
	} else {
		ret.Managers = []dto.CreatePropertyManager{
			{
				ManagerID: creatorId,
				Role:      "OWNER",
			},
		}
	}

	return ret
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
	require.Equal(t, p1.PrimaryImage, p2.PrimaryImage)
	require.Equal(t, *p1.Lat, *p2.Lat)
	require.Equal(t, *p1.Lng, *p2.Lng)
	require.Equal(t, *p1.Description, *p2.Description)
	require.Equal(t, p1.Type, p2.Type)
	require.Equal(t, len(p1.Managers), len(p2.Managers))
	require.Equal(t, len(p1.Media), len(p2.Media))
	require.Equal(t, len(p1.Features), len(p2.Features))
	require.Equal(t, len(p1.Tags), len(p2.Tags))

	require.Equal(t, len(p1.Managers), len(p2.Managers))
	pmaCmp := func(a, b model.PropertyManagerModel) int {
		return cmp.Compare[string](a.ManagerID.String(), b.ManagerID.String())
	}
	slices.SortFunc(p1.Managers, pmaCmp)
	slices.SortFunc(p2.Managers, pmaCmp)
	require.Equal(t, p1.Managers, p2.Managers)

	require.Equal(t, len(p1.Media), len(p2.Media))
	pmCmp := func(a, b model.PropertyMediaModel) int {
		return int(a.ID - b.ID)
	}
	slices.SortFunc(p1.Media, pmCmp)
	slices.SortFunc(p2.Media, pmCmp)
	require.Equal(t, p1.Media, p2.Media)

	require.Equal(t, len(p1.Features), len(p2.Features))
	pfCmp := func(a, b model.PropertyFeatureModel) int {
		return cmp.Compare[int64](a.FeatureID, b.FeatureID)
	}
	slices.SortFunc(p1.Features, pfCmp)
	slices.SortFunc(p2.Features, pfCmp)
	require.Equal(t, p1.Features, p2.Features)
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
	require.Equal(t, arg.PrimaryImage, func(mediaId int64) string {
		for _, m := range p.Media {
			if m.ID == mediaId {
				return m.Url
			}
		}
		return ""
	}(p.PrimaryImage))
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

func NewRandomPropertyDB(
	t *testing.T,
	testPropertyRepo Repo,
	testAuthRepo auth_repo.Repo,
) *model.PropertyModel {
	arg := PrepareRandomProperty(t, testAuthRepo, uuid.Nil)

	p, err := testPropertyRepo.CreateProperty(context.Background(), &arg)
	require.NoError(t, err)
	comparePropertyAndCreateDto(t, p, &arg)

	return p
}

func NewRandomPropertyModel(t *testing.T, creatorId uuid.UUID) *model.PropertyModel {
	id, err := uuid.NewRandom()
	require.NoError(t, err)

	return &model.PropertyModel{
		ID:             id,
		CreatorID:      creatorId,
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
		PrimaryImage:   1,
		Lat:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Lng:            types.Ptr[float64](random.RandomFloat64(10, 50)),
		Description:    types.Ptr[string](random.RandomAlphanumericStr(100)),
		Type:           database.PROPERTYTYPE(propertyTypes[random.RandomInt32(0, int32(len(propertyTypes)-1))]),
		IsPublic:       false,
		// CreatedAt:      time.Now(),
		// UpdatedAt:      time.Now(),
		Managers: []model.PropertyManagerModel{
			{
				PropertyID: id,
				ManagerID:  creatorId,
				Role:       "OWNER",
			},
		},
		Media: []model.PropertyMediaModel{
			{
				ID:          1,
				PropertyID:  id,
				Url:         random.RandomURL(),
				Type:        "IMAGE",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				ID:          2,
				PropertyID:  id,
				Url:         random.RandomURL(),
				Type:        "VIDEO",
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
		Features: []model.PropertyFeatureModel{
			{
				PropertyID:  id,
				FeatureID:   7,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				PropertyID:  id,
				FeatureID:   9,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
			{
				PropertyID:  id,
				FeatureID:   8,
				Description: types.Ptr[string](random.RandomAlphanumericStr(100)),
			},
		},
		Tags: []model.PropertyTagModel{
			{
				ID:         1,
				PropertyID: id,
				Tag:        random.RandomAlphanumericStr(10),
			},
			{
				ID:         2,
				PropertyID: id,
				Tag:        random.RandomAlphanumericStr(10),
			},
			{
				ID:         3,
				PropertyID: id,
				Tag:        random.RandomAlphanumericStr(10),
			},
		},
	}
}

func NewRandomPropertyDBFromArg(
	t *testing.T,
	testPropertyRepo Repo,
	arg *dto.CreateProperty,
) *model.PropertyModel {

	p, err := testPropertyRepo.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	comparePropertyAndCreateDto(t, p, arg)

	return p
}
