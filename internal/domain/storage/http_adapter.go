package storage

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/user2410/rrms-backend/internal/domain/auth"
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

	storageRoute.Use(auth.AuthorizedMiddleware(tokenMaker))

	storageRoute.Post("/presign", a.getPresignUrl())
}

func (a *adapter) getPresignUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var payload dto.PutObjectPresign
		if err := ctx.BodyParser(&payload); err != nil {
			return err
		}
		fmt.Println(payload)
		tkPayload := ctx.Locals(auth.AuthorizationPayloadKey).(*token.Payload)

		presignUrl, err := a.service.GetPresignUrl(&payload, tkPayload.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return ctx.Status(fiber.StatusOK).JSON(presignUrl)
	}
}
