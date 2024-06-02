package vnpay

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	payment_dto "github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (s *VnPayService) HandleReturn(data *payment_dto.UpdatePayment, paymentInfo string) error {
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
	return s.lService.UpdateListingStatus(id, success)
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

	return s.lService.UpdateListingExpiration(listingId, duration)
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

	return s.lService.UpdateListingPriority(listingId, int(priority))
}
