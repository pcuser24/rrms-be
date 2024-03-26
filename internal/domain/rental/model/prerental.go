package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PrerentalCoapModel struct {
	PrerentalID int64     `json:"prerentalId"`
	FullName    *string   `json:"fullName"`
	Dob         time.Time `json:"dob"`
	Job         *string   `json:"job"`
	Income      *int32    `json:"income"`
	Email       *string   `json:"email"`
	Phone       *string   `json:"phone"`
	Description *string   `json:"description"`
}

func ToPreRentalCoapModel(pr *database.PrerentalCoap) *PrerentalCoapModel {
	return &PrerentalCoapModel{
		PrerentalID: pr.PrerentalID,
		FullName:    types.PNStr(pr.FullName),
		Dob:         pr.Dob.Time,
		Job:         types.PNStr(pr.Job),
		Income:      types.PNInt32(pr.Income),
		Email:       types.PNStr(pr.Email),
		Phone:       types.PNStr(pr.Phone),
		Description: types.PNStr(pr.Description),
	}
}

type PrerentalModel struct {
	ID                   int64                    `json:"id"`
	ApplicationID        *int64                   `json:"applicationId"`
	CreatorID            uuid.UUID                `json:"creatorId"`
	TenantID             uuid.UUID                `json:"tenantId"`
	ProfileImage         string                   `json:"profileImage"`
	PropertyID           uuid.UUID                `json:"propertyId"`
	UnitID               uuid.UUID                `json:"unitId"`
	TenantType           database.TENANTTYPE      `json:"tenantType"`
	TenantIdentity       string                   `json:"tenantIdentity"`
	TenantDob            time.Time                `json:"tenantDob"`
	TenantPhone          string                   `json:"tenantPhone"`
	TenantEmail          string                   `json:"tenantEmail"`
	TenantAddress        *string                  `json:"tenantAddress"`
	ContractType         database.CONTRACTTYPE    `json:"contractType"`
	ContractContent      *string                  `json:"contractContent"`
	ContractLastUpdateAt time.Time                `json:"contractLastUpdate"`
	ContractLastUpdateBy uuid.UUID                `json:"contractLastUpdateBy"`
	LandArea             float32                  `json:"landArea"`
	PropertyArea         float32                  `json:"propertyArea"`
	StartDate            time.Time                `json:"startDate"`
	MoveinDate           time.Time                `json:"moveinDate"`
	RentalPeriod         int32                    `json:"rentalPeriod"`
	RentalPrice          float32                  `json:"rentalPrice"`
	Status               database.PRERENTALSTATUS `json:"status"`
	Note                 *string                  `json:"note"`

	Coaps []PrerentalCoapModel `json:"coaps"`
}

type PreRentalContractModel struct {
	ID                   int64                 `json:"id"`
	ContractType         database.CONTRACTTYPE `json:"contractType"`
	ContractContent      *string               `json:"contractContent"`
	ContractLastUpdateAt time.Time             `json:"contractLastUpdate"`
	ContractLastUpdateBy uuid.UUID             `json:"contractLastUpdateBy"`
}

func ToPreRentalModel(pr *database.Prerental) *PrerentalModel {
	return &PrerentalModel{
		ID:                   pr.ID,
		ApplicationID:        types.PNInt64(pr.ApplicationID),
		CreatorID:            pr.CreatorID,
		TenantID:             pr.TenantID.Bytes,
		ProfileImage:         pr.ProfileImage,
		PropertyID:           pr.PropertyID,
		UnitID:               pr.UnitID,
		TenantType:           pr.TenantType,
		TenantIdentity:       pr.TenantIdentity,
		TenantDob:            pr.TenantDob.Time,
		TenantPhone:          pr.TenantPhone,
		TenantEmail:          pr.TenantEmail,
		TenantAddress:        types.PNStr(pr.TenantAddress),
		ContractType:         pr.ContractType.CONTRACTTYPE,
		ContractContent:      types.PNStr(pr.ContractContent),
		ContractLastUpdateAt: pr.ContractLastUpdateAt.Time,
		ContractLastUpdateBy: pr.ContractLastUpdateBy.Bytes,
		LandArea:             pr.LandArea,
		PropertyArea:         pr.UnitArea,
		StartDate:            pr.StartDate.Time,
		MoveinDate:           pr.MoveinDate.Time,
		RentalPeriod:         pr.RentalPeriod,
		RentalPrice:          pr.RentalPrice,
		Status:               pr.Status,
		Note:                 types.PNStr(pr.Note),
	}
}
