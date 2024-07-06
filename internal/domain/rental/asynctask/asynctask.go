package asynctask

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/service"
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
	processor.RegisterHandler(asynctask.RENTAL_PRERENTAL_NEW, a.notifyCreatePreRental)
	processor.RegisterHandler(asynctask.RENTAL_PRERENTAL_UPDATE, a.notifyUpdatePrerental)
	processor.RegisterHandler(asynctask.RENTAL_PAYMENT_CREATE, a.notifyCreatePayment)
	processor.RegisterHandler(asynctask.RENTAL_PAYMENT_UPDATE, a.notifyUpdatePayment)
	processor.RegisterHandler(asynctask.RENTAL_CONTRACT_CREATE, a.notifyCreateContract)
	processor.RegisterHandler(asynctask.RENTAL_CONTRACT_UPDATE, a.notifyUpdateContract)
	processor.RegisterHandler(asynctask.RENTAL_COMPLAINT_CREATE, a.notifyCreateComplaint)
	processor.RegisterHandler(asynctask.RENTAL_COMPLAINT_REPLY, a.notifyReplyComplaint)
	processor.RegisterHandler(asynctask.RENTAL_COMPLAINT_STATUS_UPDATE, a.notifyUpdateComplaintStatus)
}

func (a *adapter) notifyCreatePreRental(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyCreatePreRental")
	var payload dto.NotifyCreatePreRental
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreatePreRental(payload.Rental, payload.Secret)
}

func (a *adapter) notifyUpdatePrerental(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyUpdatePrerental")
	var payload dto.NotifyUpdatePreRental
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyUpdatePreRental(payload.PreRental, payload.Rental, payload.UpdateData)
}

func (a *adapter) notifyCreatePayment(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyCreatePayment")
	var payload dto.NotifyCreateRentalPayment
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreateRentalPayment(payload.Rental, payload.RentalPayment)
}

func (a *adapter) notifyUpdatePayment(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyUpdatePayment")
	var payload dto.NotifyUpdatePayments
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyUpdatePayments(payload.Rental, payload.RentalPayment, payload.UpdateData)
}

func (a *adapter) notifyCreateContract(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyCreateContract")
	var payload dto.NotifyCreateContract
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreateContract(payload.Contract, payload.Rental)
}

func (a *adapter) notifyUpdateContract(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyUpdateContract")
	var payload dto.NotifyUpdateContract
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyUpdateContract(payload.Contract, payload.Rental, payload.Side)
}

func (a *adapter) notifyCreateComplaint(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyCreateComplaint")
	var payload dto.NotifyCreateRentalComplaint
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreateRentalComplaint(payload.Complaint, payload.Rental)
}

func (a *adapter) notifyReplyComplaint(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyReplyComplaint")
	var payload dto.NotifyCreateComplaintReply
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyCreateComplaintReply(payload.Complaint, payload.ComplaintReply, payload.Rental)
}

func (a *adapter) notifyUpdateComplaintStatus(ctx context.Context, task *asynq.Task) error {
	log.Println("notifyUpdateComplaintStatus")
	var payload dto.NotifyUpdateComplaintStatus
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return a.service.NotifyUpdateComplaintStatus(payload.Complaint, payload.Rental, payload.Status, payload.UpdatedBy)
}
