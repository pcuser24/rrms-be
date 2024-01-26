package asynctask

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/email"
)

/* -------------------------------------------------------------------------- */
/*                                 Task types                                 */
/* -------------------------------------------------------------------------- */

const (
	TaskSendEmailOnNewApplication = "application:new:succesful"
)

/* -------------------------------------------------------------------------- */
/*                                 Distributor                                */
/* -------------------------------------------------------------------------- */

type TaskDistributor interface {
	DistributeTaskSendEmailOnNewApplication(
		ctx context.Context,
		payload *dto.TaskSendEmailOnNewApplicationPayload,
		opts ...asynq.Option,
	) error
}

type taskDistributor struct {
	distributor asynctask.Distributor
}

func NewTaskDistributor(d asynctask.Distributor) TaskDistributor {
	return &taskDistributor{
		distributor: d,
	}
}

func (d *taskDistributor) DistributeTaskSendEmailOnNewApplication(
	ctx context.Context,
	payload *dto.TaskSendEmailOnNewApplicationPayload,
	opts ...asynq.Option,
) error {
	return d.distributor.DistributeTaskJSON(
		ctx,
		TaskSendEmailOnNewApplication,
		payload,
		opts...,
	)
}

/* -------------------------------------------------------------------------- */
/*                                  Processor                                 */
/* -------------------------------------------------------------------------- */

type TaskProcessor interface {
	ProcessTaskSendEmailOnNewApplication(ctx context.Context, task *asynq.Task) error
	RegisterProcessor()
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
	p.processor.RegisterHandler(
		TaskSendEmailOnNewApplication,
		p.ProcessTaskSendEmailOnNewApplication,
	)
}

func (p *taskProcessor) ProcessTaskSendEmailOnNewApplication(ctx context.Context, task *asynq.Task) error {
	return nil
}
