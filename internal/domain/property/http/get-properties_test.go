package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"go.uber.org/mock/gomock"
)

func TestGetPropertyById(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")

	property := property_repo.NewRandomPropertyModel(t, userId)

	testcases := []struct {
		name          string
		id            string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			id:   property.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(property.ID)).
					Times(1).
					Return(property.IsPublic, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
				pRepo.EXPECT().
					GetPropertyById(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)
				requireBodyMatchProperty(t, res.Body, property)
			},
		},
		{
			name:      "BadRequest/InvalidId",
			id:        "abcd",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "Forbidden/PrivateProperty",
			id:   property.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					uuid.Nil, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(property.ID)).
					Times(1).
					Return(property.IsPublic, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusForbidden, res.StatusCode)
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
			lRepo := listing_repo.NewMockRepo(ctrl)
			aRepo := application_repo.NewMockRepo(ctrl)
			rRepo := rental_repo.NewMockRepo(ctrl)
			authRepo := auth_repo.NewMockRepo(ctrl)

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo, rRepo, authRepo)

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/properties/property/%s", tc.id),
				nil)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestGetPropertyByIds(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")

	properties := []*model.PropertyModel{
		property_repo.NewRandomPropertyModel(t, userId),
		property_repo.NewRandomPropertyModel(t, userId),
	}

	testcases := []struct {
		name          string
		query         fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			query: fiber.Map{
				"propIds": fmt.Sprintf("propIds=%s&propIds=%s", properties[0].ID.String(), properties[1].ID.String()),
				"fields":  fmt.Sprintf("fields=%s", strings.Join(dto.GetRetrievableFields(), ",")),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(properties[0].ID)).
					Times(1).
					Return(properties[0].IsPublic, nil)
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(properties[1].ID)).
					Times(1).
					Return(properties[1].IsPublic, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(properties[0].ID)).
					Times(1).
					Return(properties[0].Managers, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(properties[1].ID)).
					Times(1).
					Return(properties[1].Managers, nil)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Eq([]string{properties[0].ID.String(), properties[1].ID.String()}), gomock.Eq(dto.GetRetrievableFields())).
					Times(1).
					Return([]model.PropertyModel{*properties[0], *properties[1]}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)
			},
		},
		{
			name: "OK/CustomFields",
			query: fiber.Map{
				"propIds": fmt.Sprintf("propIds=%s&propIds=%s", properties[0].ID.String(), properties[1].ID.String()),
				"fields": func() string {
					fs := dto.GetRetrievableFields()
					return fmt.Sprintf("fields=%s", strings.Join(append(
						fs[:4],
						fs[len(fs)-3:]...,
					), ","))
				}(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(properties[0].ID)).
					Times(1).
					Return(properties[0].IsPublic, nil)
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Eq(userId), gomock.Eq(properties[1].ID)).
					Times(1).
					Return(properties[1].IsPublic, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(properties[0].ID)).
					Times(1).
					Return(properties[0].Managers, nil)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(properties[1].ID)).
					Times(1).
					Return(properties[1].Managers, nil)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Eq([]string{properties[0].ID.String(), properties[1].ID.String()}), gomock.Eq(func() []string {
						fs := dto.GetRetrievableFields()
						return append(
							fs[:4],
							fs[len(fs)-3:]...,
						)
					}())).
					Times(1).
					Return([]model.PropertyModel{
						{ID: properties[0].ID, Name: properties[0].Name, Building: properties[0].Building, Project: properties[0].Project, Area: properties[0].Area, Features: properties[0].Features, Tags: properties[0].Tags, Media: properties[0].Media},
						{ID: properties[1].ID, Name: properties[1].Name, Building: properties[1].Building, Project: properties[1].Project, Area: properties[1].Area, Features: properties[1].Features, Tags: properties[1].Tags, Media: properties[1].Media},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var gotProperties []model.PropertyModel
				err = json.Unmarshal(data, &gotProperties)
				require.NoError(t, err)
				match := func(t *testing.T, p1, p2 *model.PropertyModel) {
					require.Equal(t, p1.Name, p2.Name)
					require.Equal(t, *p1.Building, *p2.Building)
					require.Equal(t, *p1.Project, *p2.Project)
					require.Equal(t, p1.Area, p2.Area)
					require.Equal(t, p1.Features, p2.Features)
					require.Equal(t, p1.Tags, p2.Tags)
					require.Equal(t, p1.Media, p2.Media)
				}
				if gotProperties[0].ID == properties[0].ID {
					match(t, &gotProperties[0], properties[0])
					match(t, &gotProperties[1], properties[1])
				} else {
					match(t, &gotProperties[1], properties[0])
					match(t, &gotProperties[0], properties[1])
				}
			},
		},
		{
			name: "BadRequest/InvalidId",
			query: fiber.Map{
				"propIds": strings.Join([]string{properties[0].ID.String(), "20354d7a-e4fe-47af-8ff6-187bcef92f3fa"}, ","),
				"fields":  "abc,name",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "BadRequest/InvalidFields",
			query: fiber.Map{
				"propIds": fmt.Sprintf("propIds=%s&propIds=%s", properties[0].ID.String(), properties[1].ID.String()),
				"fields":  "fields=abc,name",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					IsPropertyVisible(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
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
			lRepo := listing_repo.NewMockRepo(ctrl)
			aRepo := application_repo.NewMockRepo(ctrl)
			rRepo := rental_repo.NewMockRepo(ctrl)
			authRepo := auth_repo.NewMockRepo(ctrl)

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo, rRepo, authRepo)

			var queries []string
			for _, v := range tc.query {
				queries = append(queries, v.(string))
			}

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/properties/ids/?%s", strings.Join(queries, "&")),
				nil)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestGetManagedProperties(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")

	properties := []*model.PropertyModel{
		property_repo.NewRandomPropertyModel(t, userId),
		property_repo.NewRandomPropertyModel(t, userId),
	}
	properties[0].Managers[0].Role = "OWNER"
	properties[1].Managers[0].Role = "MANAGER"

	testcases := []struct {
		name          string
		query         fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			query: fiber.Map{
				"fields": strings.Join(dto.GetRetrievableFields(), ","),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetManagedProperties(gomock.Any(), gomock.Eq(userId), gomock.Eq(&dto.GetPropertiesQuery{})).
					Times(1).
					Return([]property_repo.GetManagedPropertiesRow{
						{
							PropertyID: properties[0].Managers[0].PropertyID,
							Role:       properties[0].Managers[0].Role,
						},
						{
							PropertyID: properties[1].Managers[0].PropertyID,
							Role:       properties[1].Managers[0].Role,
						},
					}, nil)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Eq([]string{properties[0].ID.String(), properties[1].ID.String()}), gomock.Eq(dto.GetRetrievableFields())).
					Times(1).
					Return([]model.PropertyModel{*properties[0], *properties[1]}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name: "OK/CustomFields",
			query: fiber.Map{
				"fields": strings.Join(append(
					dto.GetRetrievableFields()[:4],
					dto.GetRetrievableFields()[len(dto.GetRetrievableFields())-3:]...,
				), ","),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetManagedProperties(gomock.Any(), gomock.Eq(userId), gomock.Eq(&dto.GetPropertiesQuery{})).
					Times(1).
					Return([]property_repo.GetManagedPropertiesRow{
						{
							PropertyID: properties[0].Managers[0].PropertyID,
							Role:       properties[0].Managers[0].Role,
						},
						{
							PropertyID: properties[1].Managers[0].PropertyID,
							Role:       properties[1].Managers[0].Role,
						},
					}, nil)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Eq([]string{properties[0].ID.String(), properties[1].ID.String()}), gomock.Eq(append(
						dto.GetRetrievableFields()[:4],
						dto.GetRetrievableFields()[len(dto.GetRetrievableFields())-3:]...,
					))).
					Times(1).
					Return([]model.PropertyModel{
						{ID: properties[0].ID, Name: properties[0].Name, Building: properties[0].Building, Project: properties[0].Project, Area: properties[0].Area, Features: properties[0].Features, Tags: properties[0].Tags, Media: properties[0].Media},
						{ID: properties[1].ID, Name: properties[1].Name, Building: properties[1].Building, Project: properties[1].Project, Area: properties[1].Area, Features: properties[1].Features, Tags: properties[1].Tags, Media: properties[1].Media},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var gotProperties []property_service.GetManagedPropertiesItem
				err = json.Unmarshal(data, &gotProperties)
				require.NoError(t, err)
				match := func(t *testing.T, p1, p2 property_service.GetManagedPropertiesItem) {
					require.Equal(t, p1.Property.Name, p2.Property.Name)
					require.Equal(t, p1.Property.Building, p2.Property.Building)
					require.Equal(t, p1.Property.Project, p2.Property.Project)
					require.Equal(t, p1.Property.Area, p2.Property.Area)
					require.Equal(t, p1.Property.Features, p2.Property.Features)
					require.Equal(t, p1.Property.Tags, p2.Property.Tags)
					require.Equal(t, p1.Property.Media, p2.Property.Media)
				}
				if gotProperties[0].Property.ID == properties[0].ID {
					match(t, gotProperties[0], property_service.GetManagedPropertiesItem{
						Role:     properties[0].Managers[0].Role,
						Property: *properties[0],
					})
					match(t, gotProperties[1], property_service.GetManagedPropertiesItem{
						Role:     properties[1].Managers[0].Role,
						Property: *properties[1],
					})
				} else {
					match(t, gotProperties[0], property_service.GetManagedPropertiesItem{
						Role:     properties[1].Managers[0].Role,
						Property: *properties[1],
					})
					match(t, gotProperties[1], property_service.GetManagedPropertiesItem{
						Role:     properties[0].Managers[0].Role,
						Property: *properties[0],
					})
				}
			},
		},
		{
			name:      "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetManagedProperties(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					GetPropertiesByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
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
			lRepo := listing_repo.NewMockRepo(ctrl)
			aRepo := application_repo.NewMockRepo(ctrl)
			rRepo := rental_repo.NewMockRepo(ctrl)
			authRepo := auth_repo.NewMockRepo(ctrl)

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo, rRepo, authRepo)

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/properties/my-properties?fields=%s", tc.query["fields"]),
				nil)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}
