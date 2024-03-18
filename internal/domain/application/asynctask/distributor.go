package asynctask

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
)

/* -------------------------------------------------------------------------- */
/*                                 Task types                                 */
/* -------------------------------------------------------------------------- */

const (
	TaskSendEmailOnNewApplication = "application:new:succesful"
	TaskUpdateApplicationStatus   = "application:update:status"
)

/* -------------------------------------------------------------------------- */
/*                                 Distributor                                */
/* -------------------------------------------------------------------------- */

type TaskDistributor interface {
	SendEmailOnNewApplication(ctx context.Context, payload *dto.SendEmailOnNewApplicationPayload) error
	UpdateApplicationStatus(ctx context.Context, payload *dto.UpdateApplicationStatusPayload) error
}

type taskDistributor struct {
	distributor asynctask.Distributor
}

func NewTaskDistributor(d asynctask.Distributor) TaskDistributor {
	return &taskDistributor{
		distributor: d,
	}
}

func (d *taskDistributor) SendEmailOnNewApplication(ctx context.Context, payload *dto.SendEmailOnNewApplicationPayload) error {
	return d.distributor.DistributeTaskJSON(ctx, TaskSendEmailOnNewApplication, payload)
}

func (d *taskDistributor) UpdateApplicationStatus(ctx context.Context, payload *dto.UpdateApplicationStatusPayload) error {
	return d.distributor.DistributeTaskJSON(ctx, TaskUpdateApplicationStatus, payload)
}
