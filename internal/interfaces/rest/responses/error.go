package responses

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func ErrorResponse(err error) fiber.Map {
	return fiber.Map{"message": err.Error()}
}

func ValidationErrorResponse(ctx *fiber.Ctx, err validator.FieldError) {
	ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
}

func DBErrorResponse(ctx *fiber.Ctx, err *pq.Error) {
	switch err.Code[0:2] {
	case "22":
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	case "23":
		ctx.Status(http.StatusConflict).JSON(fiber.Map{"message": err.Error()})
	default:
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
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
