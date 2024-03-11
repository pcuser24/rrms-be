package utils

import "errors"

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
)

func CalculateListingPrice(priority int, postDuration int) (int64, error) {
	p := listingPriorities[priority]
	if p == 0 {
		return 0, ErrInvalidPriority
	}
	discount := listingDiscounts[postDuration]

	return int64(p-(p*discount/100)) * int64(postDuration), nil
}
