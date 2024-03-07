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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_dto "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"go.uber.org/mock/gomock"
)

func TestGetUnitById(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")
	propertyId := uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a")
	unit := unit_repo.NewRandomUnitModel(t, propertyId)

	testcases := []struct {
		name          string
		id            string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK/UnitIsPublic",
			id:   unit.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(unit.ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					GetUnitById(gomock.Any(), gomock.Eq(unit.ID)).
					Times(1).
					Return(unit, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)
				requireBodyMatchUnit(t, res.Body, unit)
			},
		},
		{
			name: "OK/UnitIsNotPublic",
			id:   unit.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(unit.ID)).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Eq(unit.ID), gomock.Eq(userId)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					GetUnitById(gomock.Any(), gomock.Eq(unit.ID)).
					Times(1).
					Return(unit, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)
				requireBodyMatchUnit(t, res.Body, unit)
			},
		},
		{
			name: "Forbidden/UnitIsNotPublic",
			id:   unit.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(unit.ID)).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Eq(unit.ID), gomock.Eq(uuid.Nil)).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					GetUnitById(gomock.Any(), gomock.Any()).
					Times(0)
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

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo)

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/units/unit/%s", tc.id),
				nil)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestGetUnitsByIds(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")
	propertyId := uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a")
	units := []*unit_model.UnitModel{
		unit_repo.NewRandomUnitModel(t, propertyId),
		unit_repo.NewRandomUnitModel(t, propertyId),
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
				"unitIds": fmt.Sprintf("unitIds=%s&unitIds=%s", units[0].ID.String(), units[1].ID.String()),
				"fields":  fmt.Sprintf("fields=%s", strings.Join(unit_dto.GetRetrievableFields(), ",")),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[0].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[1].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					GetUnitsByIds(gomock.Any(), gomock.Eq([]string{units[0].ID.String(), units[1].ID.String()}), gomock.Eq(unit_dto.GetRetrievableFields())).
					Times(1).
					Return([]unit_model.UnitModel{*units[0], *units[1]}, nil)
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
				"unitIds": fmt.Sprintf("unitIds=%s&unitIds=%s", units[0].ID.String(), units[1].ID.String()),
				"fields": func() string {
					fs := unit_dto.GetRetrievableFields()
					return fmt.Sprintf("fields=%s", strings.Join(append(
						fs[:4],
						fs[len(fs)-2:]...,
					), ","))
				}(),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[0].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[1].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					GetUnitsByIds(gomock.Any(),
						gomock.Eq([]string{units[0].ID.String(), units[1].ID.String()}),
						gomock.Eq(func() []string {
							fs := unit_dto.GetRetrievableFields()
							return append(
								fs[:4],
								fs[len(fs)-2:]...,
							)
						}()),
					).
					Times(1).
					Return([]unit_model.UnitModel{
						{ID: units[0].ID, Name: units[0].Name, PropertyID: units[0].PropertyID, Area: units[0].Area, Floor: units[0].Floor, Media: units[0].Media, Amenities: units[0].Amenities},
						{ID: units[1].ID, Name: units[1].Name, PropertyID: units[1].PropertyID, Area: units[1].Area, Floor: units[1].Floor, Media: units[1].Media, Amenities: units[1].Amenities},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var gotUnits []unit_model.UnitModel
				err = json.Unmarshal(data, &gotUnits)
				require.NoError(t, err)

				matchFn := func(t *testing.T, u1, u2 *unit_model.UnitModel) {
					require.Equal(t, u1.ID, u2.ID)
					require.Equal(t, u1.PropertyID, u2.PropertyID)
					require.Equal(t, u1.Name, u2.Name)
					require.Equal(t, u1.Area, u2.Area)
					require.Equal(t, u1.Floor, u2.Floor)
					require.ElementsMatch(t, u1.Amenities, u2.Amenities)
					require.ElementsMatch(t, u1.Media, u2.Media)
				}

				if gotUnits[0].ID == units[0].ID {
					matchFn(t, &gotUnits[0], units[0])
					matchFn(t, &gotUnits[1], units[1])
				} else {
					matchFn(t, &gotUnits[1], units[0])
					matchFn(t, &gotUnits[0], units[1])
				}
			},
		},
		{
			name: "BadRequest/InvalidFields",
			query: fiber.Map{
				"unitIds": fmt.Sprintf("unitIds=%s&unitIds=%s", units[0].ID.String(), units[1].ID.String()),
				"fields":  "fields=invalid_field",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					GetUnitsByIds(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "OK/SomeUnitsAreNotPublic",
			query: fiber.Map{
				"unitIds": fmt.Sprintf("unitIds=%s&unitIds=%s", units[0].ID.String(), units[1].ID.String()),
				"fields":  fmt.Sprintf("fields=%s", strings.Join(unit_dto.GetRetrievableFields(), ",")),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[0].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[1].ID)).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Eq(units[1].ID), gomock.Eq(userId)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					GetUnitsByIds(gomock.Any(), gomock.Eq([]string{units[0].ID.String(), units[1].ID.String()}), gomock.Eq(unit_dto.GetRetrievableFields())).
					Times(1).
					Return([]unit_model.UnitModel{*units[0], *units[1]}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)
			},
		},
		{
			name: "OK/SomeUnitsAreNotVisible",
			query: fiber.Map{
				"unitIds": fmt.Sprintf("unitIds=%s&unitIds=%s", units[0].ID.String(), units[1].ID.String()),
				"fields":  fmt.Sprintf("fields=%s", strings.Join(unit_dto.GetRetrievableFields(), ",")),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					uuid.MustParse("98ab8b6d-0427-4534-8c77-46ea8b02c8d9"), time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[0].ID)).
					Times(1).
					Return(true, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Eq(units[0].ID), gomock.Eq(uuid.MustParse("98ab8b6d-0427-4534-8c77-46ea8b02c8d9"))).
					Times(0)
				uRepo.EXPECT().
					IsPublic(gomock.Any(), gomock.Eq(units[1].ID)).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					CheckUnitManageability(gomock.Any(), gomock.Eq(units[1].ID), gomock.Eq(uuid.MustParse("98ab8b6d-0427-4534-8c77-46ea8b02c8d9"))).
					Times(1).
					Return(false, nil)
				uRepo.EXPECT().
					GetUnitsByIds(gomock.Any(), gomock.Eq([]string{units[0].ID.String()}), gomock.Eq(unit_dto.GetRetrievableFields())).
					Times(1).
					Return([]unit_model.UnitModel{*units[0]}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.NotEmpty(t, res.Body)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var gotUnits []unit_model.UnitModel
				err = json.Unmarshal(data, &gotUnits)
				require.NoError(t, err)
				require.Len(t, gotUnits, 1)
				require.Equal(t, units[0].ID, gotUnits[0].ID)
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

			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo)

			var queries []string
			for _, v := range tc.query {
				queries = append(queries, v.(string))
			}

			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/units/ids/?%s", strings.Join(queries, "&")),
				nil)
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}
