package service

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/payment/dto"
	"github.com/user2410/rrms-backend/internal/domain/payment/model"
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
)

type Service interface {
	CreatePayment(data *dto.CreatePayment) (*model.PaymentModel, error)
	GetPaymentById(id int64) (*model.PaymentModel, error)
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreatePayment(data *dto.CreatePayment) (*model.PaymentModel, error) {
	return s.repo.CreatePayment(context.Background(), data)
}

func (s *service) GetPaymentById(id int64) (*model.PaymentModel, error) {
	return s.repo.GetPaymentById(context.Background(), id)
}
