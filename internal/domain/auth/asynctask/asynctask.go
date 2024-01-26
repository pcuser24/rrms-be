package asynctask

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/email"
)

/* -------------------------------------------------------------------------- */
/*                                 Task types                                 */
/* -------------------------------------------------------------------------- */

const (
	TaskSendVerifyEmail = "auth:renew_password"
)

/* -------------------------------------------------------------------------- */
/*                                 Distributor                                */
/* -------------------------------------------------------------------------- */

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *dto.TaskSendVerifyEmailPayload,
		opts ...asynq.Option,
	) error
}

type distributor struct {
	distributor asynctask.Distributor
}

func NewTaskDistributor(d asynctask.Distributor) TaskDistributor {
	return &distributor{
		distributor: d,
	}
}

func (d *distributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *dto.TaskSendVerifyEmailPayload,
	opts ...asynq.Option,
) error {
	return d.distributor.DistributeTaskJSON(
		ctx,
		TaskSendVerifyEmail,
		payload,
		opts...,
	)
}

/* -------------------------------------------------------------------------- */
/*                                  Processor                                 */
/* -------------------------------------------------------------------------- */

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	RegisterProcessor()
}

type processor struct {
	processor asynctask.Processor
	mailer    email.EmailSender
}

func NewTaskProcessor(p asynctask.Processor, m email.EmailSender) TaskProcessor {
	return &processor{
		processor: p,
		mailer:    m,
	}
}

func (p *processor) RegisterProcessor() {
	p.processor.RegisterHandler(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
}

func (p *processor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	return nil
}
