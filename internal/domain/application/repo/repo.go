package repo

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateApplication(ctx context.Context, data *dto.CreateApplication) (*model.ApplicationModel, error)
	GetApplicationById(ctx context.Context, id int64) (*model.ApplicationModel, error)
	GetApplicationsByIds(ctx context.Context, ids []int64, fields []string) ([]model.ApplicationModel, error) // Get applications with custom fields by ids
	GetApplicationsByUserId(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]int64, error)
	GetApplicationsToUser(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]int64, error)
	CheckVisibility(ctx context.Context, id int64, uid uuid.UUID) (bool, error)
	CheckUpdatability(ctx context.Context, id int64, uid uuid.UUID) (bool, error)
	UpdateApplicationStatus(ctx context.Context, aid int64, userId uuid.UUID, status database.APPLICATIONSTATUS) (int, error)
	DeleteApplication(ctx context.Context, id int64) error
	GetRentalByApplicationId(ctx context.Context, aid int64) (rental_model.RentalModel, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateApplication(ctx context.Context, data *dto.CreateApplication) (*model.ApplicationModel, error) {

	var am *model.ApplicationModel

	res, err := r.dao.CreateApplication(ctx, *data.ToCreateApplicationDB())
	if err != nil {
		return nil, err
	}
	am = model.ToApplicationModel(&res)

	err = func() error {
		for _, m := range data.Minors {
			res, err := r.dao.CreateApplicationMinor(ctx, *m.ToCreateApplicationMinorDB(am.ID))
			if err != nil {
				return err
			}
			am.Minors = append(am.Minors, model.ToApplicationMinorModel(&res))
		}
		for _, c := range data.Coaps {
			res, err := r.dao.CreateApplicationCoap(ctx, *c.ToCreateApplicationCoapDB(am.ID))
			if err != nil {
				return err
			}
			am.Coaps = append(am.Coaps, model.ToApplicationCoapModel(&res))
		}
		for _, p := range data.Pets {
			res, err := r.dao.CreateApplicationPet(ctx, *p.ToCreateApplicationPetDB(am.ID))
			if err != nil {
				return err
			}
			am.Pets = append(am.Pets, model.ToApplicationPetModel(&res))
		}
		for _, v := range data.Vehicles {
			res, err := r.dao.CreateApplicationVehicle(ctx, *v.ToCreateApplicationVehicleDB(am.ID))
			if err != nil {
				return err
			}
			am.Vehicles = append(am.Vehicles, model.ToApplicationVehicleModel(&res))
		}

		return nil
	}()

	if err != nil {
		// rollback and ignore any error
		_ = r.dao.DeleteApplication(ctx, am.ID)
		return nil, err
	}

	return am, nil
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

func (r *repo) GetApplicationsByIds(ctx context.Context, ids []int64, fields []string) ([]model.ApplicationModel, error) {

	if len(ids) == 0 {
		return nil, nil
	}
	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"minors", "coaps", "tags", "media"}, f) {
			fkFields = append(fkFields, f)
		} else {
			nonFKFields = append(nonFKFields, f)
		}
	}
	// log.Println(nonFKFields, fkFields)

	// get non fk fields
	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("applications")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	query, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.ApplicationModel
	var i database.Application
	var scanningFields []interface{} = []interface{}{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "creator_id":
			scanningFields = append(scanningFields, &i.CreatorID)
		case "listing_id":
			scanningFields = append(scanningFields, &i.ListingID)
		case "property_id":
			scanningFields = append(scanningFields, &i.PropertyID)
		case "unit_id":
			scanningFields = append(scanningFields, &i.UnitID)
		case "listing_price":
			scanningFields = append(scanningFields, &i.ListingPrice)
		case "offered_price":
			scanningFields = append(scanningFields, &i.OfferedPrice)
		case "status":
			scanningFields = append(scanningFields, &i.Status)
		case "created_at":
			scanningFields = append(scanningFields, &i.CreatedAt)
		case "updated_at":
			scanningFields = append(scanningFields, &i.UpdatedAt)
		case "full_name":
			scanningFields = append(scanningFields, &i.FullName)
		case "email":
			scanningFields = append(scanningFields, &i.Email)
		case "phone":
			scanningFields = append(scanningFields, &i.Phone)
		case "dob":
			scanningFields = append(scanningFields, &i.Dob)
		case "profile_image":
			scanningFields = append(scanningFields, &i.ProfileImage)
		case "movein_date":
			scanningFields = append(scanningFields, &i.MoveinDate)
		case "preferred_term":
			scanningFields = append(scanningFields, &i.PreferredTerm)
		case "rental_intention":
			scanningFields = append(scanningFields, &i.RentalIntention)
		case "rh_address":
			scanningFields = append(scanningFields, &i.RhAddress)
		case "rh_city":
			scanningFields = append(scanningFields, &i.RhCity)
		case "rh_district":
			scanningFields = append(scanningFields, &i.RhDistrict)
		case "rh_ward":
			scanningFields = append(scanningFields, &i.RhWard)
		case "rh_rental_duration":
			scanningFields = append(scanningFields, &i.RhRentalDuration)
		case "rh_monthly_payment":
			scanningFields = append(scanningFields, &i.RhMonthlyPayment)
		case "rh_reason_for_leaving":
			scanningFields = append(scanningFields, &i.RhReasonForLeaving)
		case "employment_status":
			scanningFields = append(scanningFields, &i.EmploymentStatus)
		case "employment_company_name":
			scanningFields = append(scanningFields, &i.EmploymentCompanyName)
		case "employment_position":
			scanningFields = append(scanningFields, &i.EmploymentPosition)
		case "employment_monthly_income":
			scanningFields = append(scanningFields, &i.EmploymentMonthlyIncome)
		case "employment_comment":
			scanningFields = append(scanningFields, &i.EmploymentComment)
		}
	}
	for rows.Next() {
		if err := rows.Scan(scanningFields...); err != nil {
			return nil, err
		}
		items = append(items, *model.ToApplicationModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get fk fields
	for i := 0; i < len(items); i++ {
		p := &items[i]
		if slices.Contains(fkFields, "minors") {
			m, err := r.dao.GetApplicationMinors(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, mdb := range m {
				p.Minors = append(p.Minors, model.ToApplicationMinorModel(&mdb))
			}
		}
		if slices.Contains(fkFields, "coaps") {
			c, err := r.dao.GetApplicationCoaps(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, cdb := range c {
				p.Coaps = append(p.Coaps, model.ToApplicationCoapModel(&cdb))
			}
		}
		if slices.Contains(fkFields, "pets") {
			t, err := r.dao.GetApplicationPets(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, tdb := range t {
				p.Pets = append(p.Pets, model.ToApplicationPetModel(&tdb))
			}
		}
		if slices.Contains(fkFields, "vehicles") {
			v, err := r.dao.GetApplicationVehicles(ctx, p.ID)
			if err != nil {
				return nil, err
			}
			for _, vdb := range v {
				p.Vehicles = append(p.Vehicles, model.ToApplicationVehicleModel(&vdb))
			}
		}

	}
	return items, nil
}

func (r *repo) GetApplicationsByUserId(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]int64, error) {
	return r.dao.GetApplicationsByUserId(ctx, database.GetApplicationsByUserIdParams{
		UserID: pgtype.UUID{
			Bytes: uid,
			Valid: uid != uuid.Nil,
		},
		Limit:  limit,
		Offset: offset,
	})
}

func (r *repo) GetApplicationsToUser(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]int64, error) {
	return r.dao.GetApplicationsToUser(ctx, database.GetApplicationsToUserParams{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *repo) CheckVisibility(ctx context.Context, id int64, uid uuid.UUID) (bool, error) {
	return r.dao.CheckApplicationVisibility(ctx, database.CheckApplicationVisibilityParams{
		ID:        id,
		ManagerID: uid,
	})
}

func (r *repo) CheckUpdatability(ctx context.Context, id int64, uid uuid.UUID) (bool, error) {
	return r.dao.CheckApplicationUpdatabilty(ctx, database.CheckApplicationUpdatabiltyParams{
		ID:        id,
		ManagerID: uid,
	})
}

func (r *repo) UpdateApplicationStatus(ctx context.Context, aid int64, userId uuid.UUID, status database.APPLICATIONSTATUS) (int, error) {
	res, err := r.dao.UpdateApplicationStatus(ctx, database.UpdateApplicationStatusParams{
		ID:        aid,
		Status:    status,
		ManagerID: userId,
	})
	if err != nil {
		return 0, err
	}
	return len(res), nil
}

func (r *repo) DeleteApplication(ctx context.Context, id int64) error {
	return r.dao.DeleteApplication(ctx, id)
}

func (r *repo) GetRentalByApplicationId(ctx context.Context, aid int64) (rental_model.RentalModel, error) {
	res, err := r.dao.GetRentalByApplicationId(ctx, pgtype.Int8{
		Int64: aid,
		Valid: true,
	})
	if err != nil {
		return rental_model.RentalModel{}, err
	}
	return rental_model.ToRentalModel(&res), nil
}
