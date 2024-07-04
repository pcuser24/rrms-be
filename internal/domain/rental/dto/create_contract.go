package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateContract struct {
	RentalID               int64     `json:"rentalId" validate:"required"`
	AFullname              string    `json:"aFullname" validate:"required"`
	ADob                   time.Time `json:"aDob" validate:"required"`
	APhone                 string    `json:"aPhone" validate:"required"`
	AAddress               string    `json:"aAddress" validate:"required"`
	AHouseholdRegistration string    `json:"aHouseholdRegistration" validate:"required"`
	AIdentity              string    `json:"aIdentity" validate:"required"`
	AIdentityIssuedBy      string    `json:"aIdentityIssuedBy" validate:"required"`
	AIdentityIssuedAt      time.Time `json:"aIdentityIssuedAt" validate:"required"`
	ADocuments             []string  `json:"aDocuments" validate:"omitempty"`
	ABankAccount           *string   `json:"aBankAccount" validate:"omitempty"`
	ABank                  *string   `json:"aBank" validate:"omitempty"`
	ARegistrationNumber    string    `json:"aRegistrationNumber" validate:"required"`
	// BFullname              string    `json:"bFullname" validate:"required"`
	// BPhone                 string    `json:"bPhone" validate:"required"`
	PaymentMethod  string  `json:"paymentMethod" validate:"required"`
	NCopies        int32   `json:"nCopies" validate:"required"`
	CreatedAtPlace string  `json:"createdAtPlace" validate:"required"`
	Content        *string `json:"content" validate:"omitempty"`
	UserID         uuid.UUID
}

func (c *CreateContract) ToCreateContractDB() database.CreateContractParams {
	content := pgtype.Text{
		Valid: true,
	}
	if c.Content != nil {
		content.String = *c.Content
	} else {
		content.String = ""
	}

	return database.CreateContractParams{
		RentalID:  c.RentalID,
		AFullname: c.AFullname,
		ADob: pgtype.Date{
			Time:  c.ADob,
			Valid: !c.ADob.IsZero(),
		},
		APhone:                 c.APhone,
		AAddress:               c.AAddress,
		AHouseholdRegistration: c.AHouseholdRegistration,
		AIdentity:              c.AIdentity,
		AIdentityIssuedBy:      c.AIdentityIssuedBy,
		AIdentityIssuedAt: pgtype.Date{
			Time:  c.AIdentityIssuedAt,
			Valid: !c.AIdentityIssuedAt.IsZero(),
		},
		ADocuments:          c.ADocuments,
		ABankAccount:        types.StrN(c.ABankAccount),
		ABank:               types.StrN(c.ABank),
		ARegistrationNumber: c.ARegistrationNumber,
		// BFullname:           c.BFullname,
		// BPhone:              c.BPhone,
		PaymentMethod:  c.PaymentMethod,
		NCopies:        c.NCopies,
		CreatedAtPlace: c.CreatedAtPlace,
		Content:        content,
		UserID:         c.UserID,
	}
}
