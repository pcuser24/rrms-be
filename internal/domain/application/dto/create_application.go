package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateApplicationMinor struct {
	FullName    string    `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"required"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationMinor) ToCreateApplicationMinorDB(aid int64) *database.CreateApplicationMinorParams {
	return &database.CreateApplicationMinorParams{
		ApplicationID: aid,
		FullName:      m.FullName,
		Dob:           m.Dob,
		Email:         types.StrN(m.Email),
		Phone:         types.StrN(m.Phone),
		Description:   types.StrN(m.Description),
	}
}

type CreateApplicationCoapModel struct {
	FullName    string    `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"required"`
	Job         string    `json:"job" validate:"required"`
	Income      int32     `json:"income" validate:"required"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationCoapModel) ToCreateApplicationCoapDB(aid int64) *database.CreateApplicationCoapParams {
	return &database.CreateApplicationCoapParams{
		ApplicationID: aid,
		FullName:      m.FullName,
		Dob:           m.Dob,
		Job:           m.Job,
		Income:        m.Income,
		Email:         types.StrN(m.Email),
		Phone:         types.StrN(m.Phone),
		Description:   types.StrN(m.Description),
	}
}

type CreateApplicationPet struct {
	Type        string   `json:"type" validate:"required"`
	Weight      *float32 `json:"weight" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationPet) ToCreateApplicationPetDB(aid int64) *database.CreateApplicationPetParams {
	return &database.CreateApplicationPetParams{
		ApplicationID: aid,
		Type:          m.Type,
		Weight:        types.Float32N(m.Weight),
		Description:   types.StrN(m.Description),
	}
}

type CreateApplicationVehicle struct {
	Type        string  `json:"type" validate:"required"`
	Model       *string `json:"model" validate:"omitempty"`
	Code        string  `json:"code" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationVehicle) ToCreateApplicationVehicleDB(aid int64) *database.CreateApplicationVehicleParams {
	return &database.CreateApplicationVehicleParams{
		ApplicationID: aid,
		Type:          m.Type,
		Code:          m.Code,
		Model:         types.StrN(m.Model),
		Description:   types.StrN(m.Description),
	}
}

type CreateApplication struct {
	ListingID               uuid.UUID           `json:"listingId" validate:"omitempty,uuid4"`
	PropertyID              uuid.UUID           `json:"propertyId" validate:"required,uuid4"`
	UnitID                  uuid.UUID           `json:"unitId" validate:"required,uuid4"`
	ListingPrice            int64               `json:"listingPrice" validate:"required,gt=0"`
	OfferedPrice            int64               `json:"offeredPrice" validate:"required,gt=0"`
	CreatorID               uuid.UUID           `json:"creatorId"`
	TenantType              database.TENANTTYPE `json:"tenantType" validate:"required,oneof=INDIVIDUAL FAMILY ORGANIZATION"`
	FullName                string              `json:"fullName" validate:"required"`
	Email                   string              `json:"email" validate:"required,email"`
	Phone                   string              `json:"phone" validate:"required"`
	Dob                     time.Time           `json:"dob" validate:"omitempty"`
	ProfileImage            string              `json:"profileImage" validate:"required,url"`
	MoveinDate              time.Time           `json:"moveinDate" validate:"required"`
	PreferredTerm           int32               `json:"preferredTerm" validate:"required,gt=0"`
	RentalIntention         string              `json:"rentalIntention" validate:"required"`
	OrganizationName        *string             `json:"organizationName" validate:"omitempty"`
	OrganizationHqAddress   *string             `json:"organizationHqAddress" validate:"omitempty"`
	OrganizationScale       *string             `json:"organizationScale" validate:"omitempty"`
	RhAddress               *string             `json:"rhAddress" validate:"omitempty"`
	RhCity                  *string             `json:"rhCity" validate:"omitempty"`
	RhDistrict              *string             `json:"rhDistrict" validate:"omitempty"`
	RhWard                  *string             `json:"rhWard" validate:"omitempty"`
	RhRentalDuration        *int32              `json:"rhRentalDuration" validate:"omitempty,gt=0"`
	RhMonthlyPayment        *int64              `json:"rhMonthlyPayment" validate:"omitempty,gt=0"`
	RhReasonForLeaving      *string             `json:"rhReasonForLeaving" validate:"omitempty"`
	EmploymentStatus        string              `json:"employmentStatus" validate:"omitempty,oneof=UNEMPLOYED EMPLOYED SELF-EMPLOYED RETIRED STUDENT"`
	EmploymentCompanyName   *string             `json:"employmentCompanyName" validate:"omitempty"`
	EmploymentPosition      *string             `json:"employmentPosition" validate:"omitempty"`
	EmploymentMonthlyIncome *int64              `json:"employmentMonthlyIncome" validate:"omitempty,gt=0"`
	EmploymentComment       *string             `json:"employmentComment" validate:"omitempty"`

	Minors   []CreateApplicationMinor     `json:"minors" validate:"dive"`
	Coaps    []CreateApplicationCoapModel `json:"coaps" validate:"dive"`
	Pets     []CreateApplicationPet       `json:"pets" validate:"dive"`
	Vehicles []CreateApplicationVehicle   `json:"vehicles" validate:"dive"`

	ApplicationKey string `json:"k" validate:"omitempty"`
}

func (a *CreateApplication) ToCreateApplicationDB() *database.CreateApplicationParams {
	return &database.CreateApplicationParams{
		ListingID:    a.ListingID,
		PropertyID:   a.PropertyID,
		UnitID:       a.UnitID,
		ListingPrice: a.ListingPrice,
		OfferedPrice: a.OfferedPrice,
		CreatorID:    types.UUIDN(a.CreatorID),
		TenantType:   a.TenantType,
		FullName:     a.FullName,
		Email:        a.Email,
		Phone:        a.Phone,
		Dob: pgtype.Date{
			Time:  a.Dob,
			Valid: !a.Dob.IsZero(),
		},
		ProfileImage: a.ProfileImage,
		MoveinDate: pgtype.Date{
			Time:  a.MoveinDate,
			Valid: !a.MoveinDate.IsZero(),
		},
		PreferredTerm:           a.PreferredTerm,
		RentalIntention:         a.RentalIntention,
		EmploymentStatus:        a.EmploymentStatus,
		EmploymentCompanyName:   types.StrN(a.EmploymentCompanyName),
		EmploymentPosition:      types.StrN(a.EmploymentPosition),
		EmploymentMonthlyIncome: types.Int64N(a.EmploymentMonthlyIncome),
		EmploymentComment:       types.StrN(a.EmploymentComment),
		OrganizationName:        types.StrN(a.OrganizationName),
		OrganizationHqAddress:   types.StrN(a.OrganizationHqAddress),
		OrganizationScale:       types.StrN(a.OrganizationScale),
		RhAddress:               types.StrN(a.RhAddress),
		RhCity:                  types.StrN(a.RhCity),
		RhDistrict:              types.StrN(a.RhDistrict),
		RhWard:                  types.StrN(a.RhWard),
		RhRentalDuration:        types.Int32N(a.RhRentalDuration),
		RhMonthlyPayment:        types.Int64N(a.RhMonthlyPayment),
		RhReasonForLeaving:      types.StrN(a.RhReasonForLeaving),
	}
}
