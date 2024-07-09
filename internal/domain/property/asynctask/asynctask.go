package asynctask

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
)

type Adapter interface {
	Register(processor asynctask.Processor)
}

type adapter struct {
	service service.Service
}

func NewAdapter(service service.Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) Register(processor asynctask.Processor) {
	processor.RegisterHandler(asynctask.PROPERTY_VERIFICATION_CREATE, a.notifyVerificationCreate)
	processor.RegisterHandler(asynctask.PROPERTY_VERIFICATION_UPDATE, a.notifyVerificationUpdate)
}

func (a *adapter) notifyVerificationCreate(ctx context.Context, task *asynq.Task) error {
	var payload model.PropertyVerificationRequest
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreatePropertyVerificationRequestStatus(&payload)
}

func (a *adapter) notifyVerificationUpdate(ctx context.Context, task *asynq.Task) error {
	var payload dto.UpdatePropertyVerificationRequestStatusNotification
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyUpdatePropertyVerificationRequestStatus(payload.Request, payload.UpdateData)
}
