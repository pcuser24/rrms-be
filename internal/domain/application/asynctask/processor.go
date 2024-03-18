package asynctask

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/email"
)

type TaskProcessor interface {
	RegisterProcessor()
	SendEmailOnNewApplication(ctx context.Context, task *asynq.Task) error
	UpdateApplicationStatus(ctx context.Context, task *asynq.Task) error
}

type taskProcessor struct {
	processor asynctask.Processor
	mailer    email.EmailSender
}

func NewTaskProcessor(p asynctask.Processor, m email.EmailSender) TaskProcessor {
	return &taskProcessor{
		processor: p,
		mailer:    m,
	}
}

func (p *taskProcessor) RegisterProcessor() {
	p.processor.RegisterHandler(TaskSendEmailOnNewApplication, p.SendEmailOnNewApplication)
	p.processor.RegisterHandler(TaskUpdateApplicationStatus, p.UpdateApplicationStatus)
}

func (p *taskProcessor) SendEmailOnNewApplication(ctx context.Context, task *asynq.Task) error {
	var payload dto.SendEmailOnNewApplicationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Println("Error unmarshalling payload", err)
		return err
	}
	log.Println("Processing TaskSendEmailOnNewApplication", payload)
	return nil
}

func (p *taskProcessor) UpdateApplicationStatus(ctx context.Context, task *asynq.Task) error {
	var payload dto.UpdateApplicationStatusPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Println("Error unmarshalling payload", err)
		return err
	}
	log.Println("Processing TaskUpdateApplicationStatus", payload)
	return nil
}
