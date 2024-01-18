package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UserModel struct {
	ID        uuid.UUID
	Email     string
	Password  *string
	GroupID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedF  bool
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
	if u.CreatedBy != uuid.Nil {
		str := u.CreatedBy.String()
		ur.CreatedBy = &str
	}
	if u.UpdatedBy != uuid.Nil {
		str := u.UpdatedBy.String()
		ur.UpdatedBy = &str
	}
	return ur
}

func ToUserModel(ud *database.User) *UserModel {
	return &UserModel{
		ID:        ud.ID,
		Email:     ud.Email,
		Password:  types.PNStr(ud.Password),
		GroupID:   types.PNUUID(ud.GroupID),
		CreatedAt: ud.CreatedAt,
		UpdatedAt: ud.UpdatedAt,
		CreatedBy: types.PNUUID(ud.CreatedBy),
		UpdatedBy: types.PNUUID(ud.UpdatedBy),
		DeletedF:  ud.DeletedF,
	}
}
