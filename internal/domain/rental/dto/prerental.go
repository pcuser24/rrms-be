package dto

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreatePreRental = CreateRental

func (c *CreatePreRental) ToCreatePreRentalDB() (database.CreatePreRentalParams, error) {
	cdb := database.CreatePreRentalParams{
		CreatorID:             c.CreatorID,
		ApplicationID:         types.Int64N(c.ApplicationID),
		TenantID:              types.UUIDN(c.TenantID),
		ProfileImage:          c.ProfileImage,
		PropertyID:            c.PropertyID,
		UnitID:                c.UnitID,
		TenantType:            c.TenantType,
		TenantName:            c.TenantName,
		TenantPhone:           c.TenantPhone,
		TenantEmail:           c.TenantEmail,
		OrganizationName:      types.StrN(c.OrganizationName),
		OrganizationHqAddress: types.StrN(c.OrganizationHqAddress),
		StartDate: pgtype.Date{
			Time:  c.StartDate,
			Valid: !c.StartDate.IsZero(),
		},
		MoveinDate: pgtype.Date{
			Time:  c.MoveinDate,
			Valid: !c.MoveinDate.IsZero(),
		},
		RentalPeriod: c.RentalPeriod,
		PaymentType: database.NullRENTALPAYMENTTYPE{
			RENTALPAYMENTTYPE: c.PaymentType,
			Valid:             c.PaymentType != "",
		},
		RentalPrice:        c.RentalPrice,
		RentalPaymentBasis: c.RentalPaymentBasis,
		RentalIntention:    c.RentalIntention,
		NoticePeriod:       types.Int32N(c.NoticePeriod),
		GracePeriod: pgtype.Int4{
			Int32: c.GracePeriod,
			Valid: c.GracePeriod != 0,
		},
		LatePaymentPenaltyScheme: database.NullLATEPAYMENTPENALTYSCHEME{
			LATEPAYMENTPENALTYSCHEME: c.LatePaymentPenaltyScheme,
			Valid:                    c.LatePaymentPenaltyScheme != "",
		},
		LatePaymentPenaltyAmount: types.Float32N(c.LatePaymentPenaltyAmount),
		ElectricitySetupBy:       c.ElectricitySetupBy,
		ElectricityPaymentType:   types.StrN(c.ElectricityPaymentType),
		ElectricityPrice:         types.Float32N(c.ElectricityPrice),
		ElectricityCustomerCode:  types.StrN(c.ElectricityCustomerCode),
		ElectricityProvider:      types.StrN(c.ElectricityProvider),
		WaterSetupBy:             c.WaterSetupBy,
		WaterPaymentType:         types.StrN(c.WaterPaymentType),
		WaterPrice:               types.Float32N(c.WaterPrice),
		WaterCustomerCode:        types.StrN(c.WaterCustomerCode),
		WaterProvider:            types.StrN(c.WaterProvider),
		Note:                     types.StrN(c.Note),
	}
	if len(c.Coaps) > 0 {
		b, err := json.Marshal(c.Coaps)
		if err != nil {
			return cdb, err
		}
		cdb.Coaps = b
	}
	if len(c.Minors) > 0 {
		b, err := json.Marshal(c.Minors)
		if err != nil {
			return cdb, err
		}
		cdb.Minors = b
	}
	if len(c.Pets) > 0 {
		b, err := json.Marshal(c.Pets)
		if err != nil {
			return cdb, err
		}
		cdb.Pets = b
	}
	if len(c.Services) > 0 {
		b, err := json.Marshal(c.Services)
		if err != nil {
			return cdb, err
		}
		cdb.Services = b
	}
	if len(c.Policies) > 0 {
		b, err := json.Marshal(c.Policies)
		if err != nil {
			return cdb, err
		}
		cdb.Policies = b
	}

	return cdb, nil
}

type GetPreRentalResponse struct {
	PreRental *rental_model.PreRental       `json:"preRental"`
	Property  *property_model.PropertyModel `json:"property"`
	Unit      *unit_model.UnitModel         `json:"unit"`
}

type UpdatePreRental struct {
	Feedback *string `json:"feedback" validate:"omitempty"`
	State    string  `json:"state" validate:"required,oneof=APPROVED REVIEW REJECTED"`
}

type GetPreRentalsQuery struct {
	Limit  int32 `json:"limit" validate:"omitempty,min=1"`
	Offset int32 `json:"offset" validate:"omitempty,min=0"`
}
