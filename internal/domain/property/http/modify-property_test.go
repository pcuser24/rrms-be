package http

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"go.uber.org/mock/gomock"
)

func TestUpdateProperty(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")

	property := property_repo.NewRandomPropertyModel(t, userId)

	testcases := []struct {
		name          string
		id            string
		body          fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			id:   property.ID.String(),
			body: fiber.Map{
				"name":    "abc",
				"project": "xyz",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
				pRepo.EXPECT().
					UpdateProperty(gomock.Any(), gomock.Eq(&dto.UpdateProperty{
						ID:      property.ID,
						Name:    types.Ptr[string]("abc"),
						Project: types.Ptr[string]("xyz"),
					})).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name:      "Unauthorized",
			id:        property.ID.String(),
			body:      fiber.Map{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					UpdateProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "Forbidden/PrivateProperty",
			id:   property.ID.String(),
			body: fiber.Map{
				"name": "abc",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					uuid.Nil, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
				pRepo.EXPECT().
					UpdateProperty(gomock.Any(), gomock.Eq(property.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
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

			req := httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/api/properties/property/%s", tc.id),
				bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestDeleteProperty(t *testing.T) {
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
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
				pRepo.EXPECT().
					DeleteProperty(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name:      "Unauthorized",
			id:        property.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				pRepo.EXPECT().
					DeleteProperty(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
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
					GetPropertyManagers(gomock.Any(), gomock.Eq(property.ID)).
					Times(1).
					Return(property.Managers, nil)
				pRepo.EXPECT().
					DeleteProperty(gomock.Any(), gomock.Eq(property.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
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

			req := httptest.NewRequest(
				http.MethodDelete,
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
