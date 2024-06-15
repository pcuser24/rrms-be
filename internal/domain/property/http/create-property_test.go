package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"

	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"

	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
	"go.uber.org/mock/gomock"
)

func TestCreateProperty(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")

	property := property_repo.NewRandomPropertyModel(t, userId)
	defaultArg := propertyModelToCreatePropertyArg(property)

	testcases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			body: createPropertyToMap(defaultArg),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.AnyOf(defaultArg)).
					Times(1).
					Return(property, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)

				if !assert.Equal(t, http.StatusCreated, res.StatusCode) {
					data, err := io.ReadAll(res.Body)
					require.NoError(t, err)
					var errorMessage validation.ErrorMessage
					err = json.Unmarshal(data, &errorMessage)
					require.NoError(t, err)
					t.Log(errorMessage)
				}

				require.Equal(t, http.StatusCreated, res.StatusCode)
				requireBodyMatchProperty(t, res.Body, property)
			},
		},
		{
			name: "OK/MissingOptionalFields",
			body: func() fiber.Map {
				m := createPropertyToMap(defaultArg)
				delete(m, "building")
				delete(m, "project")
				delete(m, "numberOfFloors")
				delete(m, "yearBuilt")
				delete(m, "orientation")
				delete(m, "entranceWidth")
				delete(m, "facade")
				delete(m, "ward")
				delete(m, "lat")
				delete(m, "lng")
				delete(m, "description")
				return m
			}(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				arg := defaultArg
				arg.Building = nil
				arg.Project = nil
				arg.NumberOfFloors = nil
				arg.YearBuilt = nil
				arg.Orientation = nil
				arg.EntranceWidth = nil
				arg.Facade = nil
				arg.Ward = nil
				arg.Lat = nil
				arg.Lng = nil
				arg.Description = nil

				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.AnyOf(arg)).
					Times(1).
					Return(property, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)

				if !assert.Equal(t, http.StatusCreated, res.StatusCode) {
					data, err := io.ReadAll(res.Body)
					require.NoError(t, err)
					var errorMessage validation.ErrorMessage
					err = json.Unmarshal(data, &errorMessage)
					require.NoError(t, err)
					t.Log(errorMessage)
				}

				require.Equal(t, http.StatusCreated, res.StatusCode)
				requireBodyMatchProperty(t, res.Body, property)
			},
		},
		{
			name: "BadRequest/MissingRequiredFields",
			body: func() fiber.Map {
				m := createPropertyToMap(defaultArg)
				delete(m, "name")
				return m
			}(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var errorMessage validation.ErrorMessage
				err = json.Unmarshal(data, &errorMessage)
				require.NoError(t, err)
				t.Log("error message:", errorMessage)
				require.Equal(t, "[Name]: '' | Needs to implement 'required'", errorMessage.Message)
			},
		},
		{
			name: "BadRequest/InvalidFieldType",
			body: func() fiber.Map {
				m := createPropertyToMap(defaultArg)
				m["numberOfFloors"] = "abc"
				return m
			}(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var errorMessage validation.ErrorMessage
				err = json.Unmarshal(data, &errorMessage)
				require.NoError(t, err)
				t.Log("error message:", errorMessage)
			},
		},
		{
			name: "BadRequest/DeepValidation/MissingRequiredField",
			body: func() fiber.Map {
				m := createPropertyToMap(defaultArg)
				m["features"] = []struct {
					Description *string `json:"description"`
				}{
					{Description: &dto.GetRetrievableFields()[0]},
				}
				return m
			}(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var errorMessage validation.ErrorMessage
				err = json.Unmarshal(data, &errorMessage)
				require.NoError(t, err)
				t.Log("error message:", errorMessage)
				require.Equal(t, "[FeatureID]: '0' | Needs to implement 'required'", errorMessage.Message)
			},
		},
		{
			name: "BadRequest/DeepValidation/InvalidFieldType",
			body: func() fiber.Map {
				m := createPropertyToMap(defaultArg)
				m["features"] = []struct {
					FeatureID   string `json:"featureId"`
					Description int    `json:"description"`
				}{
					{
						FeatureID:   "abc",
						Description: 12,
					},
				}
				return m
			}(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					CreateProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var errorMessage validation.ErrorMessage
				err = json.Unmarshal(data, &errorMessage)
				require.NoError(t, err)
				t.Log("error message:", errorMessage)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			pRepo := property_repo.NewMockRepo(ctrl)
			uRepo := unit_repo.NewMockRepo(ctrl)

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, ctrl)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/properties/", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func requireBodyMatchProperty(t *testing.T, body io.ReadCloser, property *model.PropertyModel) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProperty model.PropertyModel
	err = json.Unmarshal(data, &gotProperty)
	require.NoError(t, err)
	require.Equal(t, *property, gotProperty)
}

func propertyModelToCreatePropertyArg(property *model.PropertyModel) *dto.CreateProperty {
	ret := dto.CreateProperty{
		CreatorID:      property.CreatorID,
		Name:           property.Name,
		Building:       property.Building,
		Project:        property.Project,
		Area:           property.Area,
		NumberOfFloors: property.NumberOfFloors,
		YearBuilt:      property.YearBuilt,
		Orientation:    property.Orientation,
		EntranceWidth:  property.EntranceWidth,
		Facade:         property.Facade,
		FullAddress:    property.FullAddress,
		District:       property.District,
		City:           property.City,
		Ward:           property.Ward,
		PrimaryImage:   property.Media[0].Url,
		Lat:            property.Lat,
		Lng:            property.Lng,
		Description:    property.Description,
		Type:           property.Type,
	}
	for _, m := range property.Managers {
		ret.Managers = append(ret.Managers, dto.CreatePropertyManager{
			ManagerID: m.ManagerID,
			Role:      m.Role,
		})
	}
	for _, m := range property.Media {
		ret.Media = append(ret.Media, dto.CreatePropertyMedia{
			Url:         m.Url,
			Type:        m.Type,
			Description: m.Description,
		})
	}
	for _, m := range property.Features {
		ret.Features = append(ret.Features, dto.CreatePropertyFeature{
			FeatureID:   m.FeatureID,
			Description: m.Description,
		})
	}
	for _, m := range property.Tags {
		ret.Tags = append(ret.Tags, dto.CreatePropertyTag{
			Tag: m.Tag,
		})
	}
	return &ret
}

func createPropertyToMap(arg *dto.CreateProperty) fiber.Map {
	return fiber.Map{
		"creatorId":      arg.CreatorID,
		"name":           arg.Name,
		"building":       arg.Building,
		"project":        arg.Project,
		"area":           arg.Area,
		"numberOfFloors": arg.NumberOfFloors,
		"yearBuilt":      arg.YearBuilt,
		"orientation":    arg.Orientation,
		"entranceWidth":  arg.EntranceWidth,
		"facade":         arg.Facade,
		"fullAddress":    arg.FullAddress,
		"district":       arg.District,
		"city":           arg.City,
		"ward":           arg.Ward,
		"primaryImage":   arg.PrimaryImage,
		"lat":            arg.Lat,
		"lng":            arg.Lng,
		"description":    arg.Description,
		"type":           arg.Type,
		"managers":       arg.Managers,
		"media":          arg.Media,
		"features":       arg.Features,
		"tags":           arg.Tags,
	}
}
