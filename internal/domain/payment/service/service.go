package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
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
	GetPaymentsOfUser(userId uuid.UUID, query *dto.GetPaymentsOfUserQuery) ([]model.PaymentModel, error)
	HandleReturn(data *dto.UpdatePayment, paymentInfo string) error
}

type PaymentService struct {
	repo repo.Repo
}

func NewPaymentService(repo repo.Repo) PaymentService {
	return PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) GetPaymentById(userId uuid.UUID, id int64) (*model.PaymentModel, error) {
	visible, err := s.repo.CheckPaymentAccessible(context.Background(), userId, id)
	if err != nil {
		return nil, err
	}
	if !visible {
		return nil, ErrInaccessiblePayment
	}

	return s.repo.GetPaymentById(context.Background(), id)
}

func (s *PaymentService) GetPaymentsOfUser(userId uuid.UUID, query *dto.GetPaymentsOfUserQuery) ([]model.PaymentModel, error) {
	var (
		limit  int32
		offset int32
	)
	if query.Limit == nil {
		limit = math.MaxInt32
	} else {
		limit = *query.Limit
	}
	if query.Offset == nil {
		offset = 0
	} else {
		offset = *query.Offset
	}
	return s.repo.GetPaymentsOfUser(context.Background(), userId, limit, offset)
}
