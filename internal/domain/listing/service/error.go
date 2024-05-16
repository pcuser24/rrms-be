package service

import "errors"

var (
	ErrUpgradeListingInvalidPriority = errors.New("invalid priority")
	ErrUnpaidPayment                 = errors.New("unpaid payment")
)
