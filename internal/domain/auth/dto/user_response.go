package dto

import (
	"time"

	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

// UserModel without password
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// CreatedBy *string   `json:"created_by"`
	// UpdatedBy *string   `json:"updated_by"`
	DeletedF bool `json:"deleted_f"`

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
