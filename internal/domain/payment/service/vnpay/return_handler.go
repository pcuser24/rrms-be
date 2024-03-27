package vnpay

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type PAYMENTTYPE string

const PAYMENTTYPE_DELIMITER = "-"
const (
	PAYMENTTYPE_PAYLISTING = "CREATELISTING"
)

var (
	ErrInvalidPaymentInfo = errors.New("invalid payment info")
	ErrInvalidPaymentType = errors.New("invalid payment type")
)

// HandleReturn handles the return from VNPAY. paymentInfo is in the form "[paymentType-paymentObject]..."
func (s *Service) handleReturn(paymentInfo string, success bool) error {
	end := strings.Index(paymentInfo, "]")
	if end == -1 || paymentInfo[0] != '[' {
		return ErrInvalidPaymentInfo
	}
	d := strings.Index(paymentInfo, PAYMENTTYPE_DELIMITER)
	if d == -1 {
		return ErrInvalidPaymentInfo
	}
	paymentType := paymentInfo[1:d]
	paymentObject := paymentInfo[d+1 : end]
	switch PAYMENTTYPE(paymentType) {
	case PAYMENTTYPE_PAYLISTING:
		return s.handlePayListing(paymentObject, success)
	default:
		return ErrInvalidPaymentType
	}
}

func (s *Service) handlePayListing(listingId string, success bool) error {
	id, err := uuid.Parse(listingId)
	if err != nil {
		return err
	}
	return s.lRepo.UpdateListingStatus(context.Background(), id, success)
}
