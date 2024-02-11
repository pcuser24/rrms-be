package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_dto "github.com/user2410/rrms-backend/internal/domain/unit/dto"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
	"go.uber.org/mock/gomock"

	"github.com/google/uuid"
)

func TestCreateUnit(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")
	propertyId := uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a")
	unit := unit_repo.NewRandomUnitModel(t, propertyId)
	defaultArg := unit_dto.CreateUnit{
		PropertyID:          propertyId,
		Name:                &unit.Name,
		Area:                unit.Area,
		Floor:               unit.Floor,
		NumberOfLivingRooms: unit.NumberOfLivingRooms,
		NumberOfBedrooms:    unit.NumberOfBedrooms,
		NumberOfBathrooms:   unit.NumberOfBathrooms,
		NumberOfToilets:     unit.NumberOfToilets,
		NumberOfKitchens:    unit.NumberOfKitchens,
		NumberOfBalconies:   unit.NumberOfBalconies,
		Type:                unit.Type,
	}
	defaultArg.Amenities = make([]unit_dto.CreateUnitAmenity, len(unit.Amenities))
	for i, a := range unit.Amenities {
		defaultArg.Amenities[i] = unit_dto.CreateUnitAmenity{
			AmenityID:   a.AmenityID,
			Description: a.Description,
		}
	}
	defaultArg.Media = make([]unit_dto.CreateUnitMedia, len(unit.Media))
	for i, m := range unit.Media {
		defaultArg.Media[i] = unit_dto.CreateUnitMedia{
			Url:         m.Url,
			Type:        m.Type,
			Description: m.Description,
		}
	}

	testcases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"propertyId":          defaultArg.PropertyID,
				"name":                *defaultArg.Name,
				"area":                defaultArg.Area,
				"floor":               *defaultArg.Floor,
				"numberOfLivingRooms": *defaultArg.NumberOfLivingRooms,
				"numberOfBedrooms":    *defaultArg.NumberOfBedrooms,
				"numberOfBathrooms":   *defaultArg.NumberOfBathrooms,
				"numberOfToilets":     *defaultArg.NumberOfToilets,
				"numberOfKitchens":    *defaultArg.NumberOfKitchens,
				"numberOfBalconies":   *defaultArg.NumberOfBalconies,
				"type":                defaultArg.Type,
				"amenities":           defaultArg.Amenities,
				"media":               defaultArg.Media,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{ManagerID: userId, PropertyID: propertyId, Role: "OWNER"},
					}, nil)
				uRepo.EXPECT().
					CreateUnit(gomock.Any(), gomock.Eq(&defaultArg)).
					Times(1).
					Return(unit, nil)
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
				require.NotEmpty(t, res.Body)

				requireBodyMatchUnit(t, res.Body, unit)
			},
		},
		{
			name: "OK/MissingOptionalFields",
			body: fiber.Map{
				"propertyId":        defaultArg.PropertyID,
				"name":              defaultArg.Name,
				"area":              defaultArg.Area,
				"floor":             defaultArg.Floor,
				"numberOfBedrooms":  *defaultArg.NumberOfBedrooms,
				"numberOfBathrooms": *defaultArg.NumberOfBathrooms,
				"numberOfBalconies": *defaultArg.NumberOfBalconies,
				"type":              defaultArg.Type,
				"amenities":         defaultArg.Amenities,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				arg := defaultArg
				arg.NumberOfLivingRooms = nil
				arg.NumberOfToilets = nil
				arg.NumberOfKitchens = nil
				arg.Media = nil

				_unit := *unit
				_unit.NumberOfLivingRooms = nil
				_unit.NumberOfToilets = nil
				_unit.NumberOfKitchens = nil
				_unit.Media = nil

				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{ManagerID: userId, PropertyID: propertyId, Role: "OWNER"},
					}, nil)

				uRepo.EXPECT().
					CreateUnit(gomock.Any(), gomock.Eq(&arg)).
					Times(1).
					Return(&_unit, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusCreated, res.StatusCode)
				require.NotEmpty(t, res.Body)

				_unit := *unit
				_unit.NumberOfLivingRooms = nil
				_unit.NumberOfToilets = nil
				_unit.NumberOfKitchens = nil
				_unit.Media = nil
				requireBodyMatchUnit(t, res.Body, &_unit)
			},
		},
		{
			name: "BadRequest/MissingRequiredFields",
			body: fiber.Map{
				"propertyId":        defaultArg.PropertyID,
				"name":              defaultArg.Name,
				"floor":             defaultArg.Floor,
				"numberOfBedrooms":  *defaultArg.NumberOfBedrooms,
				"numberOfBathrooms": *defaultArg.NumberOfBathrooms,
				"numberOfBalconies": *defaultArg.NumberOfBalconies,
				"type":              defaultArg.Type,
				"amenities":         defaultArg.Amenities,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					CreateUnit(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "BadRequest/InvalidFieldType",
			body: fiber.Map{
				"propertyId":          defaultArg.PropertyID,
				"name":                defaultArg.Name,
				"area":                defaultArg.Name,
				"floor":               defaultArg.Floor,
				"numberOfLivingRooms": *defaultArg.NumberOfLivingRooms,
				"numberOfBedrooms":    *defaultArg.NumberOfBedrooms,
				"numberOfBathrooms":   *defaultArg.NumberOfBathrooms,
				"numberOfToilets":     defaultArg.Name,
				"numberOfKitchens":    *defaultArg.NumberOfKitchens,
				"numberOfBalconies":   *defaultArg.NumberOfBalconies,
				"type":                defaultArg.Type,
				"amenities":           defaultArg.Amenities,
				"media":               defaultArg.Media,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					CreateUnit(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "Forbidden/NotManagerOfTheProperty",
			body: fiber.Map{
				"propertyId":          defaultArg.PropertyID,
				"name":                *defaultArg.Name,
				"area":                defaultArg.Area,
				"floor":               *defaultArg.Floor,
				"numberOfLivingRooms": *defaultArg.NumberOfLivingRooms,
				"numberOfBedrooms":    *defaultArg.NumberOfBedrooms,
				"numberOfBathrooms":   *defaultArg.NumberOfBathrooms,
				"numberOfToilets":     *defaultArg.NumberOfToilets,
				"numberOfKitchens":    *defaultArg.NumberOfKitchens,
				"numberOfBalconies":   *defaultArg.NumberOfBalconies,
				"type":                defaultArg.Type,
				"amenities":           defaultArg.Amenities,
				"media":               defaultArg.Media,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{ManagerID: uuid.MustParse("98ab8b6d-0427-4534-8c77-46ea8b02c8d9"), PropertyID: propertyId, Role: "OWNER"},
						{ManagerID: uuid.MustParse("f6ca05c0-fad5-46fc-a237-a8e930e7cb13"), PropertyID: propertyId, Role: "MANAGER"},
					}, nil)
				uRepo.EXPECT().
					CreateUnit(gomock.Any(), gomock.Any()).
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
			tc.buildStubs(pRepo, uRepo)

			srv := newTestServer(t, pRepo, uRepo)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/units/", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func requireBodyMatchUnit(t *testing.T, body io.ReadCloser, unit *unit_model.UnitModel) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res unit_model.UnitModel
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	unit_repo.SameUnits(t, unit, &res)
}
