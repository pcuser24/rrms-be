package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateContract struct {
	ID                        int64     `json:"id" validate:"required"`
	AFullname                 *string   `json:"aFullname" validate:"omitempty"`
	ADob                      time.Time `json:"aDob" validate:"omitempty"`
	APhone                    *string   `json:"aPhone" validate:"omitempty"`
	AAddress                  *string   `json:"aAddress" validate:"omitempty"`
	AHouseholdRegistration    *string   `json:"aHouseholdRegistration" validate:"omitempty"`
	AIdentity                 *string   `json:"aIdentity" validate:"omitempty"`
	AIdentityIssuedBy         *string   `json:"aIdentityIssuedBy" validate:"omitempty"`
	AIdentityIssuedAt         time.Time `json:"aIdentityIssuedAt" validate:"omitempty"`
	ADocuments                []string  `json:"aDocuments" validate:"omitempty"`
	ABankAccount              *string   `json:"aBankAccount" validate:"omitempty"`
	ABank                     *string   `json:"aBank" validate:"omitempty"`
	BFullname                 *string   `json:"bFullname" validate:"omitempty"`
	BOrganizationName         *string   `json:"bOrganizationName" validate:"omitempty"`
	BOrganizationHqAddress    *string   `json:"bOrganizationHqAddress" validate:"omitempty"`
	BOrganizationCode         *string   `json:"bOrganizationCode" validate:"omitempty"`
	BOrganizationCodeIssuedAt time.Time `json:"bOrganizationCodeIssuedAt" validate:"omitempty"`
	BOrganizationCodeIssuedBy *string   `json:"bOrganizationCodeIssuedBy" validate:"omitempty"`
	BDob                      *string   `json:"bDob" validate:"omitempty"`
	BPhone                    *string   `json:"bPhone" validate:"omitempty"`
	BAddress                  *string   `json:"bAddress" validate:"omitempty"`
	BHouseholdRegistration    *string   `json:"bHouseholdRegistration" validate:"omitempty"`
	BIdentity                 *string   `json:"bIdentity" validate:"omitempty"`
	BIdentityIssuedBy         *string   `json:"bIdentityIssuedBy" validate:"omitempty"`
	BIdentityIssuedAt         time.Time `json:"bIdentityIssuedAt" validate:"omitempty"`
	BBankAccount              *string   `json:"bBankAccount" validate:"omitempty"`
	BBank                     *string   `json:"bBank" validate:"omitempty"`
	BTaxCode                  *string   `json:"bTaxCode" validate:"omitempty"`
	PaymentMethod             *string   `json:"paymentMethod" validate:"omitempty"`
	PaymentDay                *int32    `json:"paymentDay" validate:"omitempty"`
	Content                   *string   `json:"content" validate:"omitempty"`
	UserID                    uuid.UUID `json:"userId" validate:"required"`
}

func (c *UpdateContract) ToUpdateContractDB() database.UpdateContractParams {
	return database.UpdateContractParams{
		ID:                        c.ID,
		AFullname:                 types.StrN(c.AFullname),
		ADob:                      types.DateN(c.ADob),
		APhone:                    types.StrN(c.APhone),
		AAddress:                  types.StrN(c.AAddress),
		AHouseholdRegistration:    types.StrN(c.AHouseholdRegistration),
		AIdentity:                 types.StrN(c.AIdentity),
		AIdentityIssuedBy:         types.StrN(c.AIdentityIssuedBy),
		AIdentityIssuedAt:         types.DateN(c.AIdentityIssuedAt),
		ADocuments:                c.ADocuments,
		ABankAccount:              types.StrN(c.ABankAccount),
		ABank:                     types.StrN(c.ABank),
		BFullname:                 types.StrN(c.BFullname),
		BOrganizationName:         types.StrN(c.BOrganizationName),
		BOrganizationHqAddress:    types.StrN(c.BOrganizationHqAddress),
		BOrganizationCode:         types.StrN(c.BOrganizationCode),
		BOrganizationCodeIssuedAt: types.DateN(c.BOrganizationCodeIssuedAt),
		BOrganizationCodeIssuedBy: types.StrN(c.BOrganizationCodeIssuedBy),
		BDob:                      types.StrN(c.BDob),
		BPhone:                    types.StrN(c.BPhone),
		BAddress:                  types.StrN(c.BAddress),
		BHouseholdRegistration:    types.StrN(c.BHouseholdRegistration),
		BIdentity:                 types.StrN(c.BIdentity),
		BIdentityIssuedBy:         types.StrN(c.BIdentityIssuedBy),
		BIdentityIssuedAt:         types.DateN(c.BIdentityIssuedAt),
		BBankAccount:              types.StrN(c.BBankAccount),
		BBank:                     types.StrN(c.BBank),
		BTaxCode:                  types.StrN(c.BTaxCode),
		PaymentMethod:             types.StrN(c.PaymentMethod),
		PaymentDay:                types.Int32N(c.PaymentDay),
		Content:                   types.StrN(c.Content),
		UserID:                    c.UserID,
	}
}

type UpdateContractContent struct {
	ID      int64                   `json:"id"`
	Content *string                 `json:"content"`
	Status  database.CONTRACTSTATUS `json:"status"`
	UserID  uuid.UUID               `json:"userId"`
}

func (c *UpdateContractContent) ToUpdateContractContentDB() database.UpdateContractContentParams {
	return database.UpdateContractContentParams{
		ID:      c.ID,
		Content: types.StrN(c.Content),
		Status:  c.Status,
		UserID:  c.UserID,
	}
}
