package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/user2410/rrms-backend/pkg/ds/set"

	repos "github.com/user2410/rrms-backend/internal/domain/_repos"

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
	GetUserById(id uuid.UUID) (*dto.UserResponse, error)
	GetUserByIds(ids []uuid.UUID) ([]dto.UserResponse, error)
	RefreshAccessToken(accessToken, refreshToken string) (*dto.LoginUserRes, error)
	Logout(id uuid.UUID) error
	UpdateUser(currentUserId, targetUserId uuid.UUID, data *dto.UpdateUser) error
}

type service struct {
	domainRepo      repos.DomainRepo
	tokenMaker      token.Maker
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(
	domainRepo repos.DomainRepo,
	tokenMaker token.Maker,
	accessTokenTTL, refreshToken time.Duration,
) Service {
	return &service{
		domainRepo:      domainRepo,
		tokenMaker:      tokenMaker,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshToken,
	}
}

func (u *service) Register(data *dto.RegisterUser) (*model.UserModel, error) {
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	data.Password = hash

	// Create a new entry in User table
	user, err := u.domainRepo.AuthRepo.CreateUser(context.Background(), data)
	if err != nil {
		return nil, err
	}

	user.Password = nil
	return user, nil
}

var ErrInvalidCredential = fmt.Errorf("invalid password")

func (u *service) Login(data *dto.LoginUser, sessionData *dto.CreateSession) (*dto.LoginUserRes, error) {
	user, err := u.domainRepo.AuthRepo.GetUserByEmail(context.Background(), data.Email)
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
		currentSession, err = u.domainRepo.AuthRepo.GetSessionById(context.Background(), sessionData.ID)
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
	session, err := u.domainRepo.AuthRepo.CreateSession(context.Background(), sessionData)
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
	return u.domainRepo.AuthRepo.GetUserByEmail(context.Background(), email)
}

func (u *service) GetUserById(id uuid.UUID) (*dto.UserResponse, error) {
	res, err := u.domainRepo.AuthRepo.GetUserById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return res.ToUserResponse(), nil
}

func (u *service) GetUserByIds(ids []uuid.UUID) ([]dto.UserResponse, error) {
	_ids := set.NewSet[uuid.UUID]()
	for _, id := range ids {
		if id != uuid.Nil {
			_ids.Add(id)
		}
	}
	users := make([]dto.UserResponse, 0, _ids.Size())
	for id := range _ids {
		user, err := u.GetUserById(id)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

var ErrInvalidSession = fmt.Errorf("invalid session")

func (u *service) RefreshAccessToken(accessToken, refreshToken string) (*dto.LoginUserRes, error) {
	accessPayload, _ := u.tokenMaker.VerifyToken(accessToken)
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

	session, err := u.domainRepo.AuthRepo.GetSessionById(context.Background(), accessPayload.ID)
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
	return u.domainRepo.AuthRepo.UpdateSessionStatus(context.Background(), id, true)
}

func (u *service) UpdateUser(currentUserId, targetUserId uuid.UUID, data *dto.UpdateUser) error {
	// Check if the user is updating his own information
	// TODO: Managed user info can be modified by managing user
	// if currentUserId != targetUserId {
	// 	return fmt.Errorf("unauthorized")
	// }

	// Preprocess data
	data.UpdatedBy = currentUserId

	if data.Password != nil {
		hash, err := utils.HashPassword(*data.Password)
		if err != nil {
			return err
		}
		data.Password = &hash
	}

	// Update user
	return u.domainRepo.AuthRepo.UpdateUser(context.Background(), targetUserId, data)
}
