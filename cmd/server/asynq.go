package server

import (
	"log"

	"github.com/hibiken/asynq"
	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	auth_asynctask "github.com/user2410/rrms-backend/internal/domain/auth/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/email"
)

func (c *serverCommand) setupAsyncTaskProcessor(
	mailer email.EmailSender,
) {
	c.asyncTaskProcessor = asynctask.NewRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: c.config.AsynqRedisAddress,
	})

	auth_asynctask.NewTaskProcessor(c.asyncTaskProcessor, mailer).RegisterProcessor()
	application_asynctask.NewTaskProcessor(c.asyncTaskProcessor, mailer).RegisterProcessor()
}

func (c *serverCommand) runAsyncTaskProcessor() {
	log.Println("Starting async task processor...")
	if err := c.asyncTaskProcessor.Start(); err != nil {
		log.Fatal("Failed to start task processor:", err)
	}
}
