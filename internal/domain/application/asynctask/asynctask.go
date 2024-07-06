package asynctask

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/domain/application/service"

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
	processor.RegisterHandler(asynctask.APPLICATION_NEW, a.notifyCreateApplication)
	processor.RegisterHandler(asynctask.APPLICATION_UPDATE, a.notifyUpdateApplication)
}

func (a *adapter) notifyCreateApplication(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyCreateApplication")
	var payload model.ApplicationModel
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.SendNotificationOnNewApplication(&payload)
}

func (a *adapter) notifyUpdateApplication(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyUpdateApplication")
	var payload dto.NotificationOnUpdateApplication
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.SendNotificationOnUpdateApplication(payload.Application, payload.Status)
}
