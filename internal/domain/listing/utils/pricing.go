package utils

import (
	"errors"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/listing/model"
)

var (
	// map priority level to base price
	listingPriorities = map[int]int{
		1: 2000,
		2: 5000,
		3: 7000,
		4: 9000,
	}
	// map post duration to discount percentage
	listingDiscounts = map[int]int{
		7:  0,
		15: 10,
		30: 20,
	}
)

var (
	ErrInvalidPriority = errors.New("invalid priority")
	ErrInvalidDuration = errors.New("invalid duration")
)

func CalculateListingPrice(priority int, postDuration int) (int64, int, error) {
	p := listingPriorities[priority]
	if p == 0 {
		return 0, 0, ErrInvalidPriority
	}
	discount := listingDiscounts[postDuration]

	return int64(p-(p*discount/100)) * int64(postDuration), discount, nil
}

func CalculateUpgradeListingPrice(l *model.ListingModel, p int) (int64, int, error) {
	if p <= 0 || p > 4 || p <= int(l.Priority) {
		return 0, 0, ErrInvalidPriority
	}
	daysLeft := time.Since(l.CreatedAt).Hours() / 24
	oldBasePrice := listingPriorities[int(l.Priority)]
	newBasePrice := listingPriorities[p]

	return int64((newBasePrice - oldBasePrice) * int(daysLeft)), 0, nil
}

func CalculateExtendListingPrice(l *model.ListingModel, d int) (int64, int, error) {
	if d <= 0 {
		return 0, 0, ErrInvalidDuration
	}
	return int64(listingPriorities[int(l.Priority)] * d), 0, nil
}
