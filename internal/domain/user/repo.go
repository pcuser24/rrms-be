package user

import (
	"database/sql"
	"github.com/user2410/rrms-backend/internal/domain/user/model"
	"github.com/user2410/rrms-backend/internal/domain/user/queries"
	"github.com/user2410/rrms-backend/internal/infrastructure/repositories/db"
)

type UserRepo interface {
	GetUserByEmail(email string) (*model.UserModel, error)
	Login(email string, password string) (*model.UserModel, error)
	InsertUser(user *model.UserModel) (*model.UserModel, error)
}

type userRepo struct {
	dao db.DAO
}

func NewUserRepo(d db.DAO) UserRepo {
	return &userRepo{
		dao: d,
	}
}

func (u *userRepo) GetUserByEmail(email string) (*model.UserModel, error) {
	rows, err := u.dao.NamedQuery(queries.GetUserByEmailSQL, map[string]interface{}{
		"email": email,
	})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := model.UserDb{}
		if err = rows.StructScan(&u); err != nil {
			return nil, err
		}
		return u.ToUserModel(), nil
	}
	return nil, nil
}

func (u *userRepo) Login(email string, password string) (*model.UserModel, error) {
	rows, err := u.dao.NamedQuery(queries.LoginSQL, map[string]interface{}{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := model.UserDb{}
		if err = rows.StructScan(&u); err != nil {
			return nil, err
		}
		return u.ToUserModel(), nil
	}
	return nil, sql.ErrNoRows
}

func (u *userRepo) InsertUser(user *model.UserModel) (*model.UserModel, error) {
	userDb := user.ToUserDb()
	result, err := u.dao.NamedQuery(queries.InsertUserSQL, map[string]interface{}{
		"email":    userDb.Email,
		"password": userDb.Password,
	})
	if err != nil {
		return nil, err
	}
	result.Close()
	if result.Next() {
		_userDB := model.UserDb{}
		if err = result.StructScan(&_userDB); err != nil {
			return nil, err
		}
		return _userDB.ToUserModel(), nil
	}
	return user, nil
}
