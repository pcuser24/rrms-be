package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

func (s *service) CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error) {
	// TODO: validate rental payment (code)
	return s.domainRepo.RentalRepo.CreateRentalPayment(context.Background(), data)
}

func (s *service) GetRentalPayment(id int64) (model.RentalPayment, error) {
	return s.domainRepo.RentalRepo.GetRentalPayment(context.Background(), id)
}

func (s *service) GetPaymentsOfRental(rentalID int64) ([]model.RentalPayment, error) {
	return s.domainRepo.RentalRepo.GetPaymentsOfRental(context.Background(), rentalID)
}

func (s *service) GetManagedRentalPayments(uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) ([]dto.GetManagedRentalPaymentsItem, error) {
	if query.Limit == nil {
		query.Limit = types.Ptr[int32](math.MaxInt32)
	}
	if query.Offset == nil {
		query.Offset = types.Ptr[int32](0)
	}
	statusSet := set.NewSet[database.RENTALPAYMENTSTATUS]()
	statusSet.AddAll(query.Status...)
	query.Status = statusSet.ToSlice()
	return s.domainRepo.RentalRepo.GetManagedRentalPayments(context.Background(), uid, query)
}

func (s *service) UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error {
	rp, err := s.domainRepo.RentalRepo.GetRentalPayment(context.Background(), id)
	if err != nil {
		return err
	}
	side, err := s.domainRepo.RentalRepo.GetRentalSide(context.Background(), rp.RentalID, userId)
	if err != nil {
		return err
	}

	var _data dto.UpdateRentalPayment

	if rp.Status != status {
		return ErrInvalidPaymentTypeTransition
	}
	switch status {
	case database.RENTALPAYMENTSTATUSPLAN:
		__data := data.(*dto.UpdatePlanRentalPayment)
		if side != "A" {
			return ErrInvalidPaymentTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:         id,
			UserID:     userId,
			Status:     __data.Status,
			Amount:     &__data.Amount,
			Discount:   __data.Discount,
			ExpiryDate: __data.ExpiryDate,
		}
	case database.RENTALPAYMENTSTATUSISSUED:
		__data := data.(*dto.UpdateIssuedRentalPayment)
		if side != "B" {
			return ErrInvalidPaymentTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:     id,
			UserID: userId,
			Status: __data.Status,
			Note:   __data.Note,
		}
	case database.RENTALPAYMENTSTATUSPENDING:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "B" {
			return ErrInvalidPaymentTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:          id,
			UserID:      userId,
			PaymentDate: __data.PaymentDate,
			Status:      database.RENTALPAYMENTSTATUSREQUEST2PAY,
		}
	case database.RENTALPAYMENTSTATUSREQUEST2PAY:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "A" {
			return ErrInvalidPaymentTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:          id,
			UserID:      userId,
			PaymentDate: __data.PaymentDate,
			Status:      database.RENTALPAYMENTSTATUSPAID,
		}
	default:
		return ErrInvalidPaymentTypeTransition
	}
	return s.domainRepo.RentalRepo.UpdateRentalPayment(context.Background(), &_data)

}
