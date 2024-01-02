package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
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
	m := ApplicationMinorModel{
		ApplicationID: db.ApplicationID,
		FullName:      db.FullName,
		Dob:           db.Dob,
	}
	if db.Email.Valid {
		m.Email = &db.Email.String
	}
	if db.Phone.Valid {
		m.Phone = &db.Phone.String
	}
	if db.Description.Valid {
		m.Description = &db.Description.String
	}
	return m
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
	m := ApplicationCoapModel{
		ApplicationID: db.ApplicationID,
		FullName:      db.FullName,
		Dob:           db.Dob,
		Job:           db.Job,
		Income:        db.Income,
	}
	if db.Email.Valid {
		m.Email = &db.Email.String
	}
	if db.Phone.Valid {
		m.Phone = &db.Phone.String
	}
	if db.Description.Valid {
		m.Description = &db.Description.String
	}
	return m
}

/***/

type ApplicationPetModel struct {
	ApplicationID int64    `json:"applicationId" validate:"required"`
	Type          string   `json:"type" validate:"required"`
	Weight        *float64 `json:"weight" validate:"omitempty"`
	Description   *string  `json:"description" validate:"omitempty"`
}

func ToApplicationPetModel(db *database.ApplicationPet) ApplicationPetModel {
	m := ApplicationPetModel{
		ApplicationID: db.ApplicationID,
		Type:          db.Type,
	}
	if db.Weight.Valid {
		m.Weight = &db.Weight.Float64
	}
	if db.Description.Valid {
		m.Description = &db.Description.String
	}
	return m
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
	m := ApplicationVehicle{
		ApplicationID: db.ApplicationID,
		Type:          db.Type,
		Code:          db.Code,
	}
	if db.Model.Valid {
		m.Model = &db.Model.String
	}
	if db.Description.Valid {
		m.Description = &db.Description.String
	}
	return m
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
	am := &ApplicationModel{
		ID:            a.ID,
		PropertyID:    a.PropertyID,
		Status:        a.Status,
		CreatorID:     a.CreatorID,
		FullName:      a.FullName,
		Email:         a.Email,
		Phone:         a.Phone,
		Dob:           a.Dob,
		ProfileImage:  a.ProfileImage,
		MoveinDate:    a.MoveinDate,
		PreferredTerm: a.PreferredTerm,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
	if a.RhAddress.Valid {
		am.RhAddress = &a.RhAddress.String
	}
	if a.RhCity.Valid {
		am.RhCity = &a.RhCity.String
	}
	if a.RhDistrict.Valid {
		am.RhDistrict = &a.RhDistrict.String
	}
	if a.RhWard.Valid {
		am.RhWard = &a.RhWard.String
	}
	if a.RhRentalDuration.Valid {
		am.RhRentalDuration = &a.RhRentalDuration.Int32
	}
	if a.RhMonthlyPayment.Valid {
		am.RhMonthlyPayment = &a.RhMonthlyPayment.Float64
	}
	if a.RhReasonForLeaving.Valid {
		am.RhReasonForLeaving = &a.RhReasonForLeaving.String
	}
	if a.EmploymentCompanyName.Valid {
		am.EmploymentCompanyName = &a.EmploymentCompanyName.String
	}
	if a.EmploymentPosition.Valid {
		am.EmploymentPosition = &a.EmploymentPosition.String
	}
	if a.EmploymentMonthlyIncome.Valid {
		am.EmploymentMonthlyIncome = &a.EmploymentMonthlyIncome.Float64
	}
	if a.EmploymentComment.Valid {
		am.EmploymentComment = &a.EmploymentComment.String
	}
	return am
}
