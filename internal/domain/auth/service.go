package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type AuthService interface {
	RegisterUser(data *dto.RegisterUser) (*model.UserModel, error)
	Login(data *dto.LoginUser, sessionData *dto.CreateSessionDto) (*LoginUserRes, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserById(id uuid.UUID) (*model.UserModel, error)
	Logout(id uuid.UUID) error
}

type authService struct {
	repo            AuthRepo
	tokenMaker      token.Maker
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repo AuthRepo, tokenMaker token.Maker, accessTokenTTL, refreshToken time.Duration) AuthService {
	return &authService{
		repo:            repo,
		tokenMaker:      tokenMaker,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshToken,
	}
}

func (u *authService) RegisterUser(data *dto.RegisterUser) (*model.UserModel, error) {
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

	return user, nil
}

type LoginUserRes struct {
	User           *model.UserModel
	AccessToken    string
	AccessPayload  *token.Payload
	RefreshToken   string
	RefreshPayload *token.Payload
	SessionID      string
}

var ErrInvalidCredential = fmt.Errorf("invalid password")

func (u *authService) Login(data *dto.LoginUser, sessionData *dto.CreateSessionDto) (*LoginUserRes, error) {
	user, err := u.repo.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	err = utils.VerifyPassword(*user.Password, data.Password)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	// Signin new user
	// Check for session existence
	var currentSession *model.SessionModel = nil
	if sessionData.ID != uuid.Nil {
		currentSession, err = u.repo.GetSessionById(context.Background(), sessionData.ID)
		if err != nil {
			return nil, err
		}
	}

	// If session exists and is not blocked
	if currentSession != nil &&
		*currentSession.UserAgent == string(sessionData.UserAgent) &&
		*currentSession.ClientIp == sessionData.ClientIp &&
		!currentSession.IsBlocked {
		// Create a new access token
		accessToken, accessPayload, err := u.tokenMaker.CreateToken(user.ID, u.accessTokenTTL, token.CreateTokenOptions{
			TokenType: token.AccessToken,
			TokenID:   currentSession.ID,
		})
		if err != nil {
			return nil, err
		}
		// Return the session
		return &LoginUserRes{
			User:           user,
			SessionID:      currentSession.ID.String(),
			RefreshToken:   currentSession.SessionToken,
			RefreshPayload: &token.Payload{ExpiredAt: currentSession.Expires},
			AccessToken:    accessToken,
			AccessPayload:  accessPayload,
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

	return &LoginUserRes{
		User:           user,
		AccessToken:    accessToken,
		AccessPayload:  accessPayload,
		RefreshToken:   refreshToken,
		RefreshPayload: refreshPayload,
		SessionID:      session.ID.String(),
	}, nil
}

func (u *authService) GetUserByEmail(email string) (*model.UserModel, error) {
	return u.repo.GetUserByEmail(context.Background(), email)
}

func (u *authService) GetUserById(id uuid.UUID) (*model.UserModel, error) {
	return u.repo.GetUserById(context.Background(), id)
}

func (u *authService) Logout(id uuid.UUID) error {
	return u.repo.UpdateSessionStatus(context.Background(), id, true)
}
