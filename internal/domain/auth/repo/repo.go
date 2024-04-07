package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Repo interface {
	CreateUser(ctx context.Context, data *dto.RegisterUser) (*model.UserModel, error)
	CreateSession(ctx context.Context, data *dto.CreateSession) (*model.SessionModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserModel, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.UserModel, error)
	GetSessionById(ctx context.Context, id uuid.UUID) (*model.SessionModel, error)
	UpdateUser(ctx context.Context, id uuid.UUID, data *dto.UpdateUser) error
	UpdateSessionStatus(ctx context.Context, id uuid.UUID, isBlocked bool) error
}

type authRepo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &authRepo{
		dao: d,
	}
}

func (u *authRepo) CreateUser(ctx context.Context, data *dto.RegisterUser) (*model.UserModel, error) {
	res, err := u.dao.CreateUser(ctx, database.CreateUserParams{
		Email:     data.Email,
		Password:  types.StrN(&data.Password),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Role:      data.Role,
	})
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}

func (u *authRepo) CreateSession(ctx context.Context, data *dto.CreateSession) (*model.SessionModel, error) {
	res, err := u.dao.CreateSession(ctx, *data.ToCreateSessionParams())
	if err != nil {
		return nil, err
	}

	return model.ToSessionModel(&res), nil
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

func (u *authRepo) GetSessionById(ctx context.Context, id uuid.UUID) (*model.SessionModel, error) {
	res, err := u.dao.GetSessionById(ctx, id)
	if err != nil {
		return nil, nil
	}
	return model.ToSessionModel(&res), err
}

func (u *authRepo) UpdateSessionStatus(ctx context.Context, id uuid.UUID, isBlocked bool) error {
	return u.dao.UpdateSessionBlockingStatus(ctx, database.UpdateSessionBlockingStatusParams{
		ID:        id,
		IsBlocked: isBlocked,
	})
}

func (u *authRepo) UpdateUser(ctx context.Context, id uuid.UUID, data *dto.UpdateUser) error {
	return u.dao.UpdateUser(ctx, database.UpdateUserParams{
		ID:        id,
		UpdatedBy: types.UUIDN(data.UpdatedBy),
		Email:     types.StrN(data.Email),
		Password:  types.StrN(data.Password),
		FirstName: types.StrN(data.FirstName),
		LastName:  types.StrN(data.LastName),
		Phone:     types.StrN(data.Phone),
		Avatar:    types.StrN(data.Avatar),
		Address:   types.StrN(data.Address),
		City:      types.StrN(data.City),
		District:  types.StrN(data.District),
		Ward:      types.StrN(data.Ward),
		Role: database.NullUSERROLE{
			USERROLE: data.Role,
			Valid:    data.Role != "",
		},
	})
}
