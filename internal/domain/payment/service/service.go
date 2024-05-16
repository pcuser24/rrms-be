package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
)

type PAYMENTTYPE string

const PAYMENTTYPE_DELIMITER = "_"
const (
	PAYMENTTYPE_CREATELISTING  PAYMENTTYPE = "CREATELISTING"
	PAYMENTTYPE_EXTENDLISTING  PAYMENTTYPE = "EXTENDLISTING"
	PAYMENTTYPE_UPGRADELISTING PAYMENTTYPE = "UPGRADELISTING"
)

var (
	ErrInvalidPaymentInfo  = errors.New("invalid payment info")
	ErrInvalidPaymentType  = errors.New("invalid payment type")
	ErrInaccessiblePayment = errors.New("inaccessible payment")
)

type Service interface {
	GetPaymentById(userId uuid.UUID, id int64) (*model.PaymentModel, error)

	HandleReturn(data *dto.UpdatePayment, paymentInfo string) error
}
