package service

import "errors"

var (
	ErrInvalidRentalExpired         = errors.New("rental expired")
	ErrInvalidPaymentTypeTransition = errors.New("invalid type transition")
)
