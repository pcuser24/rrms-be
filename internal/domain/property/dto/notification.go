package dto

import (
	"github.com/user2410/rrms-backend/internal/domain/property/model"
)

type UpdatePropertyVerificationRequestStatusNotification struct {
	Request    *model.PropertyVerificationRequest       `json:"request"`
	UpdateData *UpdatePropertyVerificationRequestStatus `json:"updateData"`
}
