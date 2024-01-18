package asynctask

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Processor interface {
	Start() error
	Shutdown()
	RegisterHandler(taskType string, handler asynq.HandlerFunc)
	// ProcessTask(context.Context, *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt) Processor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
		mux:    asynq.NewServeMux(),
	}
}

func (processor *RedisTaskProcessor) RegisterHandler(
	taskType string,
	handler asynq.HandlerFunc,
) {
	processor.mux.HandleFunc(taskType, handler)
}

func (processor *RedisTaskProcessor) Start() error {
	return processor.server.Start(processor.mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}

// func ProcessTaskSend(context.Context, *asynq.Task) error
