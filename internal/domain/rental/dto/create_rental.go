package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/internal/utils/validation"
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
	Setupby  string   `json:"setupby" validate:"required,oneof=LANDLORD TENANT"`
	Provider *string  `json:"provider" validate:"omitempty"`
	Price    *float32 `json:"price" validate:"omitempty"`
}

func (pm *CreateRentalService) ToCreateRentalServiceDB(id int64) database.CreateRentalServiceParams {
	return database.CreateRentalServiceParams{
		RentalID: id,
		Name:     pm.Name,
		SetupBy:  pm.Setupby,
		Provider: types.StrN(pm.Provider),
		Price:    types.Float32N(pm.Price),
	}
}

type CreateRentalPolicy struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (pm *CreateRentalPolicy) ToCreateRentalPolicyDB(id int64) database.CreateRentalPolicyParams {
	return database.CreateRentalPolicyParams{
		RentalID: id,
		Title:    pm.Title,
		Content:  pm.Content,
	}
}

type PreCreateRentalMedia struct {
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required,gt=0"`
	Type string `json:"type" validate:"required"`
	Url  string `json:"url"`
}

type PreCreateRental struct {
	Avatar PreCreateRentalMedia `json:"avatar" validate:"required"`
}

type CreateRental struct {
	ApplicationID            *int64 `json:"applicationId" validate:"omitempty"`
	CreatorID                uuid.UUID
	TenantID                 uuid.UUID                         `json:"tenantId" validate:"omitempty"`
	PropertyID               uuid.UUID                         `json:"propertyId" validatet:"required"`
	UnitID                   uuid.UUID                         `json:"unitId" validatet:"required"`
	ProfileImage             string                            `json:"profileImage" validate:"omitempty"`
	TenantType               database.TENANTTYPE               `json:"tenantType" validate:"required,oneof=INDIVIDUAL FAMILY ORGANIZATION"`
	TenantName               string                            `json:"tenantName" validate:"required"`
	TenantPhone              string                            `json:"tenantPhone" validate:"required"`
	TenantEmail              string                            `json:"tenantEmail" validate:"required"`
	OrganizationName         *string                           `json:"organizationName" validate:"omitempty"`
	OrganizationHqAddress    *string                           `json:"organizationHqAddress" validate:"omitempty"`
	StartDate                time.Time                         `json:"startDate" validate:"required"`
	MoveinDate               time.Time                         `json:"moveinDate" validate:"required"`
	RentalPeriod             int32                             `json:"rentalPeriod" validate:"required"`
	PaymentType              database.RENTALPAYMENTTYPE        `json:"paymentType" validate:"required,oneof=PREPAID POSTPAID"`
	RentalPrice              float32                           `json:"rentalPrice" validate:"required"`
	RentalPaymentBasis       int32                             `json:"rentalPaymentBasis" validate:"required"`
	RentalIntention          string                            `json:"rentalIntention" validate:"required"`
	NoticePeriod             *int32                            `json:"noticePeriod" validate:"omitempty,gte=0"`
	GracePeriod              int32                             `json:"gracePeriod" validate:"required,gt=0"`
	LatePaymentPenaltyScheme database.LATEPAYMENTPENALTYSCHEME `json:"latePaymentPenaltyScheme" validate:"required,oneof=FIXED PERCENT NONE"`
	LatePaymentPenaltyAmount *float32                          `json:"latePaymentPenaltyAmount" validate:"omitempty,gte=0"`
	ElectricitySetupBy       string                            `json:"electricitySetupBy" validate:"required,oneof=LANDLORD TENANT"`
	ElectricityPaymentType   *string                           `json:"electricityPaymentType" validate:"omitempty,oneof=RETAIL FIXED"`
	ElectricityPrice         *float32                          `json:"electricityPrice" validate:"omitempty"`
	ElectricityCustomerCode  *string                           `json:"electricityCustomerCode" validate:"omitempty"`
	ElectricityProvider      *string                           `json:"electricityProvider" validate:"omitempty"`
	WaterSetupBy             string                            `json:"waterSetupBy" validate:"required,oneof=LANDLORD TENANT"`
	WaterPaymentType         *string                           `json:"waterPaymentType" validate:"omitempty,oneof=RETAIL FIXED"`
	WaterPrice               *float32                          `json:"waterPrice" validate:"omitempty"`
	WaterCustomerCode        *string                           `json:"waterCustomerCode" validate:"omitempty"`
	WaterProvider            *string                           `json:"waterProvider" validate:"omitempty"`

	Note *string `json:"note"`

	Coaps    []CreateRentalCoap    `json:"coaps"`
	Minors   []CreateRentalMinor   `json:"minors"`
	Pets     []CreateRentalPet     `json:"pets"`
	Services []CreateRentalService `json:"services"`
	Policies []CreateRentalPolicy  `json:"policies"`
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
		RentalPeriod: pm.RentalPeriod,
		PaymentType: database.NullRENTALPAYMENTTYPE{
			RENTALPAYMENTTYPE: pm.PaymentType,
			Valid:             pm.PaymentType != "",
		},
		RentalPrice:        pm.RentalPrice,
		RentalPaymentBasis: pm.RentalPaymentBasis,
		RentalIntention:    pm.RentalIntention,
		NoticePeriod:       types.Int32N(pm.NoticePeriod),
		GracePeriod: pgtype.Int4{
			Int32: pm.GracePeriod,
			Valid: pm.GracePeriod != 0,
		},
		LatePaymentPenaltyScheme: database.NullLATEPAYMENTPENALTYSCHEME{
			LATEPAYMENTPENALTYSCHEME: pm.LatePaymentPenaltyScheme,
			Valid:                    pm.LatePaymentPenaltyScheme != "",
		},
		LatePaymentPenaltyAmount: types.Float32N(pm.LatePaymentPenaltyAmount),
		ElectricitySetupBy:       pm.ElectricitySetupBy,
		ElectricityPaymentType:   types.StrN(pm.ElectricityPaymentType),
		ElectricityPrice:         types.Float32N(pm.ElectricityPrice),
		ElectricityCustomerCode:  types.StrN(pm.ElectricityCustomerCode),
		ElectricityProvider:      types.StrN(pm.ElectricityProvider),
		WaterSetupBy:             pm.WaterSetupBy,
		WaterPaymentType:         types.StrN(pm.WaterPaymentType),
		WaterPrice:               types.Float32N(pm.WaterPrice),
		WaterCustomerCode:        types.StrN(pm.WaterCustomerCode),
		WaterProvider:            types.StrN(pm.WaterProvider),
		Note:                     types.StrN(pm.Note),
	}
}

func (pm *CreateRental) Validate() error {
	if errs := validation.ValidateStruct(nil, *pm); len(errs) > 0 {
		return errs[0]
	}

	if pm.ElectricitySetupBy == "LANDLORD" && pm.ElectricityPaymentType == nil {
		return errors.New("electricityPaymentType is required")
	}

	if pm.WaterSetupBy == "LANDLORD" && pm.WaterPaymentType == nil {
		return errors.New("waterPaymentType is required")
	}

	return nil
}
