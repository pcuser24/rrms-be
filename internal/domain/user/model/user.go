package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/user/dto"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

type UserModel struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Password  *string    `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	DeletedF  bool       `json:"deleted_f"`
}

// Bang user
type UserDb struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	Password  sql.NullString `json:"password"`
	GroupID   uuid.NullUUID  `json:"group_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy uuid.NullUUID  `json:"created_by"`
	UpdatedBy uuid.NullUUID  `json:"updated_by"`
	// 1: deleted, 0: not deleted
	DeletedF bool `json:"deleted_f"`
}

func (u *UserModel) ToUserDb() *UserDb {
	return &UserDb{
		ID:        u.ID,
		Email:     u.Email,
		Password:  types.StrN(u.Password),
		CreatedAt: u.CreatedAt,
		CreatedBy: types.UUIDN(u.CreatedBy),
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: types.UUIDN(u.UpdatedBy),
		DeletedF:  u.DeletedF,
	}
}

func (u *UserModel) ToUserResponse() *dto.UserResponse {
	ur := &dto.UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		CreatedBy: nil,
		UpdatedBy: nil,
		DeletedF:  u.DeletedF,
	}
	if u.CreatedBy != nil {
		str := u.CreatedBy.String()
		ur.CreatedBy = &str
	}
	if u.UpdatedBy != nil {
		str := u.UpdatedBy.String()
		ur.UpdatedBy = &str
	}
	return ur
}

func (ud *UserDb) ToUserModel() *UserModel {
	return &UserModel{
		ID:        ud.ID,
		Email:     ud.Email,
		Password:  types.PNStr(ud.Password),
		CreatedAt: ud.CreatedAt,
		CreatedBy: types.PNUUID(ud.CreatedBy),
		UpdatedAt: ud.UpdatedAt,
		UpdatedBy: types.PNUUID(ud.UpdatedBy),
		DeletedF:  ud.DeletedF,
	}
}
