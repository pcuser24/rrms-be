package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type AuthService interface {
	RegisterUser(data *dto.RegisterUser) (*RegisterUserRes, error)
	Login(data *dto.LoginUser) (*LoginUserRes, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetUserById(id uuid.UUID) (*dto.UserResponse, error)
}

type authService struct {
	repo           AuthRepo
	tokenMaker     token.Maker
	accessTokenTTL time.Duration
}

func NewUserService(repo AuthRepo, tokenMaker token.Maker, accessTokenTTL time.Duration) AuthService {
	return &authService{
		repo:           repo,
		tokenMaker:     tokenMaker,
		accessTokenTTL: accessTokenTTL,
	}
}

type RegisterUserRes struct {
	User          *dto.UserResponse
	AccessToken   string
	AccessPayload *token.Payload
	// RefreshToken string
	// RefreshPayload *token.Payload
}

func (u *authService) RegisterUser(data *dto.RegisterUser) (*RegisterUserRes, error) {
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	data.Password = hash

	user, err := u.repo.InsertUser(context.Background(), data)
	if err != nil {
		return nil, err
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(user.ID, token.AccessToken, u.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &RegisterUserRes{
		User:          user.ToUserResponse(),
		AccessToken:   accessToken,
		AccessPayload: accessPayload,
	}, nil
}

type LoginUserRes struct {
	User          *dto.UserResponse
	AccessToken   string
	AccessPayload *token.Payload
}

func (u *authService) Login(data *dto.LoginUser) (*LoginUserRes, error) {
	user, err := u.repo.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		return nil, err
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(user.ID, token.AccessToken, u.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &LoginUserRes{
		User:          user.ToUserResponse(),
		AccessToken:   accessToken,
		AccessPayload: accessPayload,
	}, nil
}

func (u *authService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (u *authService) GetUserById(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := u.repo.GetUserById(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}
