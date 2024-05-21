package http

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	listing_model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/utils/mock"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
	"go.uber.org/mock/gomock"
)

func TestCreateListing(t *testing.T) {
	userId := uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f")
	propertyId := uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a")
	unitIds := []uuid.UUID{
		uuid.MustParse("937b234a-d091-49a7-99dc-daba3dc68a9b"),
		uuid.MustParse("63037e39-5e0c-4e67-9187-d9236e1bfd28"),
		uuid.MustParse("86514de5-e4aa-4c17-9cce-02bf4c2d6a2c"),
	}
	listing := listing_repo.NewRandomListingModel(t, userId, propertyId, unitIds)
	defaultDto := listing_dto.CreateListing{
		CreatorID:         userId,
		PropertyID:        propertyId,
		Title:             listing.Title,
		Description:       listing.Description,
		FullName:          listing.FullName,
		Email:             listing.Email,
		Phone:             listing.Phone,
		ContactType:       listing.ContactType,
		Price:             listing.Price,
		PriceNegotiable:   listing.PriceNegotiable,
		SecurityDeposit:   listing.SecurityDeposit,
		LeaseTerm:         listing.LeaseTerm,
		PetsAllowed:       listing.PetsAllowed,
		NumberOfResidents: listing.NumberOfResidents,
		Priority:          listing.Priority,
		PostDuration:      15,
		Policies:          make([]listing_dto.CreateListingPolicy, len(listing.Policies)),
		Units:             make([]listing_dto.CreateListingUnit, len(listing.Units)),
	}
	for i, p := range listing.Policies {
		defaultDto.Policies[i] = listing_dto.CreateListingPolicy{
			PolicyID: p.PolicyID,
			Note:     p.Note,
		}
	}
	for i, u := range listing.Units {
		defaultDto.Units[i] = listing_dto.CreateListingUnit{
			UnitID: u.UnitID,
			Price:  u.Price,
		}
	}
	defaultArg := fiber.Map{
		"creatorId":         userId,
		"propertyId":        propertyId,
		"title":             listing.Title,
		"description":       listing.Description,
		"fullName":          listing.FullName,
		"email":             listing.Email,
		"phone":             listing.Phone,
		"contactType":       listing.ContactType,
		"price":             listing.Price,
		"priceNegotiable":   listing.PriceNegotiable,
		"securityDeposit":   listing.SecurityDeposit,
		"leaseTerm":         listing.LeaseTerm,
		"petsAllowed":       listing.PetsAllowed,
		"numberOfResidents": listing.NumberOfResidents,
		"priority":          listing.Priority,
		"postDuration":      15,
		"policies":          listing.Policies,
		"units":             listing.Units,
	}

	testcases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			body: defaultArg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{PropertyID: propertyId, ManagerID: userId, Role: "OWNER"},
					}, nil)
				_unitIds := make([]any, len(unitIds))
				for i, id := range unitIds {
					_unitIds[i] = id
				}
				uRepo.EXPECT().
					CheckUnitOfProperty(gomock.Any(), gomock.Eq(propertyId), mock.InMatcher(_unitIds)).
					Times(3).
					Return(true, nil)
				lRepo.EXPECT().
					CreateListing(gomock.Any(), gomock.Eq(&defaultDto)).
					Times(1).
					Return(listing, nil)
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

				requireBodyMatchListing(t, res.Body, listing)
			},
		},
		{
			name: "OK/MissingOptionalFields",
			body: (func() fiber.Map {
				// copy defaultArg
				arg := maps.Clone(defaultArg)
				delete(arg, "numberOfResidents")
				delete(arg, "policies")
				return arg
			})(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{PropertyID: propertyId, ManagerID: userId, Role: "OWNER"},
					}, nil)
				_unitIds := make([]any, len(unitIds))
				for i, id := range unitIds {
					_unitIds[i] = id
				}
				uRepo.EXPECT().
					CheckUnitOfProperty(gomock.Any(), gomock.Eq(propertyId), mock.InMatcher(_unitIds)).
					Times(3).
					Return(true, nil)
				dto := defaultDto
				dto.NumberOfResidents = nil
				dto.Policies = nil
				lRepo.EXPECT().
					CreateListing(gomock.Any(), gomock.Eq(&dto)).
					Times(1).
					Return(listing, nil)
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

				requireBodyMatchListing(t, res.Body, listing)
			},
		},
		{
			name: "BadRequest/MissingRequiredFields",
			body: (func() fiber.Map {
				// copy defaultArg
				arg := maps.Clone(defaultArg)
				delete(arg, "propertyId")
				return arg
			})(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Any()).
					Times(0)
				uRepo.EXPECT().
					CheckUnitOfProperty(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				lRepo.EXPECT().
					CreateListing(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "Forbidden/PrivateProperty",
			body: defaultArg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{}, nil)
				uRepo.EXPECT().
					CheckUnitOfProperty(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				lRepo.EXPECT().
					CreateListing(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusForbidden, res.StatusCode)
			},
		},
		{
			name: "Forbidden/UnitNotOfProperty",
			body: defaultArg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				auth_http.AddAuthorization(t, request,
					tokenMaker, auth_http.AuthorizationTypeBearer,
					userId, time.Minute, token.CreateTokenOptions{TokenType: token.AccessToken})
			},
			buildStubs: func(pRepo *property_repo.MockRepo, uRepo *unit_repo.MockRepo, lRepo *listing_repo.MockRepo) {
				pRepo.EXPECT().
					GetPropertyManagers(gomock.Any(), gomock.Eq(propertyId)).
					Times(1).
					Return([]property_model.PropertyManagerModel{
						{PropertyID: propertyId, ManagerID: userId, Role: "OWNER"},
					}, nil)
				_unitIds := make([]any, len(unitIds))
				for i, id := range unitIds {
					_unitIds[i] = id
				}
				uRepo.EXPECT().
					CheckUnitOfProperty(gomock.Any(), gomock.Eq(propertyId), mock.InMatcher(_unitIds)).
					Times(1).
					Return(false, nil)
				lRepo.EXPECT().
					CreateListing(gomock.Any(), gomock.Any()).
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
			paymentRepo := payment_repo.NewMockRepo(ctrl)
			rRepo := rental_repo.NewMockRepo(ctrl)
			authRepo := auth_repo.NewMockRepo(ctrl)

			tc.buildStubs(pRepo, uRepo, lRepo)

			srv := newTestServer(t, pRepo, uRepo, lRepo, aRepo, paymentRepo, rRepo, authRepo)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/listings/", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func requireBodyMatchListing(t *testing.T, body io.Reader, listing *listing_model.ListingModel) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var got listing_model.ListingModel
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, listing, &got)
}
