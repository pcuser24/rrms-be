package model

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/pkg/utils/types"
	"time"
)

type UserModel struct {
	ID        uuid.UUID  `json:"id"`
	Email     *string    `json:"email"`
	Password  *string    `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy *int64     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int64     `json:"updated_by"`
	DeletedF  bool       `json:"deleted_f"`
}

type UserDb struct {
	ID        sql.NullString `db:"id"`
	Email     sql.NullString `db:"email"`
	Password  sql.NullString `db:"password"`
	CreatedAt sql.NullString `db:"created_at"`
	CreatedBy sql.NullInt64  `db:"created_by"`
	UpdatedAt sql.NullString `db:"updated_at"`
	UpdatedBy sql.NullInt64  `db:"updated_by"`
	DeletedF  sql.NullBool   `db:"deleted_f"`
}

func (u *UserModel) ToUserDb() *UserDb {
	return &UserDb{
		ID:        types.UUIDN(u.ID),
		Email:     types.StrN(u.Email),
		Password:  types.StrN(u.Password),
		CreatedAt: types.TimeNStr(u.CreatedAt),
		CreatedBy: types.Int64N(u.CreatedBy),
		UpdatedAt: types.TimeNStr(u.UpdatedAt),
		UpdatedBy: types.Int64N(u.UpdatedBy),
		DeletedF:  sql.NullBool{Bool: u.DeletedF, Valid: true},
	}
}

func (ud *UserDb) ToUserModel() *UserModel {
	return &UserModel{
		ID:        types.NUUID(ud.ID),
		Email:     types.PNStr(ud.Email),
		Password:  types.PNStr(ud.Password),
		CreatedAt: types.NStrTime(ud.CreatedAt),
		CreatedBy: types.PNInt64(ud.CreatedBy),
		UpdatedAt: types.NStrTime(ud.UpdatedAt),
		UpdatedBy: types.PNInt64(ud.UpdatedBy),
		DeletedF:  types.NBool(ud.DeletedF),
	}
}
