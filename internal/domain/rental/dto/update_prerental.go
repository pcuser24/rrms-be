package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdatePreRental struct {
	TenantID       uuid.UUID                `json:"tenantId" validate:"omitempty"`
	ProfileImage   *string                  `json:"profileImage" validate:"omitempty"`
	TenantType     database.TENANTTYPE      `json:"tenantType" validate:"omitempty"`
	TenantName     *string                  `json:"tenantName" validate:"omitempty"`
	TenantDOB      time.Time                `json:"tenantDob" validate:"omitempty"`
	TenantIdentity *string                  `json:"tenantIdentity" validate:"omitempty"`
	TenantPhone    *string                  `json:"tenantPhone" validate:"omitempty"`
	TenantEmail    *string                  `json:"tenantEmail" validate:"omitempty"`
	TenantAddress  *string                  `json:"tenantAddress" validate:"omitempty"`
	LandArea       *float32                 `json:"landArea" validate:"omitempty"`
	UnitArea       *float32                 `json:"unitArea" validate:"omitempty"`
	StartDate      time.Time                `json:"startDate" validate:"omitempty"`
	MoveinDate     time.Time                `json:"moveinDate" validate:"omitempty"`
	RentalPeriod   *int32                   `json:"rentalPeriod" validate:"omitempty"`
	RentalPrice    *float32                 `json:"rentalPrice" validate:"omitempty"`
	Note           *string                  `json:"note" validate:"omitempty"`
	Status         database.PRERENTALSTATUS `json:"status" validate:"omitempty"`
}

func (pm *UpdatePreRental) ToUpdatePreRentalDB(id int64) database.UpdatePreRentalParams {
	return database.UpdatePreRentalParams{
		ID:           id,
		TenantID:     types.UUIDN(pm.TenantID),
		ProfileImage: types.StrN(pm.ProfileImage),
		TenantType: database.NullTENANTTYPE{
			TENANTTYPE: pm.TenantType,
			Valid:      pm.TenantType != "",
		},
		TenantName: types.StrN(pm.TenantName),
		TenantDob: pgtype.Date{
			Time:  pm.TenantDOB,
			Valid: !pm.TenantDOB.IsZero(),
		},
		TenantIdentity: types.StrN(pm.TenantIdentity),
		TenantPhone:    types.StrN(pm.TenantPhone),
		TenantEmail:    types.StrN(pm.TenantEmail),
		TenantAddress:  types.StrN(pm.TenantAddress),
		StartDate: pgtype.Date{
			Time:  pm.StartDate,
			Valid: !pm.StartDate.IsZero(),
		},
		MoveinDate: pgtype.Date{
			Time:  pm.MoveinDate,
			Valid: !pm.MoveinDate.IsZero(),
		},
		RentalPeriod: types.Int32N(pm.RentalPeriod),
		RentalPrice:  types.Float32N(pm.RentalPrice),
		Note:         types.StrN(pm.Note),
		Status: database.NullPRERENTALSTATUS{
			PRERENTALSTATUS: pm.Status,
			Valid:           pm.Status != "",
		},
	}
}

type UpdatePreRentalContract struct {
	ContractType         database.CONTRACTTYPE `json:"contractType" validate:"required"`
	ContractContent      *string               `json:"contractContent" validate:"required"`
	ContractLastUpdateBy uuid.UUID             `json:"contract_last_update_by" validate:"required"`
}

func (pm *UpdatePreRentalContract) ToUpdatePreRentalContractDB(id int64) database.UpdatePreRentalContractParams {
	return database.UpdatePreRentalContractParams{
		ID: id,
		ContractType: database.NullCONTRACTTYPE{
			CONTRACTTYPE: pm.ContractType,
			Valid:        pm.ContractType != "",
		},
		ContractContent: types.StrN(pm.ContractContent),
		ContractLastUpdateBy: pgtype.UUID{
			Bytes: pm.ContractLastUpdateBy,
			Valid: pm.ContractLastUpdateBy != uuid.Nil,
		},
	}
}
