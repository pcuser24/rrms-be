package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type CreateApplicationMinorModel struct {
	FullName    string    `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"required"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationMinorModel) ToCreateApplicationMinorDB(aid int64) *database.CreateApplicationMinorParams {
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

type CreateApplicationPetModel struct {
	Type        string   `json:"type" validate:"required"`
	Weight      *float32 `json:"weight" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationPetModel) ToCreateApplicationPetDB(aid int64) *database.CreateApplicationPetParams {
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

type CreateApplicationDto struct {
	ListingID                uuid.UUID   `json:"listingId" validate:"required,uuid4"`
	PropertyID               uuid.UUID   `json:"propertyId" validate:"required,uuid4"`
	UnitIds                  []uuid.UUID `json:"unitIds" validate:"required,dive,uuid4"`
	CreatorID                uuid.UUID   `json:"creatorId" validate:"required,uuid4"`
	FullName                 string      `json:"fullName" validate:"required"`
	Email                    string      `json:"email" validate:"required,email"`
	Phone                    string      `json:"phone" validate:"required"`
	Dob                      time.Time   `json:"dob" validate:"required"`
	ProfileImage             string      `json:"profileImage" validate:"required,url"`
	MoveinDate               time.Time   `json:"moveinDate" validate:"required"`
	PreferredTerm            int32       `json:"preferredTerm" validate:"required,gt=0"`
	RhAddress                *string     `json:"rhAddress" validate:"omitempty"`
	RhCity                   *string     `json:"rhCity" validate:"omitempty"`
	RhDistrict               *string     `json:"rhDistrict" validate:"omitempty"`
	RhWard                   *string     `json:"rhWard" validate:"omitempty"`
	RhRentalDuration         *int32      `json:"rhRentalDuration" validate:"omitempty,gt=0"`
	RhMonthlyPayment         *float32    `json:"rhMonthlyPayment" validate:"omitempty,gt=0"`
	RhReasonForLeaving       *string     `json:"rhReasonForLeaving" validate:"omitempty"`
	EmploymentStatus         string      `json:"employmentStatus" validate:"required,oneof=UNEMPLOYED EMPLOYED SELF-EMPLOYED RETIRED STUDENT"`
	EmploymentCompanyName    *string     `json:"employmentCompanyName" validate:"omitempty"`
	EmploymentPosition       *string     `json:"employmentPosition" validate:"omitempty"`
	EmploymentMonthlyIncome  *float32    `json:"employmentMonthlyIncome" validate:"omitempty,gt=0"`
	EmploymentComment        *string     `json:"employmentComment" validate:"omitempty"`
	EmploymentProofsOfIncome []string    `json:"employmentProofsOfIncome" validate:"omitempty"`
	IdentityType             string      `json:"identityType" validate:"required,oneof=ID CITIZENIDENTIFICATION PASSPORT DRIVERLICENSE"`
	IdentityNumber           string      `json:"identityNumber" validate:"required"`
	IdentityIssuedDate       time.Time   `json:"identityIssuedDate" validate:"required"`
	IdentityIssuedBy         string      `json:"identityIssuedBy" validate:"required"`

	Minors   []CreateApplicationMinorModel `json:"minors" validate:"dive"`
	Coaps    []CreateApplicationCoapModel  `json:"coaps" validate:"dive"`
	Pets     []CreateApplicationPetModel   `json:"pets" validate:"dive"`
	Vehicles []CreateApplicationVehicle    `json:"vehicles" validate:"dive"`
}

func (a *CreateApplicationDto) ToCreateApplicationDB() *database.CreateApplicationParams {
	return &database.CreateApplicationParams{
		PropertyID:               a.PropertyID,
		UnitIds:                  a.UnitIds,
		CreatorID:                a.CreatorID,
		FullName:                 a.FullName,
		Email:                    a.Email,
		Phone:                    a.Phone,
		Dob:                      a.Dob,
		ProfileImage:             a.ProfileImage,
		MoveinDate:               a.MoveinDate,
		PreferredTerm:            a.PreferredTerm,
		EmploymentStatus:         a.EmploymentStatus,
		EmploymentCompanyName:    types.StrN(a.EmploymentCompanyName),
		EmploymentPosition:       types.StrN(a.EmploymentPosition),
		EmploymentMonthlyIncome:  types.Float32N(a.EmploymentMonthlyIncome),
		EmploymentComment:        types.StrN(a.EmploymentComment),
		EmploymentProofsOfIncome: a.EmploymentProofsOfIncome,
		RhAddress:                types.StrN(a.RhAddress),
		RhCity:                   types.StrN(a.RhCity),
		RhDistrict:               types.StrN(a.RhDistrict),
		RhWard:                   types.StrN(a.RhWard),
		RhRentalDuration:         types.Int32N(a.RhRentalDuration),
		RhMonthlyPayment:         types.Float32N(a.RhMonthlyPayment),
		RhReasonForLeaving:       types.StrN(a.RhReasonForLeaving),
		IdentityType:             a.IdentityType,
		IdentityNumber:           a.IdentityNumber,
		IdentityIssuedDate:       a.IdentityIssuedDate,
		IdentityIssuedBy:         a.IdentityIssuedBy,
	}
}
