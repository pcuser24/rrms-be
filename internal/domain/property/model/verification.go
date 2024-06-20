package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type PropertyVerificationRequest struct {
	ID                        int64                               `json:"id"`
	CreatorID                 uuid.UUID                           `json:"creatorId"`
	PropertyID                uuid.UUID                           `json:"propertyId"`
	VideoUrl                  string                              `json:"videoUrl"`
	HouseOwnershipCertificate *string                             `json:"houseOwnershipCertificate"`
	CertificateOfLanduseRight *string                             `json:"certificateOfLanduseRight"`
	FrontIdcard               string                              `json:"frontIdcard"`
	BackIdcard                string                              `json:"backIdcard"`
	Note                      *string                             `json:"note"`
	Feedback                  *string                             `json:"feedback"`
	Status                    database.PROPERTYVERIFICATIONSTATUS `json:"status"`
	CreatedAt                 time.Time                           `json:"createdAt"`
	UpdatedAt                 time.Time                           `json:"updatedAt"`
}

func ToPropertyVerificationRequest(p *database.PropertyVerificationRequest) PropertyVerificationRequest {
	return PropertyVerificationRequest{
		ID:                        p.ID,
		CreatorID:                 p.CreatorID,
		PropertyID:                p.PropertyID,
		VideoUrl:                  p.VideoUrl,
		HouseOwnershipCertificate: types.PNStr(p.HouseOwnershipCertificate),
		CertificateOfLanduseRight: types.PNStr(p.CertificateOfLanduseRight),
		FrontIdcard:               p.FrontIdcard,
		BackIdcard:                p.BackIdcard,
		Note:                      types.PNStr(p.Note),
		Feedback:                  types.PNStr(p.Feedback),
		Status:                    p.Status,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}
