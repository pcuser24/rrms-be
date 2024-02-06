package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/auth/asynctask"

	repo2 "github.com/user2410/rrms-backend/internal/domain/auth/repo"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Service interface {
	Register(data *dto.RegisterUser) (*model.UserModel, error)
	Login(data *dto.LoginUser, sessionData *dto.CreateSession) (*dto.LoginUserRes, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserById(id uuid.UUID) (*model.UserModel, error)
	RefreshAccessToken(accessToken, refreshToken string) (*dto.LoginUserRes, error)
	Logout(id uuid.UUID) error
}

type service struct {
	repo            repo2.Repo
	tokenMaker      token.Maker
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	taskDistributor asynctask.TaskDistributor
}

func NewService(
	repo repo2.Repo,
	tokenMaker token.Maker,
	accessTokenTTL, refreshToken time.Duration,
	taskDistributor asynctask.TaskDistributor,
) Service {
	return &service{
		repo:            repo,
		tokenMaker:      tokenMaker,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshToken,
		taskDistributor: taskDistributor,
	}
}

func (u *service) Register(data *dto.RegisterUser) (*model.UserModel, error) {
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = hash

	// Create a new entry in User table
	user, err := u.repo.InsertUser(context.Background(), data)
	if err != nil {
		return nil, err
	}

	user.Password = nil
	return user, nil
}

var ErrInvalidCredential = fmt.Errorf("invalid password")

func (u *service) Login(data *dto.LoginUser, sessionData *dto.CreateSession) (*dto.LoginUserRes, error) {
	user, err := u.repo.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		return nil, err
	}

	err = utils.VerifyPassword(*user.Password, data.Password)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	// Check for session existence
	var currentSession *model.SessionModel = nil
	if sessionData.ID != uuid.Nil {
		currentSession, err = u.repo.GetSessionById(context.Background(), sessionData.ID)
		if err != nil && err != database.ErrRecordNotFound {
			return nil, err
		}
	}

	// If session exists and is not blocked
	if currentSession != nil &&
		*currentSession.UserAgent == string(sessionData.UserAgent) &&
		*currentSession.ClientIp == sessionData.ClientIp &&
		!currentSession.IsBlocked &&
		currentSession.Expires.After(time.Now()) {
		// Create a new access token
		accessToken, accessPayload, err := u.tokenMaker.CreateToken(user.ID, u.accessTokenTTL, token.CreateTokenOptions{
			TokenType: token.AccessToken,
			TokenID:   currentSession.ID,
		})
		if err != nil {
			return nil, err
		}
		// Return the session
		return &dto.LoginUserRes{
			User:         *user.ToUserResponse(),
			SessionID:    currentSession.ID,
			RefreshToken: currentSession.SessionToken,
			RefreshExp:   currentSession.Expires,
			AccessToken:  accessToken,
			AccessExp:    accessPayload.ExpiredAt,
		}, nil
	}

	// Otherwise create a new session
	// 1. Create refresh token
	refreshToken, refreshPayload, err := u.tokenMaker.CreateToken(user.ID, u.refreshTokenTTL, token.CreateTokenOptions{
		TokenType: token.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	// 2. Create access token
	accessToken, accessPayload, err := u.tokenMaker.CreateToken(user.ID, u.accessTokenTTL, token.CreateTokenOptions{
		TokenType: token.AccessToken,
		TokenID:   refreshPayload.ID,
	})
	if err != nil {
		return nil, err
	}
	sessionData.ID = refreshPayload.ID
	sessionData.UserId = user.ID
	sessionData.SessionToken = refreshToken
	sessionData.Expires = refreshPayload.ExpiredAt
	sessionData.CreatedAt = refreshPayload.IssuedAt
	// 3. Create a new session
	session, err := u.repo.CreateSession(context.Background(), sessionData)
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserRes{
		User:         *user.ToUserResponse(),
		AccessToken:  accessToken,
		AccessExp:    accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshExp:   refreshPayload.ExpiredAt,
		SessionID:    session.ID,
	}, nil
}

func (u *service) GetUserByEmail(email string) (*model.UserModel, error) {
	return u.repo.GetUserByEmail(context.Background(), email)
}

func (u *service) GetUserById(id uuid.UUID) (*model.UserModel, error) {
	return u.repo.GetUserById(context.Background(), id)
}

var ErrInvalidSession = fmt.Errorf("invalid session")

func (u *service) RefreshAccessToken(accessToken, refreshToken string) (*dto.LoginUserRes, error) {
	accessPayload, err := u.tokenMaker.VerifyToken(accessToken)
	if accessPayload == nil { // Invalid access token
		return nil, token.ErrInvalidToken
	}

	// If access token is not expired, return it
	if time.Now().Before(accessPayload.ExpiredAt) {
		return &dto.LoginUserRes{
			AccessToken: accessToken,
			AccessExp:   accessPayload.ExpiredAt,
		}, nil
	}

	session, err := u.repo.GetSessionById(context.Background(), accessPayload.ID)
	if err != nil {
		return nil, err
	}

	if session.SessionToken != refreshToken {
		return nil, ErrInvalidSession
	}

	if session.IsBlocked {
		return nil, ErrInvalidSession
	}

	if session.UserId != accessPayload.UserID {
		return nil, ErrInvalidSession
	}

	if time.Now().After(session.Expires) {
		return nil, ErrInvalidSession
	}

	newAccessToken, newAccessPayload, err := u.tokenMaker.CreateToken(
		accessPayload.UserID,
		u.accessTokenTTL,
		token.CreateTokenOptions{
			TokenType: token.AccessToken,
			TokenID:   accessPayload.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserRes{
		AccessToken: newAccessToken,
		AccessExp:   newAccessPayload.ExpiredAt,
	}, nil
}

func (u *service) Logout(id uuid.UUID) error {
	return u.repo.UpdateSessionStatus(context.Background(), id, true)
}
