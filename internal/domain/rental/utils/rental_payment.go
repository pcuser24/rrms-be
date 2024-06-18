package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/utils"
)

// enum for rental payment type: "RENTAL", "ELECTRICITY", "WATER", "SERVICES"
type RentalPaymentType string

const (
	RENTALPAYMENTTYPERENTAL      RentalPaymentType = "RENTAL"
	RENTALPAYMENTTYPEDEPOSIT     RentalPaymentType = "DEPOSIT"
	RENTALPAYMENTTYPEELECTRICITY RentalPaymentType = "ELECTRICITY"
	RENTALPAYMENTTYPEWATER       RentalPaymentType = "WATER"
	RENTALPAYMENTTYPESERVICE     RentalPaymentType = "SERVICE"
	RENTALPAYMENTTYPEMAINTENANCE RentalPaymentType = "MAINTENANCE"
)

var (
	// map of rental payment type to its corresponding service name
	mapRentalPaymentTypeToServiceName map[RentalPaymentType]string

	ErrInvalidRentalPaymentCode = fmt.Errorf("invalid rental payment code")
)

func init() {
	// read service name from json config file "./services.json"
	basepath := utils.GetBasePath()
	file, err := os.ReadFile(fmt.Sprintf("%s/internal/config/rental_services.json", basepath))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &mapRentalPaymentTypeToServiceName)
	if err != nil {
		panic(err)
	}
	log.Println("Service names loaded successfully", mapRentalPaymentTypeToServiceName)
}

func GetRentalPaymentCode(rentalID int64, paymentType RentalPaymentType, serviceId int64, startDate, endDate time.Time) string {
	if paymentType == RENTALPAYMENTTYPESERVICE && serviceId != 0 {
		return fmt.Sprintf("%d_%s_%d_%02d%d%02d%d", rentalID, paymentType, serviceId, startDate.Month(), startDate.Year(), endDate.Month(), endDate.Year())
	}
	return fmt.Sprintf("%d_%s_%02d%d%02d%d", rentalID, paymentType, startDate.Month(), startDate.Year(), endDate.Month(), endDate.Year())
}

func GetServiceName(rpCode string, rServices []model.RentalService) (string, error) {
	parts := strings.Split(rpCode, "_")
	if len(parts) < 3 {
		return "", ErrInvalidRentalPaymentCode
	}
	if parts[1] == string(RENTALPAYMENTTYPESERVICE) {
		if len(parts) != 4 {
			return "", ErrInvalidRentalPaymentCode
		}
		var (
			svcID   string = parts[2]
			svcName string = parts[2]
		)
		for _, rs := range rServices {
			if svcID == fmt.Sprintf("%d", rs.ID) {
				svcName = rs.Name
				break
			}
		}
		return fmt.Sprintf("%s %s", mapRentalPaymentTypeToServiceName[RentalPaymentType(parts[1])], svcName), nil
	}
	return mapRentalPaymentTypeToServiceName[RentalPaymentType(parts[1])], nil
}

func GetRentalPaymentPrice(startDate, endDate time.Time, paymentBasis int32, rentalPrice float64) float64 {
	// calculate rental duration in days
	rentalDuration := endDate.Sub(startDate).Hours() / 24

	// calculate rental price based on number of days
	if startDate.AddDate(0, int(paymentBasis), 0).After(endDate) {
		// prorate rental price
		return rentalPrice * rentalDuration / float64(paymentBasis*30)
	} else {
		// full rental price
		return rentalPrice
	}
}
