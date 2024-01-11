package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

/***/

type ApplicationMinorModel struct {
	ApplicationID int64     `json:"applicationId" validate:"required"`
	FullName      string    `json:"fullName" validate:"required"`
	Dob           time.Time `json:"dob" validate:"required"`
	Email         *string   `json:"email" validate:"omitempty,email"`
	Phone         *string   `json:"phone" validate:"omitempty"`
	Description   *string   `json:"description" validate:"omitempty"`
}

func ToApplicationMinorModel(db *database.ApplicationMinor) ApplicationMinorModel {
	return ApplicationMinorModel{
		ApplicationID: db.ApplicationID,
		FullName:      db.FullName,
		Email:         types.PNStr(db.Email),
		Phone:         types.PNStr(db.Phone),
		Description:   types.PNStr(db.Description),
		Dob:           db.Dob,
	}
}

/***/

type ApplicationCoapModel struct {
	ApplicationID int64     `json:"applicationId" validate:"required"`
	FullName      string    `json:"fullName" validate:"required"`
	Dob           time.Time `json:"dob" validate:"required"`
	Job           string    `json:"job" validate:"required"`
	Income        int32     `json:"income" validate:"required"`
	Email         *string   `json:"email" validate:"omitempty,email"`
	Phone         *string   `json:"phone" validate:"omitempty"`
	Description   *string   `json:"description" validate:"omitempty"`
}

func ToApplicationCoapModel(db *database.ApplicationCoap) ApplicationCoapModel {
	return ApplicationCoapModel{
		ApplicationID: db.ApplicationID,
		FullName:      db.FullName,
		Dob:           db.Dob,
		Job:           db.Job,
		Email:         types.PNStr(db.Email),
		Phone:         types.PNStr(db.Phone),
		Description:   types.PNStr(db.Description),
		Income:        db.Income,
	}
}

/***/

type ApplicationPetModel struct {
	ApplicationID int64    `json:"applicationId" validate:"required"`
	Type          string   `json:"type" validate:"required"`
	Weight        *float64 `json:"weight" validate:"omitempty"`
	Description   *string  `json:"description" validate:"omitempty"`
}

func ToApplicationPetModel(db *database.ApplicationPet) ApplicationPetModel {
	return ApplicationPetModel{
		ApplicationID: db.ApplicationID,
		Weight:        types.PNFloat64(db.Weight),
		Description:   types.PNStr(db.Description),
		Type:          db.Type,
	}
}

/***/

type ApplicationVehicle struct {
	ApplicationID int64   `json:"applicationId" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Model         *string `json:"model" validate:"omitempty"`
	Code          string  `json:"code" validate:"required"`
	Description   *string `json:"description" validate:"omitempty"`
}

func ToApplicationVehicleModel(db *database.ApplicationVehicle) ApplicationVehicle {
	return ApplicationVehicle{
		ApplicationID: db.ApplicationID,
		Type:          db.Type,
		Model:         types.PNStr(db.Model),
		Code:          db.Code,
		Description:   types.PNStr(db.Description),
	}
}

type ApplicationModel struct {
	ID                       int64                      `json:"id"`
	ListingID                uuid.UUID                  `json:"listingId"`
	PropertyID               uuid.UUID                  `json:"propertyId"`
	UnitIds                  []uuid.UUID                `json:"unitIds"`
	Status                   database.APPLICATIONSTATUS `json:"status"`
	CreatorID                uuid.UUID                  `json:"creatorId"`
	FullName                 string                     `json:"fullName"`
	Email                    string                     `json:"email"`
	Phone                    string                     `json:"phone"`
	Dob                      time.Time                  `json:"dob"`
	ProfileImage             string                     `json:"profileImage"`
	MoveinDate               time.Time                  `json:"moveinDate"`
	PreferredTerm            int32                      `json:"preferredTerm"`
	RhAddress                *string                    `json:"rhAddress"`
	RhCity                   *string                    `json:"rhCity"`
	RhDistrict               *string                    `json:"rhDistrict"`
	RhWard                   *string                    `json:"rhWard"`
	RhRentalDuration         *int32                     `json:"rhRentalDuration"`
	RhMonthlyPayment         *float64                   `json:"rhMonthlyPayment"`
	RhReasonForLeaving       *string                    `json:"rhReasonForLeaving"`
	EmploymentStatus         string                     `json:"employmentStatus"`
	EmploymentCompanyName    *string                    `json:"employmentCompanyName"`
	EmploymentPosition       *string                    `json:"employmentPosition"`
	EmploymentMonthlyIncome  *float64                   `json:"employmentMonthlyIncome"`
	EmploymentComment        *string                    `json:"employmentComment"`
	EmploymentProofsOfIncome []string                   `json:"employmentProofsOfIncome"`
	IdentityType             string                     `json:"identityType"`
	IdentityNumber           string                     `json:"identityNumber"`
	IdentityIssuedDate       time.Time                  `json:"identityIssuedDate"`
	IdentityIssuedBy         string                     `json:"identityIssuedBy"`
	CreatedAt                time.Time                  `json:"createdAt"`
	UpdatedAt                time.Time                  `json:"updatedAt"`

	Minors   []ApplicationMinorModel `json:"minors"`
	Coaps    []ApplicationCoapModel  `json:"coaps"`
	Pets     []ApplicationPetModel   `json:"pets"`
	Vehicles []ApplicationVehicle    `json:"vehicles"`
}

func ToApplicationModel(a *database.Application) *ApplicationModel {
	return &ApplicationModel{
		ID:                       a.ID,
		CreatorID:                a.CreatorID,
		ListingID:                a.ListingID,
		PropertyID:               a.PropertyID,
		UnitIds:                  a.UnitIds,
		Status:                   a.Status,
		FullName:                 a.FullName,
		Email:                    a.Email,
		Phone:                    a.Phone,
		Dob:                      a.Dob,
		ProfileImage:             a.ProfileImage,
		MoveinDate:               a.MoveinDate,
		PreferredTerm:            a.PreferredTerm,
		RhAddress:                types.PNStr(a.RhAddress),
		RhCity:                   types.PNStr(a.RhCity),
		RhDistrict:               types.PNStr(a.RhDistrict),
		RhWard:                   types.PNStr(a.RhWard),
		RhRentalDuration:         types.PNInt32(a.RhRentalDuration),
		RhMonthlyPayment:         types.PNFloat64(a.RhMonthlyPayment),
		RhReasonForLeaving:       types.PNStr(a.RhReasonForLeaving),
		EmploymentStatus:         a.EmploymentStatus,
		EmploymentCompanyName:    types.PNStr(a.EmploymentCompanyName),
		EmploymentPosition:       types.PNStr(a.EmploymentPosition),
		EmploymentMonthlyIncome:  types.PNFloat64(a.EmploymentMonthlyIncome),
		EmploymentComment:        types.PNStr(a.EmploymentComment),
		EmploymentProofsOfIncome: a.EmploymentProofsOfIncome,
		IdentityType:             a.IdentityType,
		IdentityNumber:           a.IdentityNumber,
		IdentityIssuedDate:       a.IdentityIssuedDate,
		IdentityIssuedBy:         a.IdentityIssuedBy,
		CreatedAt:                a.CreatedAt,
		UpdatedAt:                a.UpdatedAt,
	}
}
