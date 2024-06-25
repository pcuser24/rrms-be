package http

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/validation"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

func (a *adapter) getRecentListings() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		type Query struct {
			Fields []string `query:"fields" validate:"listingFields"`
			Limit  int32    `query:"limit" validate:"required,min=0"`
		}
		var query Query
		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if len(query.Fields) == 1 {
			fieldSet := set.NewSet[string]()
			fieldSet.AddAll(strings.Split(query.Fields[0], ",")...)
			query.Fields = fieldSet.ToSlice()
		}
		validator := validator.New()
		validator.RegisterValidation(listing_dto.ListingFieldsLocalKey, listing_dto.ValidateQuery)
		if errs := validation.ValidateStruct(validator, query); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.service.GetRecentListings(query.Limit, query.Fields)
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}

func (a *adapter) getListingSuggestions() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		listingId, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid listingId"})
		}
		limitQuery := ctx.Query("limit")
		var limit int64
		if limitQuery != "" {
			limit, err = (strconv.ParseInt(limitQuery, 10, 32))
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid limit"})
			}
		} else {
			limit = 16
		}

		res, err := a.service.GetSimilarListingsToListing(listingId, int(limit))
		if err != nil {
			if errors.Is(err, database.ErrRecordNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return ctx.Status(fiber.StatusOK).JSON(res)
	}
}
