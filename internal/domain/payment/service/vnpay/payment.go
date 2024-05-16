package vnpay

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (s *VnPayService) GetPaymentById(userId uuid.UUID, id int64) (*model.PaymentModel, error) {
	visible, err := s.repo.CheckPaymentAccessible(context.Background(), userId, id)
	if err != nil {
		return nil, err
	}
	if !visible {
		return nil, service.ErrInaccessiblePayment
	}

	return s.repo.GetPaymentById(context.Background(), id)
}

func (s *VnPayService) HandleReturn(data *dto.UpdatePayment, paymentInfo string) error {
	end := strings.Index(paymentInfo, "]")
	if end == -1 || paymentInfo[0] != '[' {
		return service.ErrInvalidPaymentInfo
	}
	d := strings.Index(paymentInfo, service.PAYMENTTYPE_DELIMITER)
	if d == -1 {
		return service.ErrInvalidPaymentInfo
	}
	paymentType := paymentInfo[1:d]
	paymentObject := paymentInfo[d+1 : end]
	success := (*data.Status == database.PAYMENTSTATUSSUCCESS)
	switch service.PAYMENTTYPE(paymentType) {
	case service.PAYMENTTYPE_CREATELISTING:
		return s.handlePayCreateListing(paymentObject, success)
	case service.PAYMENTTYPE_EXTENDLISTING:
		return s.handlePayExtendListing(paymentObject)
	case service.PAYMENTTYPE_UPGRADELISTING:
		return s.handlePayUpgradeListing(paymentObject)
	default:
		return service.ErrInvalidPaymentType
	}
}

func (s *VnPayService) handlePayCreateListing(listingId string, success bool) error {
	id, err := uuid.Parse(listingId)
	if err != nil {
		return err
	}
	return s.lRepo.UpdateListingStatus(context.Background(), id, success)
}

func (s *VnPayService) handlePayExtendListing(paymentObject string) error {
	d := strings.Index(paymentObject, service.PAYMENTTYPE_DELIMITER)
	if d == -1 {
		return service.ErrInvalidPaymentInfo
	}
	listingId, err := uuid.Parse(paymentObject[:d])
	if err != nil {
		return service.ErrInvalidPaymentInfo
	}
	duration, err := strconv.ParseInt(paymentObject[d+1:], 10, 64)
	if err != nil {
		return service.ErrInvalidPaymentInfo
	}

	return s.lRepo.UpdateListingExpiration(context.Background(), listingId, duration)
}

func (s *VnPayService) handlePayUpgradeListing(paymentObject string) error {
	d := strings.Index(paymentObject, service.PAYMENTTYPE_DELIMITER)
	if d == -1 {
		return service.ErrInvalidPaymentInfo
	}
	listingId, err := uuid.Parse(paymentObject[:d])
	if err != nil {
		return service.ErrInvalidPaymentInfo
	}
	priority, err := strconv.ParseInt(paymentObject[d+1:], 10, 64)
	if err != nil {
		return service.ErrInvalidPaymentInfo
	}

	return s.lRepo.UpdateListingPriority(context.Background(), listingId, int(priority))
}
