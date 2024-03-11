package http

import (
	"errors"
	"io"
	"maps"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	auth_http "github.com/user2410/rrms-backend/internal/domain/auth/http"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/responses"
	"github.com/user2410/rrms-backend/internal/utils/token"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

// Tạo URL Thanh toán
// URL Thanh toán là địa chỉ URL mang thông tin thanh toán.
// Website TMĐT gửi sang Cổng thanh toán VNPAY các thông tin này khi xử lý giao dịch thanh toán trực tuyến cho Khách mua hàng.
func (a *adapter) vnpCreatePaymentUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ipAddr := ctx.IP()

		payload := new(dto.VNPCreatePaymentUrl)
		if err := ctx.BodyParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		tkPayload := ctx.Locals(auth_http.AuthorizationPayloadKey).(*token.Payload)
		paymentId, err := strconv.ParseInt(ctx.Params("paymentId"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		url, err := a.vnpayService.CreatePaymentUrl(ipAddr, tkPayload.UserID, paymentId, payload)
		if err != nil {
			if errors.Is(err, vnpay.ErrInvalidHash) || errors.Is(err, vnpay.ErrBadStatusPayment) {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		// return ctx.Redirect(url)
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"url": url})
	}
}

// URL VNPAY gọi về sau khi thanh toán
// TODO: make a separates frontend only accept this URL
func (a *adapter) vnpReturn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(dto.VNPReturnQuery)
		if err := ctx.QueryParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		queries := maps.Clone(ctx.Queries())
		err := a.vnpayService.Return(queries)
		if err != nil {
			if errors.Is(err, vnpay.ErrInvalidHash) {
				return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"code": "97"})
			}
			if dbErr, ok := err.(*pgconn.PgError); ok {
				return responses.DBErrorResponse(ctx, dbErr)
			}
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"code": queries["vnp_ResponseCode"]})
	}
}

// Nhận kết quả thanh toán từ VNPAY
// Trên URL VNPAY gọi về có mang thông tin thanh toán để căn cứ vào kết quả đó Website TMĐT xử lý các bước tiếp theo (ví dụ: cập nhật kết quả thanh toán vào Database …)
func (a *adapter) vnpIpn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(dto.VNPIpnQuery)
		if err := ctx.QueryParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		queries := maps.Clone(ctx.Queries())
		ret := a.vnpayService.Ipn(queries)
		return ctx.Status(fiber.StatusOK).JSON(ret)
	}
}

func (a *adapter) vnpQuerydr() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ipAddr := ctx.IP()

		payload := new(dto.VNPQuerydr)
		if err := ctx.QueryParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.vnpayService.Querydr(ipAddr, payload)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(res.StatusCode).Type(res.Header.Get("Content-Type")).Send(body)
	}
}

func (a *adapter) vnpRefund() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ipAddr := ctx.IP()

		payload := new(dto.VNPRefund)
		if err := ctx.QueryParser(payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errs := validation.ValidateStruct(nil, payload); len(errs) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": validation.GetValidationError(errs)})
		}

		res, err := a.vnpayService.Refund(ipAddr, payload)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.Status(res.StatusCode).Type(res.Header.Get("Content-Type")).Send(body)
	}
}
