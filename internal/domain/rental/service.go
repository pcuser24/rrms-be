package rental

import (
	"context"

	"github.com/google/uuid"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/domain/rental/contract"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Service interface {
	CreatePreRental(data *dto.CreatePreRental, creatorID uuid.UUID) (*model.PrerentalModel, error)
	GetPreRental(id int64) (*model.PrerentalModel, error)
	GetPreRentalContract(id int64) (*model.PreRentalContractModel, error)
	UpdatePreRental(data *dto.UpdatePreRental, id int64) error
	UpdatePreRentalContract(data *dto.UpdatePreRentalContract, id int64) error
	PrepareRentalContract(id int64, data *dto.PreparePreRentalContract) (*model.PreRentalContractModel, error)
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

func (s *service) CreatePreRental(data *dto.CreatePreRental, creatorID uuid.UUID) (*model.PrerentalModel, error) {
	data.CreatorID = creatorID
	return s.rRepo.CreatePreRental(context.Background(), data)
}

func (s *service) GetPreRental(id int64) (*model.PrerentalModel, error) {
	return s.rRepo.GetPreRental(context.Background(), id)
}

func (s *service) GetPreRentalContract(id int64) (*model.PreRentalContractModel, error) {
	c, err := s.rRepo.GetPreRentalContract(context.Background(), id)
	if err != nil {
		return c, err
	}
	if c.ContractType == "" || c.ContractContent == nil {
		return nil, database.ErrRecordNotFound
	}
	return c, nil
}

func (s *service) UpdatePreRental(data *dto.UpdatePreRental, id int64) error {
	return s.rRepo.UpdatePreRental(context.Background(), data, id)
}

func (s *service) UpdatePreRentalContract(data *dto.UpdatePreRentalContract, id int64) error {
	return s.rRepo.UpdatePreRentalContract(context.Background(), data, id)
}

func (s *service) PrepareRentalContract(id int64, data *dto.PreparePreRentalContract) (*model.PreRentalContractModel, error) {
	var (
		u   dto.UpdatePreRentalContract
		res model.PreRentalContractModel
	)

	pr, err := s.rRepo.GetPreRental(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if pr.ContractType != "" && pr.ContractContent != nil {
		return &model.PreRentalContractModel{
			ContractType:         pr.ContractType,
			ContractContent:      pr.ContractContent,
			ContractLastUpdateAt: pr.ContractLastUpdateAt,
			ContractLastUpdateBy: pr.ContractLastUpdateBy,
		}, nil
	}

	if data.ContractType != database.CONTRACTTYPEDIGITAL {
		u.ContractContent = data.ContractContent
		u.ContractType = data.ContractType
		res.ContractContent = data.ContractContent
		res.ContractType = data.ContractType
		_ = s.rRepo.UpdatePreRentalContract(context.Background(), &u, id)
	}

	// fetch neccessary data for digital contract
	var (
		owner *auth_model.UserModel
		a     *application_model.ApplicationModel
	)
	if pr.ApplicationID != nil {
		a, err = s.aRepo.GetApplicationById(context.Background(), *pr.ApplicationID)
		if err != nil {
			return nil, err
		}
	}

	p, err := s.pRepo.GetPropertyById(context.Background(), pr.PropertyID)
	if err != nil {
		return nil, err
	}

	units := make([]unit_model.UnitModel, 0, len(a.Units))
	for _, u := range a.Units {
		unit, err := s.uRepo.GetUnitById(context.Background(), u.UnitID)
		if err != nil {
			return nil, err
		}
		units = append(units, *unit)
	}

	for _, m := range p.Managers {
		if m.Role == "OWNER" {
			owner, err = s.authRepo.GetUserById(context.Background(), m.ManagerID)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	// update contract
	con, err := contract.RenderContractTemplate(pr, a, p, units, owner)
	if err != nil {
		return nil, err
	}
	u.ContractContent = &con
	u.ContractType = database.CONTRACTTYPEDIGITAL
	res.ContractContent = &con
	res.ContractType = database.CONTRACTTYPEDIGITAL

	_ = s.rRepo.UpdatePreRentalContract(context.Background(), &u, id)

	return &res, nil
}
