package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreatePreRentalCoap struct {
	FullName    *string   `json:"fullName"`
	Dob         time.Time `json:"dob"`
	Job         *string   `json:"job"`
	Income      *int32    `json:"income"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	Description *string   `json:"description"`
}

func (pm *CreatePreRentalCoap) ToCreatePreRentalCoapDB(prerentalID int64) database.CreatePreRentalCoapParams {
	return database.CreatePreRentalCoapParams{
		PrerentalID: prerentalID,
		FullName:    types.StrN(pm.FullName),
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

type CreatePreRental struct {
	ApplicationID   *int64                `json:"applicationId" validate:"omitempty"`
	CreatorID       uuid.UUID             `json:"creatorId" validate:"required"`
	TenantID        uuid.UUID             `json:"tenantId"`
	ProfileImage    string                `json:"profileImage" validate:"required"`
	PropertyID      uuid.UUID             `json:"propertyId" validate:"required"`
	UnitID          uuid.UUID             `json:"unitId" validate:"required"`
	TenantType      database.TENANTTYPE   `json:"tenantType" validate:"required"`
	TenantName      string                `json:"tenantName" validate:"required"`
	TenantDOB       time.Time             `json:"tenantDob" validate:"required"`
	TenantIdentity  string                `json:"tenantIdentity" validate:"required"`
	TenantPhone     string                `json:"tenantPhone" validate:"required"`
	TenantEmail     string                `json:"tenantEmail" validate:"required,email"`
	TenantAddress   *string               `json:"tenantAddress" validate:"omitempty"`
	ContractType    database.CONTRACTTYPE `json:"contractType" validate:"omitempty"`
	ContractContent *string               `json:"contractContent" validate:"omitempty"`
	LandArea        float32               `json:"landArea" validate:"required"`
	UnitArea        float32               `json:"unitArea" validate:"required"`
	StartDate       time.Time             `json:"startDate" validate:"omitempty"`
	MoveinDate      time.Time             `json:"moveinDate" validate:"required"`
	RentalPeriod    int32                 `json:"rentalPeriod" validate:"required"`
	RentalPrice     float32               `json:"rentalPrice" validate:"required"`
	Note            *string               `json:"note" validate:"omitempty"`

	Coaps []CreatePreRentalCoap `json:"coaps" validate:"omitempty,dive"`
}

func (pm *CreatePreRental) ToCreatePreRentalDB() database.CreatePreRentalParams {
	return database.CreatePreRentalParams{
		CreatorID:     pm.CreatorID,
		ApplicationID: types.Int64N(pm.ApplicationID),
		TenantID:      types.UUIDN(pm.TenantID),
		ProfileImage:  pm.ProfileImage,
		PropertyID:    pm.PropertyID,
		UnitID:        pm.UnitID,
		TenantType:    pm.TenantType,
		TenantName:    pm.TenantName,
		TenantDob: pgtype.Date{
			Time:  pm.TenantDOB,
			Valid: !pm.TenantDOB.IsZero(),
		},
		TenantIdentity: pm.TenantIdentity,
		TenantPhone:    pm.TenantPhone,
		TenantEmail:    pm.TenantEmail,
		TenantAddress:  types.StrN(pm.TenantAddress),
		ContractType: database.NullCONTRACTTYPE{
			CONTRACTTYPE: pm.ContractType,
			Valid:        pm.ContractType != "",
		},
		ContractContent: types.StrN(pm.ContractContent),
		LandArea:        pm.LandArea,
		UnitArea:        pm.UnitArea,
		StartDate: pgtype.Date{
			Time:  pm.StartDate,
			Valid: !pm.StartDate.IsZero(),
		},
		MoveinDate: pgtype.Date{
			Time:  pm.MoveinDate,
			Valid: !pm.MoveinDate.IsZero(),
		},
		RentalPeriod: pm.RentalPeriod,
		RentalPrice:  pm.RentalPrice,
		Note:         types.StrN(pm.Note),
	}
}

type PreparePreRentalContract struct {
	ContractType    database.CONTRACTTYPE `json:"contractType" validate:"required"`
	ContractContent *string               `json:"contractContent"`
}
