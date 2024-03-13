package storage

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/auth/http"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/storage/dto"

	// "github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/token"
)

type Adapter interface {
	RegisterServer(route *fiber.Router, tokenMaker token.Maker)
}

type adapter struct {
	service Service
}

func NewAdapter(service Service) Adapter {
	return &adapter{
		service: service,
	}
}

func (a *adapter) RegisterServer(route *fiber.Router, tokenMaker token.Maker) {
	storageRoute := (*route).Group("/storage")

	// storageRoute.Use(http.AuthorizedMiddleware(tokenMaker))

	storageRoute.Post("/presign", a.getPresignUrl())
}

func (a *adapter) getPresignUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PutObjectPresignRequest
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}

		userId := uuid.New()
		tkPayload, ok := ctx.Locals(http.AuthorizationPayloadKey).(*token.Payload)
		if ok {
			userId = tkPayload.UserID
		}

		presignUrl, err := a.service.GetPresignUrl(&payload, userId)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return ctx.Status(fiber.StatusCreated).JSON(presignUrl)
	}
}
