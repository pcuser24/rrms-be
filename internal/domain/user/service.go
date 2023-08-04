package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/user/dto"
	"github.com/user2410/rrms-backend/internal/domain/user/model"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type UserService interface {
	RegisterUser(data dto.RegisterUser) (*RegisterUserRes, error)
	Login(data *dto.LoginUser) (*LoginUserRes, error)
	GetUserByEmail(email string) (*model.UserModel, error)
	GetUserById(id uuid.UUID) (*model.UserModel, error)
}

type userService struct {
	repo           UserRepo
	tokenMaker     token.Maker
	accessTokenTTL time.Duration
}

func NewUserService(repo UserRepo, tokenMaker token.Maker, accessTokenTTL time.Duration) UserService {
	return &userService{
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

func (u *userService) RegisterUser(data dto.RegisterUser) (*RegisterUserRes, error) {
	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}
	user := &model.UserModel{
		Email:    data.Email,
		Password: types.Ptr[string](hash),
	}

	user, err = u.repo.InsertUser(user)
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

func (u *userService) Login(data *dto.LoginUser) (*LoginUserRes, error) {
	user, err := u.repo.GetUserByEmail(data.Email)
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

func (u *userService) InsertUser(user *model.UserModel) (*model.UserModel, error) {
	return u.repo.InsertUser(user)
}

func (u *userService) GetUserByEmail(email string) (*model.UserModel, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *userService) GetUserById(id uuid.UUID) (*model.UserModel, error) {
	return u.repo.GetUserById(id)
}
