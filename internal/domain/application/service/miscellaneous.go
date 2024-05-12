package service

import (
	"context"

	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (s *service) GetRentalByApplicationId(aid int64) (rental_model.RentalModel, error) {
	return s.aRepo.GetRentalByApplicationId(context.Background(), aid)
}
