package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
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
	GetUsersByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]model.UserModel, error)
	GetSessionById(ctx context.Context, id uuid.UUID) (*model.SessionModel, error)
	UpdateUser(ctx context.Context, id uuid.UUID, data *dto.UpdateUser) error
	UpdateSessionStatus(ctx context.Context, id uuid.UUID, isBlocked bool) error
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (u *repo) CreateUser(ctx context.Context, data *dto.RegisterUser) (*model.UserModel, error) {
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

func (u *repo) CreateSession(ctx context.Context, data *dto.CreateSession) (*model.SessionModel, error) {
	res, err := u.dao.CreateSession(ctx, *data.ToCreateSessionParams())
	if err != nil {
		return nil, err
	}

	return model.ToSessionModel(&res), nil
}

func (u *repo) GetUserByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	res, err := u.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}

func (u *repo) GetUserById(ctx context.Context, id uuid.UUID) (*model.UserModel, error) {
	res, err := u.dao.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.ToUserModel(&res), nil
}

func (u *repo) GetSessionById(ctx context.Context, id uuid.UUID) (*model.SessionModel, error) {
	res, err := u.dao.GetSessionById(ctx, id)
	if err != nil {
		return nil, nil
	}
	return model.ToSessionModel(&res), err
}

func (u *repo) UpdateSessionStatus(ctx context.Context, id uuid.UUID, isBlocked bool) error {
	return u.dao.UpdateSessionBlockingStatus(ctx, database.UpdateSessionBlockingStatusParams{
		ID:        id,
		IsBlocked: isBlocked,
	})
}

func (u *repo) UpdateUser(ctx context.Context, id uuid.UUID, data *dto.UpdateUser) error {
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

func (r *repo) GetUsersByIds(ctx context.Context, ids []uuid.UUID, fields []string) ([]model.UserModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var nonFKFields []string = []string{"id"}
	nonFKFields = append(nonFKFields, fields...)
	// log.Println(nonFKFields, fkFields)

	// get non fk fields
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("\"User\"")
	ib.Where(ib.In("id::text", sqlbuilder.List(func() []string {
		var res []string
		for _, id := range ids {
			res = append(res, id.String())
		}
		return res
	}())))
	query, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.UserModel
	var i database.User
	var scanningFields []interface{} = []interface{}{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "first_name":
			scanningFields = append(scanningFields, &i.FirstName)
		case "last_name":
			scanningFields = append(scanningFields, &i.LastName)
		case "email":
			scanningFields = append(scanningFields, &i.Email)
		case "phone":
			scanningFields = append(scanningFields, &i.Phone)
		case "avatar":
			scanningFields = append(scanningFields, &i.Avatar)
		case "address":
			scanningFields = append(scanningFields, &i.Address)
		case "city":
			scanningFields = append(scanningFields, &i.City)
		case "district":
			scanningFields = append(scanningFields, &i.District)
		case "ward":
			scanningFields = append(scanningFields, &i.Ward)
		case "role":
			scanningFields = append(scanningFields, &i.Role)
		case "created_at":
			scanningFields = append(scanningFields, &i.CreatedAt)
		case "updated_at":
			scanningFields = append(scanningFields, &i.UpdatedAt)
		case "created_by":
			scanningFields = append(scanningFields, &i.CreatedBy)
		case "updated_by":
			scanningFields = append(scanningFields, &i.UpdatedBy)
		case "deleted_f":
			scanningFields = append(scanningFields, &i.DeletedF)
		}
	}
	for rows.Next() {
		if err := rows.Scan(scanningFields...); err != nil {
			return nil, err
		}
		items = append(items, *model.ToUserModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
