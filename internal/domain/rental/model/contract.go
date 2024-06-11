package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type ContractModel struct {
	ID                     int64     `json:"id"`
	RentalID               int64     `json:"rentalId"`
	AFullname              string    `json:"aFullname"`
	ADob                   time.Time `json:"aDob"`
	APhone                 string    `json:"aPhone"`
	AAddress               string    `json:"aAddress"`
	AHouseholdRegistration string    `json:"aHouseholdRegistration"`
	AIdentity              string    `json:"aIdentity"`
	AIdentityIssuedBy      string    `json:"aIdentityIssuedBy"`
	AIdentityIssuedAt      time.Time `json:"aIdentityIssuedAt"`
	ADocuments             []string  `json:"aDocuments"`
	ABankAccount           *string   `json:"aBankAccount"`
	ABank                  *string   `json:"aBank"`
	ARegistrationNumber    string    `json:"aRegistrationNumber"`
	BFullname              string    `json:"bFullName"`

	BOrganizationName         *string   `json:"bOrganizationName"`
	BOrganizationHqAddress    *string   `json:"bOrganizationHqAddress"`
	BOrganizationCode         *string   `json:"bOrganizationCode"`
	BOrganizationCodeIssuedAt time.Time `json:"bOrganizationCodeIssuedAt"`
	BOrganizationCodeIssuedBy *string   `json:"bOrganizationCodeIssuedBy"`
	BDob                      *string   `json:"bDob"`
	BPhone                    string    `json:"bPhone"`
	BAddress                  *string   `json:"bAddress"`
	BHouseholdRegistration    *string   `json:"bHouseholdRegistration"`
	BIdentity                 *string   `json:"bIdentity"`
	BIdentityIssuedBy         *string   `json:"bIdentityIssuedBy"`
	BIdentityIssuedAt         time.Time `json:"bIdentityIssuedAt"`
	BBankAccount              *string   `json:"bBankAccount"`
	BBank                     *string   `json:"bBank"`
	BTaxCode                  *string   `json:"bTaxCode"`

	PaymentMethod  string                  `json:"paymentMethod"`
	NCopies        int32                   `json:"nCopies"`
	CreatedAtPlace string                  `json:"createdAtPlace"`
	Content        string                  `json:"content"`
	Status         database.CONTRACTSTATUS `json:"status"`
	CreatedAt      time.Time               `json:"createdAt"`
	CreatedBy      uuid.UUID               `json:"createdBy"`
	UpdatedAt      time.Time               `json:"updatedAt"`
	UpdatedBy      uuid.UUID               `json:"updatedBy"`
}

func ToContractModel(db *database.Contract) *ContractModel {
	return &ContractModel{
		ID:                        db.ID,
		RentalID:                  db.RentalID,
		AFullname:                 db.AFullname,
		ADob:                      db.ADob.Time,
		APhone:                    db.APhone,
		AAddress:                  db.AAddress,
		AHouseholdRegistration:    db.AHouseholdRegistration,
		AIdentity:                 db.AIdentity,
		AIdentityIssuedBy:         db.AIdentityIssuedBy,
		AIdentityIssuedAt:         db.AIdentityIssuedAt.Time,
		ADocuments:                db.ADocuments,
		ABankAccount:              types.PNStr(db.ABankAccount),
		ABank:                     types.PNStr(db.ABank),
		ARegistrationNumber:       db.ARegistrationNumber,
		BFullname:                 db.BFullname,
		BOrganizationName:         types.PNStr(db.BOrganizationName),
		BOrganizationHqAddress:    types.PNStr(db.BOrganizationHqAddress),
		BOrganizationCode:         types.PNStr(db.BOrganizationCode),
		BOrganizationCodeIssuedAt: db.BOrganizationCodeIssuedAt.Time,
		BOrganizationCodeIssuedBy: types.PNStr(db.BOrganizationCodeIssuedBy),
		BDob:                      types.PNStr(db.BDob),
		BPhone:                    db.BPhone,
		BAddress:                  types.PNStr(db.BAddress),
		BHouseholdRegistration:    types.PNStr(db.BHouseholdRegistration),
		BIdentity:                 types.PNStr(db.BIdentity),
		BIdentityIssuedBy:         types.PNStr(db.BIdentityIssuedBy),
		BIdentityIssuedAt:         db.BIdentityIssuedAt.Time,
		BBankAccount:              types.PNStr(db.BBankAccount),
		BBank:                     types.PNStr(db.BBank),
		BTaxCode:                  types.PNStr(db.BTaxCode),
		PaymentMethod:             db.PaymentMethod,
		NCopies:                   db.NCopies,
		CreatedAtPlace:            db.CreatedAtPlace,
		Content:                   db.Content,
		Status:                    db.Status,
		CreatedAt:                 db.CreatedAt,
		CreatedBy:                 db.CreatedBy,
		UpdatedAt:                 db.UpdatedAt,
		UpdatedBy:                 db.UpdatedBy,
	}
}
