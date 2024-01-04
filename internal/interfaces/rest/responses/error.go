package responses

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func ErrorResponse(err error) fiber.Map {
	return fiber.Map{"message": err.Error()}
}

func ValidationErrorResponse(ctx *fiber.Ctx, err validator.FieldError) {
	ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
}

func DBErrorResponse(ctx *fiber.Ctx, err *pgconn.PgError) error {
	fmt.Println("db error:", err.Code, err.Message)
	switch err.Code[0:2] {
	case "22", "42":
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Message})
	case "23":
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Message})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
}

func DBTXErrorResponse(ctx *fiber.Ctx, err *database.TXError) error {
	if err.Err != nil {
		pgErr := err.Err.(*pgconn.PgError)
		return DBErrorResponse(ctx, pgErr)
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
}

// http error status response

type BadRequestError struct {
	Message string `json:"message"`
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", e.Message)
}

func HttpErrorResponse(e error) {

}
