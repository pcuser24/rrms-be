package dto

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type CreateApplicationMinorModel struct {
	FullName    string    `json:"fullName" validate:"required"`
	Dob         time.Time `json:"dob" validate:"required"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Phone       *string   `json:"phone" validate:"omitempty"`
	Description *string   `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationMinorModel) ToCreateApplicationMinorDB(aid int64) *database.CreateApplicationMinorParams {
	db := &database.CreateApplicationMinorParams{
		ApplicationID: aid,
		FullName:      m.FullName,
		Dob:           m.Dob,
	}
	if m.Email != nil {
		db.Email = sql.NullString{
			String: *m.Email,
			Valid:  true,
		}
	}
	if m.Phone != nil {
		db.Phone = sql.NullString{
			String: *m.Phone,
			Valid:  true,
		}
	}
	if m.Description != nil {
		db.Description = sql.NullString{
			String: *m.Description,
			Valid:  true,
		}
	}

	return db
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
	db := &database.CreateApplicationCoapParams{
		ApplicationID: aid,
		FullName:      m.FullName,
		Dob:           m.Dob,
		Job:           m.Job,
		Income:        m.Income,
	}
	if m.Email != nil {
		db.Email = sql.NullString{
			String: *m.Email,
			Valid:  true,
		}
	}
	if m.Phone != nil {
		db.Phone = sql.NullString{
			String: *m.Phone,
			Valid:  true,
		}
	}
	if m.Description != nil {
		db.Description = sql.NullString{
			String: *m.Description,
			Valid:  true,
		}
	}

	return db
}

type CreateApplicationPetModel struct {
	Type        string   `json:"type" validate:"required"`
	Weight      *float64 `json:"weight" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationPetModel) ToCreateApplicationPetDB(aid int64) *database.CreateApplicationPetParams {
	db := &database.CreateApplicationPetParams{
		ApplicationID: aid,
		Type:          m.Type,
	}
	if m.Weight != nil {
		db.Weight = sql.NullFloat64{
			Float64: *m.Weight,
			Valid:   true,
		}
	}
	if m.Description != nil {
		db.Description = sql.NullString{
			String: *m.Description,
			Valid:  true,
		}
	}

	return db
}

type CreateApplicationVehicle struct {
	Type        string  `json:"type" validate:"required"`
	Model       *string `json:"model" validate:"omitempty"`
	Code        string  `json:"code" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}

func (m *CreateApplicationVehicle) ToCreateApplicationVehicleDB(aid int64) *database.CreateApplicationVehicleParams {
	db := &database.CreateApplicationVehicleParams{
		ApplicationID: aid,
		Type:          m.Type,
		Code:          m.Code,
	}
	if m.Model != nil {
		db.Model = sql.NullString{
			String: *m.Model,
			Valid:  true,
		}
	}
	if m.Description != nil {
		db.Description = sql.NullString{
			String: *m.Description,
			Valid:  true,
		}
	}

	return db
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
	RhMonthlyPayment         *float64    `json:"rhMonthlyPayment" validate:"omitempty,gt=0"`
	RhReasonForLeaving       *string     `json:"rhReasonForLeaving" validate:"omitempty"`
	EmploymentStatus         string      `json:"employmentStatus" validate:"required,oneof=UNEMPLOYED EMPLOYED SELF-EMPLOYED RETIRED STUDENT"`
	EmploymentCompanyName    *string     `json:"employmentCompanyName" validate:"omitempty"`
	EmploymentPosition       *string     `json:"employmentPosition" validate:"omitempty"`
	EmploymentMonthlyIncome  *float64    `json:"employmentMonthlyIncome" validate:"omitempty,gt=0"`
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
	adb := &database.CreateApplicationParams{
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
		EmploymentProofsOfIncome: a.EmploymentProofsOfIncome,
		IdentityType:             a.IdentityType,
		IdentityNumber:           a.IdentityNumber,
		IdentityIssuedDate:       a.IdentityIssuedDate,
		IdentityIssuedBy:         a.IdentityIssuedBy,
	}
	if a.RhAddress != nil {
		adb.RhAddress = sql.NullString{
			String: *a.RhAddress,
			Valid:  true,
		}
	}
	if a.RhCity != nil {
		adb.RhCity = sql.NullString{
			String: *a.RhCity,
			Valid:  true,
		}
	}
	if a.RhDistrict != nil {
		adb.RhDistrict = sql.NullString{
			String: *a.RhDistrict,
			Valid:  true,
		}
	}
	if a.RhWard != nil {
		adb.RhWard = sql.NullString{
			String: *a.RhWard,
			Valid:  true,
		}
	}
	if a.RhRentalDuration != nil {
		adb.RhRentalDuration = sql.NullInt32{
			Int32: *a.RhRentalDuration,
			Valid: true,
		}
	}
	if a.RhMonthlyPayment != nil {
		adb.RhMonthlyPayment = sql.NullFloat64{
			Float64: *a.RhMonthlyPayment,
			Valid:   true,
		}
	}
	if a.RhReasonForLeaving != nil {
		adb.RhReasonForLeaving = sql.NullString{
			String: *a.RhReasonForLeaving,
			Valid:  true,
		}
	}
	if a.EmploymentCompanyName != nil {
		adb.EmploymentCompanyName = sql.NullString{
			String: *a.EmploymentCompanyName,
			Valid:  true,
		}
	}
	if a.EmploymentPosition != nil {
		adb.EmploymentPosition = sql.NullString{
			String: *a.EmploymentPosition,
			Valid:  true,
		}
	}
	if a.EmploymentMonthlyIncome != nil {
		adb.EmploymentMonthlyIncome = sql.NullFloat64{
			Float64: *a.EmploymentMonthlyIncome,
			Valid:   true,
		}
	}
	if a.EmploymentComment != nil {
		adb.EmploymentComment = sql.NullString{
			String: *a.EmploymentComment,
			Valid:  true,
		}
	}
	return adb
}
