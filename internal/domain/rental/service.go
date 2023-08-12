package rental

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

type Service interface {
	GetAllRentalPolicies() ([]model.RentalPolicyModel, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAllRentalPolicies() ([]model.RentalPolicyModel, error) {
	return s.repo.GetAllRentalPolicies(context.Background())
}
