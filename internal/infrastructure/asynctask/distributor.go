package asynctask

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type Distributor interface {
	DistributeTask(
		ctx context.Context,
		taskType string,
		payload []byte,
		opts ...asynq.Option,
	) error
	DistributeTaskJSON(
		ctx context.Context,
		taskType string,
		payload any,
		opts ...asynq.Option,
	) error
	Close() error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) Distributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}

func (distributor *RedisTaskDistributor) Close() error {
	return distributor.client.Close()
}

var (
	ErrMarshalPayload = fmt.Errorf("failed to marshal task payload")
	ErrEnqueueTask    = fmt.Errorf("failed to enqueue task")
)

func (d *RedisTaskDistributor) DistributeTask(
	ctx context.Context,
	taskType string,
	payload []byte,
	opts ...asynq.Option,
) error {
	task := asynq.NewTask(taskType, payload, opts...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return ErrEnqueueTask
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (d *RedisTaskDistributor) DistributeTaskJSON(
	ctx context.Context,
	taskType string,
	payload any,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return ErrMarshalPayload
	}

	return d.DistributeTask(ctx, taskType, jsonPayload, opts...)
}

const (
	APPLICATION_NEW    = "applications/new"
	APPLICATION_UPDATE = "applications/application/update"

	RENTAL_PRERENTAL_NEW           = "rentals/prerental/new"
	RENTAL_PRERENTAL_UPDATE        = "rentals/prerental/update"
	RENTAL_PAYMENT_CREATE          = "rentals/payment/create"
	RENTAL_PAYMENT_UPDATE          = "rentals/payment/update"
	RENTAL_CONTRACT_CREATE         = "rentals/contract/create"
	RENTAL_CONTRACT_UPDATE         = "rentals/contract/update"
	RENTAL_COMPLAINT_CREATE        = "rentals/complaint/create"
	RENTAL_COMPLAINT_REPLY         = "rentals/complaint/reply"
	RENTAL_COMPLAINT_STATUS_UPDATE = "rentals/complaint/status/update"
)
