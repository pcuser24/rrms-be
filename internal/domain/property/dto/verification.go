package dto

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Media struct {
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required,gt=0"`
	Type string `json:"type" validate:"required"`
	Url  string `json:"url"`
}

type PreCreatePropertyVerificationRequest struct {
	PropertyID                uuid.UUID `json:"propertyId" validate:"required,uuid4"`
	HouseOwnershipCertificate *Media    `json:"houseOwnershipCertificate" validate:"omitempty"`
	CertificateOfLanduseRight *Media    `json:"certificateOfLanduseRight" validate:"omitempty"`
	FrontIdcard               Media     `json:"frontIdcard" validate:"required"`
	BackIdcard                Media     `json:"backIdcard" validate:"required"`
}

type CreatePropertyVerificationRequest struct {
	CreatorID                 uuid.UUID `json:"creatorId" validate:"required,uuid4"`
	PropertyID                uuid.UUID `json:"propertyId" validate:"required,uuid4"`
	VideoUrl                  string    `json:"videoUrl" validate:"required,url"`
	HouseOwnershipCertificate *string   `json:"houseOwnershipCertificate" validate:"omitempty,url"`
	CertificateOfLanduseRight *string   `json:"certificateOfLanduseRight" validate:"omitempty,url"`
	FrontIdcard               string    `json:"frontIdcard" validate:"required,url"`
	BackIdcard                string    `json:"backIdcard" validate:"required,url"`
	Note                      *string   `json:"note" validate:"omitempty"`
}

func (c *CreatePropertyVerificationRequest) ToCreatePropertyVerificationRequestParams() database.CreatePropertyVerificationRequestParams {
	return database.CreatePropertyVerificationRequestParams{
		CreatorID:                 c.CreatorID,
		PropertyID:                c.PropertyID,
		VideoUrl:                  c.VideoUrl,
		HouseOwnershipCertificate: types.StrN(c.HouseOwnershipCertificate),
		CertificateOfLanduseRight: types.StrN(c.CertificateOfLanduseRight),
		FrontIdcard:               c.FrontIdcard,
		BackIdcard:                c.BackIdcard,
		Note:                      types.StrN(c.Note),
	}
}

// type UpdatePropertyVerificationRequestParams struct {
// 	ID                        int64                               `json:"id" validate:"required"`
// 	VideoUrl                  *string                             `json:"videoUrl" validate:"omitempty"`
// 	HouseOwnershipCertificate *string                             `json:"houseOwnershipCertificate" validate:"omitempty"`
// 	CertificateOfLanduseRight *string                             `json:"certificateOfLanduseRight" validate:"omitempty"`
// 	FrontIdcard               *string                             `json:"frontIdcard" validate:"omitempty"`
// 	BackIdcard                *string                             `json:"backIdcard" validate:"omitempty"`
// 	Note                      *string                             `json:"note" validate:"omitempty"`
// 	Status                    database.PROPERTYVERIFICATIONSTATUS `json:"status" validate:"omitempty,oneof=APPROVED PENDING REJECTED"`
// }

// func (u *UpdatePropertyVerificationRequestParams) ToUpdatePropertyVerificationRequestParams() database.UpdatePropertyVerificationRequestParams {
// 	return database.UpdatePropertyVerificationRequestParams{
// 		ID:                        u.ID,
// 		VideoUrl:                  types.StrN(u.VideoUrl),
// 		HouseOwnershipCertificate: types.StrN(u.HouseOwnershipCertificate),
// 		CertificateOfLanduseRight: types.StrN(u.CertificateOfLanduseRight),
// 		FrontIdcard:               types.StrN(u.FrontIdcard),
// 		BackIdcard:                types.StrN(u.BackIdcard),
// 		Note:                      types.StrN(u.Note),
// 		Status: database.NullPROPERTYVERIFICATIONSTATUS{
// 			PROPERTYVERIFICATIONSTATUS: u.Status,
// 			Valid:                      u.Status != "",
// 		},
// 	}
// }

type UpdatePropertyVerificationRequestStatus struct {
	Status   database.PROPERTYVERIFICATIONSTATUS `json:"status" validate:"required,oneof=APPROVED PENDING REJECTED"`
	Feedback *string                             `json:"feedback" validate:"omitempty"`
}

type GetPropertyVerificationRequestsQuery struct {
	Limit  int32  `query:"limit" validate:"required,min=1,max=100"`
	Offset int32  `query:"offset" validate:"omitempty,min=0"`
	SortBy string `query:"orderBy" validate:"omitempty,oneof=created_at updated_at"`
	Order  string `query:"order" validate:"omitempty,oneof=asc desc"`

	// other filters
	CreatorID  uuid.UUID                             `query:"creatorId" validate:"omitempty"`
	PropertyID uuid.UUID                             `query:"propertyId" validate:"omitempty"`
	Status     []database.PROPERTYVERIFICATIONSTATUS `query:"status" validate:"omitempty,dive,oneof=APPROVED PENDING REJECTED"`
}

type GetPropertyVerificationRequestsResponse struct {
	FullCount int64                               `json:"fullCount"`
	Items     []model.PropertyVerificationRequest `json:"items"`
}
