package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type RentalCoapModel struct {
	RentalID    int64     `json:"rentalId"`
	FullName    *string   `json:"fullName"`
	Dob         time.Time `json:"dob"`
	Job         *string   `json:"job"`
	Income      *int32    `json:"income"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	Description *string   `json:"description"`
}

func ToRentalCoapModel(pr *database.RentalCoap) RentalCoapModel {
	return RentalCoapModel{
		RentalID:    pr.RentalID,
		FullName:    types.PNStr(pr.FullName),
		Dob:         pr.Dob.Time,
		Job:         types.PNStr(pr.Job),
		Income:      types.PNInt32(pr.Income),
		Email:       types.PNStr(pr.Email),
		Phone:       types.PNStr(pr.Phone),
		Description: types.PNStr(pr.Description),
	}
}

type RentalMinor struct {
	RentalID    int64     `json:"rentalId"`
	FullName    string    `json:"fullName"`
	Dob         time.Time `json:"dob"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	Description *string   `json:"description"`
}

func ToRentalMinor(pr *database.RentalMinor) RentalMinor {
	return RentalMinor{
		RentalID:    pr.RentalID,
		FullName:    pr.FullName,
		Dob:         pr.Dob.Time,
		Email:       types.PNStr(pr.Email),
		Phone:       types.PNStr(pr.Phone),
		Description: types.PNStr(pr.Description),
	}
}

type RentalPet struct {
	RentalID    int64    `json:"rental_id"`
	Type        string   `json:"type"`
	Weight      *float32 `json:"weight"`
	Description *string  `json:"description"`
}

func ToRentalPet(pr *database.RentalPet) RentalPet {
	return RentalPet{
		RentalID:    pr.RentalID,
		Type:        pr.Type,
		Weight:      types.PNFloat32(pr.Weight),
		Description: types.PNStr(pr.Description),
	}
}

type RentalService struct {
	RentalID int64  `json:"rental_id"`
	Name     string `json:"name"`
	// The party who set up the service, either "LANDLORD" or "TENANT"
	SetupBy  string   `json:"setupBy"`
	Provider *string  `json:"provider"`
	Price    *float32 `json:"price"`
}

func ToRentalService(pr *database.RentalService) RentalService {
	return RentalService{
		RentalID: pr.RentalID,
		Name:     pr.Name,
		SetupBy:  pr.SetupBy,
		Provider: types.PNStr(pr.Provider),
		Price:    types.PNFloat32(pr.Price),
	}
}

type RentalModel struct {
	ID                     int64               `json:"id"`
	CreatorID              uuid.UUID           `json:"creatorId"`
	PropertyID             uuid.UUID           `json:"propertyId"`
	UnitID                 uuid.UUID           `json:"unitId"`
	ApplicationID          *int64              `json:"applicationId"`
	ProfileImage           string              `json:"profileImage"`
	TenantID               uuid.UUID           `json:"tenantId"`
	TenantType             database.TENANTTYPE `json:"tenantType"`
	TenantName             string              `json:"tenantName"`
	TenantPhone            string              `json:"tenantPhone"`
	TenantEmail            string              `json:"tenantEmail"`
	OrganizationName       *string             `json:"organizationName" validate:"omitempty"`
	OrganizationHqAddress  *string             `json:"organizationHqAddress" validate:"omitempty"`
	StartDate              time.Time           `json:"startDate"`
	MoveinDate             time.Time           `json:"moveinDate"`
	RentalPeriod           int32               `json:"rentalPeriod"`
	RentalPrice            float32             `json:"rentalPrice"`
	RentalIntention        string              `json:"rentalIntention"`
	Deposit                float32             `json:"deposit"`
	DepositPaid            bool                `json:"depositPaid"`
	ElectricityPaymentType string              `json:"electricityPaymentType"`
	ElectricityPrice       *float32            `json:"electricityPrice"`
	WaterPaymentType       string              `json:"waterPaymentType"`
	WaterPrice             *float32            `json:"waterPrice"`
	Note                   *string             `json:"note"`
	CreatedAt              time.Time           `json:"createdAt"`
	UpdatedAt              time.Time           `json:"updatedAt"`

	Coaps    []RentalCoapModel `json:"coaps"`
	Minors   []RentalMinor     `json:"minors"`
	Pets     []RentalPet       `json:"pets"`
	Services []RentalService   `json:"services"`
}

func ToRentalModel(pr *database.Rental) *RentalModel {
	return &RentalModel{
		ID:                     pr.ID,
		CreatorID:              pr.CreatorID,
		PropertyID:             pr.PropertyID,
		UnitID:                 pr.UnitID,
		ApplicationID:          types.PNInt64(pr.ApplicationID),
		ProfileImage:           pr.ProfileImage,
		TenantID:               pr.TenantID.Bytes,
		TenantType:             pr.TenantType,
		TenantName:             pr.TenantName,
		TenantPhone:            pr.TenantPhone,
		TenantEmail:            pr.TenantEmail,
		OrganizationName:       types.PNStr(pr.OrganizationName),
		OrganizationHqAddress:  types.PNStr(pr.OrganizationHqAddress),
		StartDate:              pr.StartDate.Time,
		MoveinDate:             pr.MoveinDate.Time,
		RentalPeriod:           pr.RentalPeriod,
		RentalPrice:            pr.RentalPrice,
		RentalIntention:        pr.RentalIntention,
		Deposit:                pr.Deposit,
		DepositPaid:            pr.DepositPaid,
		ElectricityPaymentType: pr.ElectricityPaymentType,
		ElectricityPrice:       types.PNFloat32(pr.ElectricityPrice),
		WaterPaymentType:       pr.WaterPaymentType,
		WaterPrice:             types.PNFloat32(pr.WaterPrice),
		Note:                   types.PNStr(pr.Note),
		CreatedAt:              pr.CreatedAt,
		UpdatedAt:              pr.UpdatedAt,
	}
}
