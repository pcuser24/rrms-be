package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type UpdateRental struct {
	ID                 int64               `json:"id" validate:"required"`
	TenantID           uuid.UUID           `json:"tenantId" validate:"omitempty"`
	ProfileImage       *string             `json:"profileImage" validate:"omitempty"`
	TenantType         database.TENANTTYPE `json:"tenantType" validate:"omitempty,oneof=INDIVIDUAL FAMILY ORGANIZATION"`
	TenantName         *string             `json:"tenantName" validate:"omitempty"`
	TenantPhone        *string             `json:"tenantPhone" validate:"omitempty"`
	TenantEmail        *string             `json:"tenantEmail" validate:"omitempty"`
	StartDate          time.Time           `json:"startDate" validate:"omitempty"`
	MoveinDate         time.Time           `json:"moveinDate" validate:"omitempty"`
	RentalPeriod       *int32              `json:"rentalPeriod" validate:"omitempty"`
	RentalPrice        *float32            `json:"rentalPrice" validate:"omitempty"`
	RentalPaymentBasis *int32              `json:"rentalPaymentBasis" validate:"omitempty"`

	ElectricitySetupBy             *string               `json:"electricitySetupBy" validate:"omitempty"`
	ElectricityPaymentType         *string               `json:"electricityPaymentType" validate:"omitempty"`
	ElectricityProvider            *string               `json:"electricityProvider" validate:"omitempty"`
	ElectricityCustomerCode        *string               `json:"electricityCustomerCode" validate:"omitempty"`
	ElectricityPrice               *float32              `json:"electricityPrice" validate:"omitempty"`
	WaterSetupBy                   *string               `json:"waterSetupBy" validate:"omitempty"`
	WaterPaymentType               *string               `json:"waterPaymentType" validate:"omitempty"`
	WaterCustomerCode              *string               `json:"waterCustomerCode" validate:"omitempty"`
	WaterProvider                  *string               `json:"waterProvider" validate:"omitempty"`
	WaterPrice                     *float32              `json:"waterPrice" validate:"omitempty"`
	RentalPaymentGracePeriod       *int32                `json:"rentalPaymentGracePeriod" validate:"omitempty"`
	RentalPaymentLateFeePercentage *float32              `json:"rentalPaymentLateFeePercentage" validate:"omitempty"`
	Status                         database.RENTALSTATUS `json:"status" validate:"omitempty,oneof=INPROGRESS END"`
	Note                           *string               `json:"note" validate:"omitempty"`
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
		RentalPeriod:            types.Int32N(pm.RentalPeriod),
		RentalPrice:             types.Float32N(pm.RentalPrice),
		RentalPaymentBasis:      types.Int32N(pm.RentalPaymentBasis),
		ElectricitySetupBy:      types.StrN(pm.ElectricitySetupBy),
		ElectricityPaymentType:  types.StrN(pm.ElectricityPaymentType),
		ElectricityProvider:     types.StrN(pm.ElectricityProvider),
		ElectricityCustomerCode: types.StrN(pm.ElectricityCustomerCode),
		ElectricityPrice:        types.Float32N(pm.ElectricityPrice),
		WaterSetupBy:            types.StrN(pm.WaterSetupBy),
		WaterPaymentType:        types.StrN(pm.WaterPaymentType),
		WaterCustomerCode:       types.StrN(pm.WaterCustomerCode),
		WaterProvider:           types.StrN(pm.WaterProvider),
		WaterPrice:              types.Float32N(pm.WaterPrice),
		Note:                    types.StrN(pm.Note),
		Status: database.NullRENTALSTATUS{
			RENTALSTATUS: pm.Status,
			Valid:        pm.Status != "",
		},
	}
}
