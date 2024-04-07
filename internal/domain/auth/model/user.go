package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UserModel struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  *string
	GroupID   uuid.UUID `json:"groupId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedBy uuid.UUID `json:"updatedBy"`
	DeletedF  bool      `json:"deletedF"`

	// User info
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Phone     *string           `json:"phone"`
	Avatar    *string           `json:"avatar"`
	Address   *string           `json:"address"`
	City      *string           `json:"city"`
	District  *string           `json:"district"`
	Ward      *string           `json:"ward"`
	Role      database.USERROLE `json:"role"`
}

func (u *UserModel) ToUserResponse() *dto.UserResponse {
	ur := &dto.UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		// CreatedBy: nil,
		// UpdatedBy: nil,
		DeletedF:  u.DeletedF,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Address:   u.Address,
		City:      u.City,
		District:  u.District,
		Ward:      u.Ward,
		Role:      u.Role,
	}
	// if u.CreatedBy != uuid.Nil {
	// 	str := u.CreatedBy.String()
	// 	ur.CreatedBy = &str
	// }
	// if u.UpdatedBy != uuid.Nil {
	// 	str := u.UpdatedBy.String()
	// 	ur.UpdatedBy = &str
	// }
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
		FirstName: ud.FirstName,
		LastName:  ud.LastName,
		Phone:     types.PNStr(ud.Phone),
		Avatar:    types.PNStr(ud.Avatar),
		Address:   types.PNStr(ud.Address),
		City:      types.PNStr(ud.City),
		District:  types.PNStr(ud.District),
		Ward:      types.PNStr(ud.Ward),
		Role:      ud.Role,
	}
}
