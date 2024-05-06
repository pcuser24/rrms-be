package service

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (s *service) CreateContract(data *dto.CreateContract) (*model.ContractModel, error) {
	return s.rRepo.CreateContract(context.Background(), data)
}

func (s *service) GetRentalContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByRentalID(context.Background(), id)
}

func (s *service) PingRentalContract(id int64) (any, error) {
	return s.rRepo.PingRentalContract(context.Background(), id)
}

func (s *service) GetContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByID(context.Background(), id)
}

func (s *service) UpdateContract(data *dto.UpdateContract) error {
	return s.rRepo.UpdateContract(context.Background(), data)
}

func (s *service) UpdateContractContent(data *dto.UpdateContractContent) error {
	return s.rRepo.UpdateContractContent(context.Background(), data)
}
