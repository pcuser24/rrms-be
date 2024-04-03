package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateRental struct {
	ID                     int64               `json:"id"`
	TenantID               uuid.UUID           `json:"tenant_id"`
	ProfileImage           *string             `json:"profile_image"`
	TenantType             database.TENANTTYPE `json:"tenant_type"`
	TenantName             *string             `json:"tenant_name"`
	TenantPhone            *string             `json:"tenant_phone"`
	TenantEmail            *string             `json:"tenant_email"`
	StartDate              time.Time           `json:"start_date"`
	MoveinDate             time.Time           `json:"movein_date"`
	RentalPeriod           *int32              `json:"rental_period"`
	RentalPrice            *float32            `json:"rental_price"`
	ElectricityPaymentType *string             `json:"electricity_payment_type"`
	ElectricityPrice       *float32            `json:"electricity_price"`
	WaterPaymentType       *string             `json:"water_payment_type"`
	WaterPrice             *float32            `json:"water_price"`
	Note                   *string             `json:"note"`
}

func (pm *UpdateRental) ToUpdateRentalDB(id int64) database.UpdateRentalParams {
	return database.UpdateRentalParams{
		ID:           id,
		TenantID:     types.UUIDN(pm.TenantID),
		ProfileImage: types.StrN(pm.ProfileImage),
		TenantType: database.NullTENANTTYPE{
			TENANTTYPE: pm.TenantType,
			Valid:      pm.TenantType != "",
		},
		TenantName:  types.StrN(pm.TenantName),
		TenantPhone: types.StrN(pm.TenantPhone),
		TenantEmail: types.StrN(pm.TenantEmail),
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
	}
}

// type UpdateRentalContract struct {
// 	ContractType         database.CONTRACTTYPE `json:"contractType" validate:"required"`
// 	ContractContent      *string               `json:"contractContent" validate:"required"`
// 	ContractLastUpdateBy uuid.UUID             `json:"contract_last_update_by" validate:"required"`
// }

// func (pm *UpdateRentalContract) ToUpdateRentalContractDB(id int64) database.UpdateRentalContractParams {
// 	return database.UpdateRentalContractParams{
// 		ID: id,
// 		ContractType: database.NullCONTRACTTYPE{
// 			CONTRACTTYPE: pm.ContractType,
// 			Valid:        pm.ContractType != "",
// 		},
// 		ContractContent: types.StrN(pm.ContractContent),
// 		ContractLastUpdateBy: pgtype.UUID{
// 			Bytes: pm.ContractLastUpdateBy,
// 			Valid: pm.ContractLastUpdateBy != uuid.Nil,
// 		},
// 	}
// }
