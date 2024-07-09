package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	"github.com/user2410/rrms-backend/internal/utils"
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

func init() {
	// read service name from json config file "./services.json"
	basepath := utils.GetBasePath()
	file, err := os.Open(fmt.Sprintf("%s/internal/config/listing.json", basepath))
	if err != nil {
		panic(err)
	}

	type Data struct {
		PrioritiesToPrice map[string]int `json:"prioritiesToPrice"`
		DurationDiscounts map[string]int `json:"durationDiscounts"`
	}
	var data Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	for key, value := range data.PrioritiesToPrice {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			panic(err)
		}
		listingPriorities[intKey] = value
	}

	for key, value := range data.DurationDiscounts {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			panic(err)
		}
		listingDiscounts[intKey] = value
	}
}

var (
	ErrInvalidPriority = errors.New("invalid priority")
	ErrInvalidDuration = errors.New("invalid duration")
)

func CalculateListingPrice(priority int, postDuration int) (float32, int, int, error) {
	p := listingPriorities[priority]
	if p == 0 {
		return 0, 0, 0, ErrInvalidPriority
	}
	discount := listingDiscounts[postDuration]

	return float32((p - (p*discount)/100) * postDuration), p, discount, nil
}

func CalculateUpgradeListingPrice(l *model.ListingModel, p int) (float32, int, error) {
	if p <= 0 || p > 4 || p <= int(l.Priority) {
		return 0, 0, ErrInvalidPriority
	}
	daysLeft := time.Since(l.CreatedAt).Hours() / 24
	oldBasePrice := listingPriorities[int(l.Priority)]
	newBasePrice := listingPriorities[p]

	return float32((newBasePrice - oldBasePrice) * int(daysLeft)), 0, nil
}

func CalculateExtendListingPrice(l *model.ListingModel, d int) (float32, int, error) {
	if d <= 0 {
		return 0, 0, ErrInvalidDuration
	}
	// return int64(listingPriorities[int(l.Priority)] * d), 0, nil
	return float32(listingPriorities[int(l.Priority)] * d), 0, nil
}
