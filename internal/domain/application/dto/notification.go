package dto

import (
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type NotificationOnUpdateApplication struct {
	Application *model.ApplicationModel    `json:"application"`
	Status      database.APPLICATIONSTATUS `json:"status"`
}
