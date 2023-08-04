package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	db "github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type AuthRepo interface {
	InsertUser(ctx context.Context, data *dto.RegisterUser) (*model.UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserModel, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.UserModel, error)
}

type authRepo struct {
	dao db.DAO
}

func NewUserRepo(d db.DAO) AuthRepo {
	return &authRepo{
		dao: d,
	}
}

func (u *authRepo) InsertUser(ctx context.Context, data *dto.RegisterUser) (*model.UserModel, error) {
	res, err := u.dao.InsertUser(ctx, db.InsertUserParams{
		Email:    data.Email,
		Password: types.StrN(&data.Password),
	})
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}

func (u *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	res, err := u.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}

func (u *authRepo) GetUserById(ctx context.Context, id uuid.UUID) (*model.UserModel, error) {
	res, err := u.dao.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}
