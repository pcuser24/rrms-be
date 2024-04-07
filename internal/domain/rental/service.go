package rental

import (
	"context"

	"github.com/google/uuid"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
)

type Service interface {
	CreateRental(data *dto.CreateRental, creatorID uuid.UUID) (*model.RentalModel, error)
	GetRental(id int64) (*model.RentalModel, error)
	UpdateRental(data *dto.UpdateRental, id int64) error
	// PrepareRentalContract(id int64, data *dto.PrepareRentalContract) (*model.RentalContractModel, error)
	CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error)

	CreateContract(data *dto.CreateContract) (*model.ContractModel, error)
	GetRentalContract(id int64) (*model.ContractModel, error)
	PingRentalContract(id int64) (any, error)
	GetContract(id int64) (*model.ContractModel, error)
	UpdateContract(data *dto.UpdateContract) error
	UpdateContractContent(data *dto.UpdateContractContent) error
}

type service struct {
	authRepo auth_repo.Repo
	aRepo    application_repo.Repo
	lRepo    listing_repo.Repo
	pRepo    property_repo.Repo
	rRepo    repo.Repo
	uRepo    unit_repo.Repo
}

func NewService(rRepo repo.Repo, authRepo auth_repo.Repo, aRepo application_repo.Repo, lRepo listing_repo.Repo, pRepo property_repo.Repo, uRepo unit_repo.Repo) Service {
	return &service{
		authRepo: authRepo,
		aRepo:    aRepo,
		rRepo:    rRepo,
		lRepo:    lRepo,
		pRepo:    pRepo,
		uRepo:    uRepo,
	}
}

func (s *service) CreateRental(data *dto.CreateRental, userId uuid.UUID) (*model.RentalModel, error) {
	// TODO: validate applicationId, propertyId, unitId
	data.CreatorID = userId
	return s.rRepo.CreateRental(context.Background(), data)
}

func (s *service) GetRental(id int64) (*model.RentalModel, error) {
	return s.rRepo.GetRental(context.Background(), id)
}

func (s *service) UpdateRental(data *dto.UpdateRental, id int64) error {
	return s.rRepo.UpdateRental(context.Background(), data, id)
}

// func (s *service) UpdateRentalContract(data *dto.UpdateRentalContract, id int64) error {
// 	return s.rRepo.UpdateRentalContract(context.Background(), data, id)
// }

// func (s *service) PrepareRentalContract(id int64, data *dto.PrepareRentalContract) (*model.RentalContractModel, error) {
// 	var (
// 		u   dto.UpdateRentalContract
// 		res model.RentalContractModel
// 	)

// 	pr, err := s.rRepo.GetRental(context.Background(), id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if pr.ContractType != "" && pr.ContractContent != nil {
// 		return &model.RentalContractModel{
// 			ContractType:         pr.ContractType,
// 			ContractContent:      pr.ContractContent,
// 			ContractLastUpdateAt: pr.ContractLastUpdateAt,
// 			ContractLastUpdateBy: pr.ContractLastUpdateBy,
// 		}, nil
// 	}

// 	if data.ContractType != database.CONTRACTTYPEDIGITAL {
// 		u.ContractContent = data.ContractContent
// 		u.ContractType = data.ContractType
// 		res.ContractContent = data.ContractContent
// 		res.ContractType = data.ContractType
// 		_ = s.rRepo.UpdateRentalContract(context.Background(), &u, id)
// 	}

// 	// fetch neccessary data for digital contract
// 	var (
// 		owner *auth_model.UserModel
// 		a     *application_model.ApplicationModel
// 	)
// 	if pr.ApplicationID != nil {
// 		a, err = s.aRepo.GetApplicationById(context.Background(), *pr.ApplicationID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	p, err := s.pRepo.GetPropertyById(context.Background(), pr.PropertyID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	unit, err := s.uRepo.GetUnitById(context.Background(), a.UnitID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, m := range p.Managers {
// 		if m.Role == "OWNER" {
// 			owner, err = s.authRepo.GetUserById(context.Background(), m.ManagerID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			break
// 		}
// 	}

// 	// update contract
// 	con, err := contract.RenderContractTemplate(pr, a, p, unit, owner)
// 	if err != nil {
// 		return nil, err
// 	}
// 	u.ContractContent = &con
// 	u.ContractType = database.CONTRACTTYPEDIGITAL
// 	res.ContractContent = &con
// 	res.ContractType = database.CONTRACTTYPEDIGITAL

// 	_ = s.rRepo.UpdateRentalContract(context.Background(), &u, id)

// 	return &res, nil
// }

func (s *service) CheckRentalVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.rRepo.CheckRentalVisibility(context.Background(), id, userId)
}

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
