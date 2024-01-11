package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateApplication(ctx context.Context, data *dto.CreateApplicationDto) (*model.ApplicationModel, error)
	GetApplicationById(ctx context.Context, id int64) (*model.ApplicationModel, error)
	GetApplicationsByUserId(ctx context.Context, uid uuid.UUID) ([]model.ApplicationModel, error)
	GetApplicationsToUser(ctx context.Context, uid uuid.UUID) ([]model.ApplicationModel, error)
	UpdateApplicationStatus(ctx context.Context, id int64, status database.APPLICATIONSTATUS) error
	DeleteApplication(ctx context.Context, id int64) error
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateApplication(ctx context.Context, data *dto.CreateApplicationDto) (*model.ApplicationModel, error) {

	res, err := r.dao.QueryTx(ctx, func(d database.DAO) (interface{}, error) {
		var am *model.ApplicationModel

		res, err := d.CreateApplication(ctx, *data.ToCreateApplicationDB())
		if err != nil {
			return nil, err
		}
		am = model.ToApplicationModel(&res)

		for _, m := range data.Minors {
			res, err := d.CreateApplicationMinor(ctx, *m.ToCreateApplicationMinorDB(am.ID))
			if err != nil {
				return nil, err
			}
			am.Minors = append(am.Minors, model.ToApplicationMinorModel(&res))
		}
		for _, c := range data.Coaps {
			res, err := d.CreateApplicationCoap(ctx, *c.ToCreateApplicationCoapDB(am.ID))
			if err != nil {
				return nil, err
			}
			am.Coaps = append(am.Coaps, model.ToApplicationCoapModel(&res))
		}
		for _, p := range data.Pets {
			res, err := d.CreateApplicationPet(ctx, *p.ToCreateApplicationPetDB(am.ID))
			if err != nil {
				return nil, err
			}
			am.Pets = append(am.Pets, model.ToApplicationPetModel(&res))
		}
		for _, v := range data.Vehicles {
			res, err := d.CreateApplicationVehicle(ctx, *v.ToCreateApplicationVehicleDB(am.ID))
			if err != nil {
				return nil, err
			}
			am.Vehicles = append(am.Vehicles, model.ToApplicationVehicleModel(&res))
		}

		return am, nil
	})
	if err != nil {
		return nil, err
	}
	a := res.(*model.ApplicationModel)

	return a, nil
}

func (r *repo) GetApplicationById(ctx context.Context, id int64) (*model.ApplicationModel, error) {
	res, err := r.dao.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a := model.ToApplicationModel(&res)

	applicationMinors, err := r.dao.GetApplicationMinors(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, am := range applicationMinors {
		a.Minors = append(a.Minors, model.ToApplicationMinorModel(&am))
	}

	applicationCoaps, err := r.dao.GetApplicationCoaps(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, ac := range applicationCoaps {
		a.Coaps = append(a.Coaps, model.ToApplicationCoapModel(&ac))
	}

	applicationPets, err := r.dao.GetApplicationPets(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, ap := range applicationPets {
		a.Pets = append(a.Pets, model.ToApplicationPetModel(&ap))
	}

	applicationVehicles, err := r.dao.GetApplicationVehicles(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, av := range applicationVehicles {
		a.Vehicles = append(a.Vehicles, model.ToApplicationVehicleModel(&av))
	}

	return a, nil
}

func (r *repo) GetApplicationsByUserId(ctx context.Context, uid uuid.UUID) ([]model.ApplicationModel, error) {
	res, err := r.dao.GetApplicationsByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}
	var applications []model.ApplicationModel
	for _, a := range res {
		applications = append(applications, *model.ToApplicationModel(&a))
	}
	return applications, nil
}

func (r *repo) GetApplicationsToUser(ctx context.Context, uid uuid.UUID) ([]model.ApplicationModel, error) {
	res, err := r.dao.GetApplicationsToUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	var applications []model.ApplicationModel
	for _, a := range res {
		applications = append(applications, *model.ToApplicationModel(&a))
	}
	return applications, nil
}

func (r *repo) UpdateApplicationStatus(ctx context.Context, id int64, status database.APPLICATIONSTATUS) error {
	return r.dao.UpdateApplicationStatus(ctx, database.UpdateApplicationStatusParams{
		ID:     id,
		Status: status,
	})
}

func (r *repo) DeleteApplication(ctx context.Context, id int64) error {
	return r.dao.DeleteApplication(ctx, id)
}
