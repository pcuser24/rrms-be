package service

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

func (s *service) CreateRentalPayment(data *dto.CreateRentalPayment) (model.RentalPayment, error) {
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), data.RentalID)
	if err != nil {
		return model.RentalPayment{}, err
	}

	// TODO: validate rental payment (code)
	res, err := s.domainRepo.RentalRepo.CreateRentalPayment(context.Background(), data)
	if err != nil {
		return model.RentalPayment{}, err
	}

	err = s.notifyCreateRentalPayment(&rental, &res)
	if err != nil {
		// TODO: log error
	}

	return res, nil
}

func (s *service) GetRentalPayment(id int64) (model.RentalPayment, error) {
	return s.domainRepo.RentalRepo.GetRentalPayment(context.Background(), id)
}

func (s *service) GetPaymentsOfRental(rentalID int64) ([]model.RentalPayment, error) {
	err := s.domainRepo.RentalRepo.UpdateFinePaymentsOfRental(context.Background(), rentalID)
	if err != nil {
		return nil, err
	}
	return s.domainRepo.RentalRepo.GetPaymentsOfRental(context.Background(), rentalID)
}

func (s *service) GetManagedRentalPayments(uid uuid.UUID, query *dto.GetManagedRentalPaymentsQuery) ([]model.RentalPayment, error) {
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

// UpdateRentalPayment updates the rental payment with the given id, where status is the assumed current status of the rental payment and data.Status, if present, is the new status of the rental payment.
func (s *service) UpdateRentalPayment(id int64, userId uuid.UUID, data dto.IUpdateRentalPayment, status database.RENTALPAYMENTSTATUS) error {
	rp, err := s.domainRepo.RentalRepo.GetRentalPayment(context.Background(), id)
	if err != nil {
		return err
	}
	r, err := s.domainRepo.RentalRepo.GetRental(context.Background(), rp.RentalID)
	if err != nil {
		return err
	}

	side, err := s.domainRepo.RentalRepo.GetRentalSide(context.Background(), rp.RentalID, userId)
	if err != nil {
		return err
	}

	var _data dto.UpdateRentalPayment

	payFine := func(r *model.RentalModel, rp *model.RentalPayment) dto.UpdateRentalPayment {
		var amount float32 = 0
		switch r.LatePaymentPenaltyScheme {
		case database.LATEPAYMENTPENALTYSCHEMEPERCENT:
			amount = rp.MustPay * (1 + *r.LatePaymentPenaltyAmount/100)
		case database.LATEPAYMENTPENALTYSCHEMEFIXED:
			amount = rp.MustPay + *r.LatePaymentPenaltyAmount
		default:
			amount = rp.MustPay
		}

		return dto.UpdateRentalPayment{
			Fine:   &amount,
			Status: database.RENTALPAYMENTSTATUSPAYFINE,
		}
	}

	if rp.Status != status {
		return ErrInvalidPaymentTypeTransition
	}

	var willNotify bool = true
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
		if _data.Status == database.RENTALPAYMENTSTATUSPAID || _data.Status == database.RENTALPAYMENTSTATUSCANCELLED {
			willNotify = false
		}
		// log.Println("Send notification to tenant: Issue new rental payment")
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
		// log.Println("Send notification to managers:", utils.Ternary(__data.Status == database.RENTALPAYMENTSTATUSPENDING, "Tenant agree with the issued payment", "Tenant request a payment review"))
	case database.RENTALPAYMENTSTATUSPENDING:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "B" {
			return ErrInvalidPaymentTypeTransition
		}
		// update status to PAYFINE if rp.expiryDate + r.gracePeriod day < today
		if time.Now().After(rp.ExpiryDate.AddDate(0, 0, int(r.GracePeriod))) && rp.MustPay > 0 {
			_data = payFine(&r, &rp)
			// log.Println("Send notification to tenant: Pay fine")
			// log.Println("Send notification to managers: Tenant is late for a payment")
		} else {
			_data = dto.UpdateRentalPayment{
				ID:          id,
				UserID:      userId,
				PaymentDate: __data.PaymentDate,
				Payamount:   types.Ptr(__data.PayAmount),
				Status:      database.RENTALPAYMENTSTATUSREQUEST2PAY,
			}
			// log.Println("Send notification to managers: Tenant has update his payment status, review now")
		}
	case database.RENTALPAYMENTSTATUSREQUEST2PAY:
		__data := data.(*dto.UpdatePendingRentalPayment)
		if side != "A" {
			return ErrInvalidPaymentTypeTransition
		}
		// update status to PAYFINE if rp.expiryDate + r.gracePeriod day < today
		if time.Now().After(rp.ExpiryDate.AddDate(0, 0, int(r.GracePeriod))) && rp.MustPay > 0 {
			_data = payFine(&r, &rp)
			// log.Println("Send notification to tenant: Pay fine")
			// log.Println("Send notification to managers: Tenant is late for a payment")
		} else {
			_data = dto.UpdateRentalPayment{
				ID:          id,
				UserID:      userId,
				PaymentDate: __data.PaymentDate,
				Payamount:   types.Ptr(__data.PayAmount),
				Paid:        types.Ptr(rp.Paid + __data.PayAmount),
				Status: utils.Ternary(
					rp.Paid+__data.PayAmount < rp.MustPay,
					database.RENTALPAYMENTSTATUSPARTIALLYPAID,
					database.RENTALPAYMENTSTATUSPAID,
				),
			}
			// log.Println("Send notification to tenant: your payment is recorded")
		}
	case database.RENTALPAYMENTSTATUSPARTIALLYPAID:
		__data := data.(*dto.UpdatePartiallyPaidRentalPayment)
		if side != "B" {
			return ErrInvalidPaymentTypeTransition
		}
		// update status to PAYFINE if rp.expiryDate + r.gracePeriod day < today
		if time.Now().After(rp.ExpiryDate.AddDate(0, 0, int(r.GracePeriod))) && rp.MustPay > 0 {
			_data = payFine(&r, &rp)
			// log.Println("Send notification to tenant: Pay fine")
			// log.Println("Send notification to managers: Tenant is late for a payment")
		} else {
			_data = dto.UpdateRentalPayment{
				ID:          id,
				UserID:      userId,
				Payamount:   types.Ptr(__data.PayAmount),
				PaymentDate: __data.PaymentDate,
				Status:      database.RENTALPAYMENTSTATUSREQUEST2PAY,
			}
			// log.Println("Send notification to managers: Tenant has update his payment status, review now")
		}
	case database.RENTALPAYMENTSTATUSPAYFINE:
		if side != "A" {
			return ErrInvalidPaymentTypeTransition
		}
		_data = dto.UpdateRentalPayment{
			ID:        id,
			UserID:    userId,
			Payamount: rp.Fine,
			Paid:      rp.Fine,
			Status:    database.RENTALPAYMENTSTATUSPAID,
		}
		// log.Println("Send notification to tenant: your fine payment is done, good job")
	default:
		return ErrInvalidPaymentTypeTransition
	}

	err = s.domainRepo.RentalRepo.UpdateRentalPayment(context.Background(), &_data)
	if err != nil {
		return err
	}
	if willNotify {
		err = s.notifyUpdatePayments(&r, &rp, &_data)
		if err != nil {
			// TODO: log error
		}
	}
	return nil
}
