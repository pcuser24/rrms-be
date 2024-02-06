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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/auth/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/domain/auth/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	randomEmail := random.RandomEmail()
	randomPassword := random.RandomAlphanumericStr(10)
	hashedPassword, err := utils.HashPassword(randomPassword)
	require.NoError(t, err)

	testcases := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(repo *repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"email":    randomEmail,
				"password": randomPassword,
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&model.UserModel{
						ID:        uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f"),
						Email:     randomEmail,
						Password:  &hashedPassword,
						CreatedAt: time.Now(),
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusCreated, res.StatusCode)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				var user model.UserModel
				err = json.Unmarshal(data, &user)
				require.NoError(t, err)

				require.NotZero(t, user.ID)
				require.Equal(t, randomEmail, user.Email)
				require.Nil(t, user.Password)
				require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)
			},
		},
		{
			name: "BadRequest/MissingField",
			body: fiber.Map{
				"email": randomEmail,
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "BadRequest/InvalidEmail",
			body: fiber.Map{
				"email":    "invalid-email",
				"password": randomPassword,
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
		{
			name: "BadRequest/InvalidPassword",
			body: fiber.Map{
				"email":    randomEmail,
				"password": "short",
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					InsertUser(gomock.Any(), gomock.Any()).
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

			repo := repo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			a := asynctask.NewMockTaskDistributor(ctrl)

			srv := newTestServer(t, repo, a)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/auth/credential/register", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestLogin(t *testing.T) {
	randomPassword := random.RandomAlphanumericStr(10)
	hashedPassword, err := utils.HashPassword(randomPassword)
	require.NoError(t, err)

	user := &model.UserModel{
		ID:        uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f"),
		Email:     random.RandomEmail(),
		Password:  &hashedPassword,
		GroupID:   uuid.Nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testcases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(repo *repo.MockRepo)
		checkResponse func(t *testing.T, res *http.Response)
	}{
		{
			name: "OK/NoCurrentSession",
			body: fiber.Map{
				"email":    user.Email,
				"password": randomPassword,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				repo.EXPECT().
					CreateSession(
						gomock.Any(),
						gomock.Any(),
						// gomock.Eq(&dto.CreateSession{
						// 	ID:           uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), // Unknown
						// 	UserId:       user.ID,
						// 	SessionToken: "seSsIoN-ToKeN", // Unknown
						// 	Expires:      time.Now().Add(time.Hour),
						// 	UserAgent: []byte(""),
						// 	ClientIp: "0.0.0.0",
						// }),
					).
					Times(1).
					Return(&model.SessionModel{
						ID:           uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"),
						SessionToken: "seSsIoN-ToKeN",
						UserId:       user.ID,
						CreatedAt:    time.Now(),
						Expires:      time.Now().Add(15 * time.Minute),
						UserAgent:    types.Ptr[string](""),
						ClientIp:     types.Ptr[string]("0.0.0.0"),
						IsBlocked:    false,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				var loginRes dto.LoginUserRes
				err = json.Unmarshal(data, &loginRes)
				require.NoError(t, err)

				require.Equal(t, uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), loginRes.SessionID)
				require.Equal(t, user.Email, loginRes.User.Email)
				require.NotEmpty(t, loginRes.AccessToken)
				require.WithinDuration(t, loginRes.AccessExp, time.Now().Add(time.Minute), time.Second)
				require.NotEmpty(t, loginRes.RefreshToken)
				require.WithinDuration(t, loginRes.RefreshExp, time.Now().Add(time.Hour), time.Second)
			},
		},
		{
			name: "OK/WithActiveSession",
			body: fiber.Map{
				"email":    user.Email,
				"password": randomPassword,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, req,
					tokenMaker, AuthorizationTypeBearer,
					user.ID, time.Minute, token.CreateTokenOptions{TokenID: uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), TokenType: token.AccessToken})
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				repo.EXPECT().
					GetSessionById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&model.SessionModel{
						ID:           uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"),
						SessionToken: "seSsIoN-ToKeN",
						UserId:       user.ID,
						UserAgent:    types.Ptr[string](""),
						ClientIp:     types.Ptr[string]("0.0.0.0"),
						CreatedAt:    time.Now(),
						IsBlocked:    false,
						Expires:      time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				var loginRes dto.LoginUserRes
				err = json.Unmarshal(data, &loginRes)
				require.NoError(t, err)

				require.Equal(t, uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), loginRes.SessionID)
				require.Equal(t, user.Email, loginRes.User.Email)
				require.NotEmpty(t, loginRes.AccessToken)
				require.WithinDuration(t, loginRes.AccessExp, time.Now().Add(time.Minute), time.Second)
				require.Equal(t, "seSsIoN-ToKeN", loginRes.RefreshToken)
				require.WithinDuration(t, loginRes.RefreshExp, time.Now().Add(time.Hour), time.Second)
			},
		},
		{
			name: "OK/WithBlockedSession",
			body: fiber.Map{
				"email":    user.Email,
				"password": randomPassword,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				AddAuthorization(t, req,
					tokenMaker, AuthorizationTypeBearer,
					user.ID, time.Minute, token.CreateTokenOptions{TokenID: uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), TokenType: token.AccessToken})
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
				repo.EXPECT().
					GetSessionById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&model.SessionModel{
						ID:           uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"),
						SessionToken: "seSsIoN-ToKeN",
						UserId:       user.ID,
						UserAgent:    types.Ptr[string](""),
						ClientIp:     types.Ptr[string]("0.0.0.0"),
						CreatedAt:    time.Now(),
						IsBlocked:    true,
						Expires:      time.Now().Add(time.Hour),
					}, nil)
				repo.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&model.SessionModel{
						ID:           uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"),
						SessionToken: "seSsIoN-ToKeN",
						UserId:       user.ID,
						CreatedAt:    time.Now(),
						Expires:      time.Now().Add(15 * time.Minute),
						UserAgent:    types.Ptr[string](""),
						ClientIp:     types.Ptr[string]("0.0.0.0"),
						IsBlocked:    false,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, res.StatusCode)

				data, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				var loginRes dto.LoginUserRes
				err = json.Unmarshal(data, &loginRes)
				require.NoError(t, err)

				require.Equal(t, uuid.MustParse("efefb62c-0b8e-4868-bfe7-470fcf9f409a"), loginRes.SessionID)
				require.Equal(t, user.Email, loginRes.User.Email)
				require.NotEmpty(t, loginRes.AccessToken)
				require.WithinDuration(t, loginRes.AccessExp, time.Now().Add(time.Minute), time.Second)
				require.NotEmpty(t, loginRes.RefreshToken)
				require.WithinDuration(t, loginRes.RefreshExp, time.Now().Add(time.Hour), time.Second)
			},
		},
		{
			name: "Notfound/NoUserWithSuchEmail",
			body: fiber.Map{
				"email":    user.Email,
				"password": randomPassword,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(nil, database.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusNotFound, res.StatusCode)
			},
		},
		{
			name: "Unauthorized/WrongPassword",
			body: fiber.Map{
				"email":    user.Email,
				"password": "wrong-password",
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
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

			repo := repo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			a := asynctask.NewMockTaskDistributor(ctrl)

			srv := newTestServer(t, repo, a)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/auth/credential/login", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, srv.tokenMaker)
			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	// randomPassword := random.RandomAlphanumericStr(10)
	// hashedPassword, err := utils.HashPassword(randomPassword)
	// require.NoError(t, err)

	// user := &model.UserModel{
	// 	ID:        uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f"),
	// 	Email:     random.RandomEmail(),
	// 	Password:  &hashedPassword,
	// 	GroupID:   uuid.Nil,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	userId := uuid.MustParse("a2712956-2cbd-4c75-a57f-9d1d33a7fdcc")

	// accessToken, accessPayload, err := tokenMaker.CreateToken(userId, time.Minute, token.CreateTokenOptions{
	// 	TokenType: token.AccessToken,
	// 	TokenID:   uuid.MustParse("12c7ab40-26f7-4c77-859d-b93707de7430"),
	// })

	testcases := []struct {
		name          string
		body          func(t *testing.T, tokenMaker token.Maker) fiber.Map
		buildStubs    func(repo *repo.MockRepo)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK/NonExpiredToken",
			body: func(t *testing.T, tokenMaker token.Maker) fiber.Map {
				accessToken, accessPayload, err := tokenMaker.CreateToken(userId, time.Minute, token.CreateTokenOptions{
					TokenType: token.AccessToken,
					TokenID:   uuid.MustParse("12c7ab40-26f7-4c77-859d-b93707de7430"),
				})
				require.NoError(t, err)
				return fiber.Map{
					"refresh_token":  "refresh",
					"access_token":   accessToken,
					"access_payload": accessPayload,
				}
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetSessionById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.NotEmpty(t, response)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotEmpty(t, response.Body)

				data, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				var loginRes dto.LoginUserRes
				err = json.Unmarshal(data, &loginRes)
				require.NoError(t, err)
				require.NotEmpty(t, loginRes.AccessToken)
				require.WithinDuration(t, time.Now().Add(time.Minute), loginRes.AccessExp, time.Second)
			},
		},
		{
			name: "OK/ExpiredToken",
			body: func(t *testing.T, tokenMaker token.Maker) fiber.Map {
				return fiber.Map{
					"refresh_token": "refresh",
					"access_token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyYzdhYjQwLTI2ZjctNGM3Ny04NTlkLWI5MzcwN2RlNzQzMCIsInR5cGUiOiJhY2Nlc3MiLCJzdWIiOiJhMjcxMjk1Ni0yY2JkLTRjNzUtYTU3Zi05ZDFkMzNhN2ZkY2MiLCJpYXQiOiIyMDIzLTEyLTI4VDEwOjU4OjE5LjA0NjQ2MDIzNSswNzowMCIsImV4cCI6IjIwMjMtMTItMjhUMjI6NTg6MTkuMDQ2NDYwMzc2KzA3OjAwIn0.MpEQTBmvEgLbR5GhqaKlpK-cWKhLiGs-kaXY_erIVzY",
				}
			},
			buildStubs: func(repo *repo.MockRepo) {
				repo.EXPECT().
					GetSessionById(gomock.Any(), gomock.Eq(uuid.MustParse("12c7ab40-26f7-4c77-859d-b93707de7430"))).
					Times(1).
					Return(&model.SessionModel{
						ID:           uuid.MustParse("12c7ab40-26f7-4c77-859d-b93707de7430"),
						UserId:       userId,
						SessionToken: "refresh",
						IsBlocked:    false,
						Expires:      time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.NotEmpty(t, response)
				require.Equal(t, http.StatusOK, response.StatusCode)

				data, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				var loginRes dto.LoginUserRes
				err = json.Unmarshal(data, &loginRes)
				require.NoError(t, err)
				require.NotEmpty(t, loginRes.AccessToken)
				require.WithinDuration(t, time.Now(), loginRes.AccessExp, time.Hour)
			},
		},
	}

	for i := range testcases {
		tc := &testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repo.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			a := asynctask.NewMockTaskDistributor(ctrl)

			srv := newTestServer(t, repo, a)

			data, err := json.Marshal(tc.body(t, srv.tokenMaker))
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/api/auth/credential/refresh", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")

			res, err := srv.router.GetFibApp().Test(req)
			assert.NoError(t, err)

			tc.checkResponse(t, res)
		})
	}
}
