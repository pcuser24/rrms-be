package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateRentalCoap struct {
	FullName    *string   `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"omitempty"`
	Job         *string   `json:"job" validate:"omitempty"`
	Income      *int32    `json:"income" validate:"omitempty"`
	Email       *string   `json:"email" validate:"omitempty"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (pm *CreateRentalCoap) ToCreateRentalCoapDB(rentalID int64) database.CreateRentalCoapParams {
	return database.CreateRentalCoapParams{
		RentalID: rentalID,
		FullName: types.StrN(pm.FullName),
		Dob: pgtype.Date{
			Time:  pm.Dob,
			Valid: !pm.Dob.IsZero(),
		},
		Job:         types.StrN(pm.Job),
		Income:      types.Int32N(pm.Income),
		Email:       types.StrN(pm.Email),
		Phone:       types.StrN(pm.Phone),
		Description: types.StrN(pm.Description),
	}
}

type CreateRentalMinor struct {
	FullName    string    `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"omitempty"`
	Email       *string   `json:"email" validate:"omitempty"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (pm *CreateRentalMinor) ToCreateRentalMinorDB(id int64) database.CreateRentalMinorParams {
	return database.CreateRentalMinorParams{
		RentalID: id,
		FullName: pm.FullName,
		Dob: pgtype.Date{
			Time:  pm.Dob,
			Valid: !pm.Dob.IsZero(),
		},
		Email:       types.StrN(pm.Email),
		Phone:       types.StrN(pm.Phone),
		Description: types.StrN(pm.Description),
	}
}

type CreateRentalPet struct {
	Type        string   `json:"type"`
	Weight      *float32 `json:"weight" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
}

func (pm *CreateRentalPet) ToCreateRentalPetDB(id int64) database.CreateRentalPetParams {
	return database.CreateRentalPetParams{
		RentalID:    id,
		Type:        pm.Type,
		Weight:      types.Float32N(pm.Weight),
		Description: types.StrN(pm.Description),
	}
}

type CreateRentalService struct {
	Name     string   `json:"name" validate:"required"`
	Setupby  string   `json:"setupby" validate:"required"`
	Provider *string  `json:"provider" validate:"omitempty"`
	Price    *float32 `json:"price" validate:"omitempty"`
}

func (pm *CreateRentalService) ToCreateRentalServiceDB(id int64) database.CreateRentalServiceParams {
	return database.CreateRentalServiceParams{
		RentalID: id,
		Name:     pm.Name,
		Setupby:  pm.Setupby,
		Provider: types.StrN(pm.Provider),
		Price:    types.Float32N(pm.Price),
	}
}

type CreateRental struct {
	ApplicationID          *int64 `json:"applicationId" validate:"omitempty"`
	CreatorID              uuid.UUID
	TenantID               uuid.UUID           `json:"tenantId" validate:"omitempty"`
	PropertyID             uuid.UUID           `json:"propertyId" validatet:"required"`
	UnitID                 uuid.UUID           `json:"unitId" validatet:"required"`
	ProfileImage           string              `json:"profileImage" validate:"omitempty"`
	TenantType             database.TENANTTYPE `json:"tenantType" validate:"required"`
	TenantName             string              `json:"tenantName" validate:"required"`
	TenantPhone            string              `json:"tenantPhone" validate:"required"`
	TenantEmail            string              `json:"tenantEmail" validate:"required"`
	OrganizationName       *string             `json:"organizationName" validate:"omitempty"`
	OrganizationHqAddress  *string             `json:"organizationHqAddress" validate:"omitempty"`
	StartDate              time.Time           `json:"startDate" validate:"required"`
	MoveinDate             time.Time           `json:"moveinDate" validate:"required"`
	RentalPeriod           int32               `json:"rentalPeriod" validate:"required"`
	RentalPrice            float32             `json:"rentalPrice" validate:"required"`
	RentalIntention        string              `json:"rentalIntention" validate:"required"`
	Deposit                float32             `json:"deposit" validate:"required"`
	DepositPaid            bool                `json:"depositPaid" validate:"required"`
	ElectricityPaymentType string              `json:"electricityPaymentType" validate:"required"`
	ElectricityPrice       *float32            `json:"electricityPrice" validate:"omitempty"`
	WaterPaymentType       string              `json:"waterPaymentType" validate:"required"`
	WaterPrice             *float32            `json:"waterPrice" validate:"omitempty"`
	Note                   *string             `json:"note"`

	Coaps    []CreateRentalCoap    `json:"coaps"`
	Minors   []CreateRentalMinor   `json:"minors"`
	Pets     []CreateRentalPet     `json:"pets"`
	Services []CreateRentalService `json:"services"`
}

func (pm *CreateRental) ToCreateRentalDB() database.CreateRentalParams {
	return database.CreateRentalParams{
		CreatorID:             pm.CreatorID,
		ApplicationID:         types.Int64N(pm.ApplicationID),
		TenantID:              types.UUIDN(pm.TenantID),
		ProfileImage:          pm.ProfileImage,
		PropertyID:            pm.PropertyID,
		UnitID:                pm.UnitID,
		TenantType:            pm.TenantType,
		TenantName:            pm.TenantName,
		TenantPhone:           pm.TenantPhone,
		TenantEmail:           pm.TenantEmail,
		OrganizationName:      types.StrN(pm.OrganizationName),
		OrganizationHqAddress: types.StrN(pm.OrganizationHqAddress),
		StartDate: pgtype.Date{
			Time:  pm.StartDate,
			Valid: !pm.StartDate.IsZero(),
		},
		MoveinDate: pgtype.Date{
			Time:  pm.MoveinDate,
			Valid: !pm.MoveinDate.IsZero(),
		},
		RentalPeriod:           pm.RentalPeriod,
		RentalPrice:            pm.RentalPrice,
		RentalIntention:        pm.RentalIntention,
		Deposit:                pm.Deposit,
		DepositPaid:            pm.DepositPaid,
		ElectricityPaymentType: pm.ElectricityPaymentType,
		ElectricityPrice:       types.Float32N(pm.ElectricityPrice),
		WaterPaymentType:       pm.WaterPaymentType,
		WaterPrice:             types.Float32N(pm.WaterPrice),
		Note:                   types.StrN(pm.Note),
	}
}
