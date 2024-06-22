package repo

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreatePreRental(ctx context.Context, data *dto.CreatePreRental) (model.PreRental, error) {
	params, err := data.ToCreatePreRentalDB()
	if err != nil {
		return model.PreRental{}, err
	}
	pr, err := r.dao.CreatePreRental(ctx, params)
	if err != nil {
		return model.PreRental{}, err
	}
	return model.ToPreRentalModel(&pr)
}

func (r *repo) CreateRental(ctx context.Context, data *dto.CreateRental) (model.RentalModel, error) {
	prdb, err := r.dao.CreateRental(ctx, data.ToCreateRentalDB())
	if err != nil {
		return model.RentalModel{}, err
	}
	prm := model.ToRentalModel(&prdb)

	err = func() error {
		for _, items := range data.Coaps {
			coapdb, err := r.dao.CreateRentalCoap(ctx, items.ToCreateRentalCoapDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Coaps = append(prm.Coaps, model.ToRentalCoapModel(&coapdb))
		}
		for _, items := range data.Minors {
			minordb, err := r.dao.CreateRentalMinor(ctx, items.ToCreateRentalMinorDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Minors = append(prm.Minors, model.ToRentalMinor(&minordb))
		}
		for _, items := range data.Pets {
			petdb, err := r.dao.CreateRentalPet(ctx, items.ToCreateRentalPetDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Pets = append(prm.Pets, model.ToRentalPet(&petdb))
		}
		for _, items := range data.Services {
			servicedb, err := r.dao.CreateRentalService(ctx, items.ToCreateRentalServiceDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Services = append(prm.Services, model.ToRentalService(&servicedb))
		}
		for _, items := range data.Policies {
			policydb, err := r.dao.CreateRentalPolicy(ctx, items.ToCreateRentalPolicyDB(prdb.ID))
			if err != nil {
				return err
			}
			prm.Policies = append(prm.Policies, model.RentalPolicy(policydb))
		}
		return nil
	}()
	if err != nil {
		_err := r.dao.DeleteRental(ctx, prdb.ID)
		return model.RentalModel{}, errors.Join(err, _err)
	}

	return prm, nil
}

func (r *repo) GetPreRental(ctx context.Context, id int64) (model.PreRental, error) {
	prdb, err := r.dao.GetPreRental(ctx, id)
	if err != nil {
		return model.PreRental{}, err
	}
	return model.ToPreRentalModel(&prdb)
}

func (r *repo) GetPreRentalsToTenant(ctx context.Context, userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]model.PreRental, error) {
	prs, err := r.dao.GetPreRentalsToTenant(ctx, database.GetPreRentalsToTenantParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		return nil, err
	}

	var items []model.PreRental
	for _, pr := range prs {
		prm, err := model.ToPreRentalModel(&pr)
		if err != nil {
			return nil, err
		}
		items = append(items, prm)
	}

	return items, nil
}

func (r *repo) GetManagedPreRentals(ctx context.Context, userId uuid.UUID, query *dto.GetPreRentalsQuery) ([]model.PreRental, error) {
	prs, err := r.dao.GetManagedPreRentals(ctx, database.GetManagedPreRentalsParams{
		UserID: userId,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		return nil, err
	}

	var items []model.PreRental
	for _, pr := range prs {
		prm, err := model.ToPreRentalModel(&pr)
		if err != nil {
			return nil, err
		}
		items = append(items, prm)
	}

	return items, nil
}

func (r *repo) MovePreRentalToRental(ctx context.Context, id int64) (model.RentalModel, error) {
	pr, err := r.dao.GetPreRental(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	params, err := dto.FromPreRentalDBToCreateRental(&pr)
	if err != nil {
		return model.RentalModel{}, err
	}
	rental, err := r.CreateRental(ctx, &params)
	if err != nil {
		return model.RentalModel{}, err
	}
	err = r.dao.DeletePreRental(ctx, id)
	return rental, err
}

func (r *repo) RemovePreRental(ctx context.Context, id int64) error {
	return r.dao.DeletePreRental(ctx, id)
}

func (r *repo) GetRental(ctx context.Context, id int64) (model.RentalModel, error) {
	prdb, err := r.dao.GetRental(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	prm := model.ToRentalModel(&prdb)

	coapdb, err := r.dao.GetRentalCoapsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range coapdb {
		prm.Coaps = append(prm.Coaps, model.ToRentalCoapModel(&item))
	}

	minordb, err := r.dao.GetRentalMinorsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range minordb {
		prm.Minors = append(prm.Minors, model.ToRentalMinor(&item))
	}

	petdb, err := r.dao.GetRentalPetsByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range petdb {
		prm.Pets = append(prm.Pets, model.ToRentalPet(&item))
	}

	servicedb, err := r.dao.GetRentalServicesByRentalID(ctx, id)
	if err != nil {
		return model.RentalModel{}, err
	}
	for _, item := range servicedb {
		prm.Services = append(prm.Services, model.ToRentalService(&item))
	}

	return prm, nil
}

func (r *repo) GetRentalSide(ctx context.Context, id int64, userId uuid.UUID) (string, error) {
	return r.dao.GetRentalSide(ctx, database.GetRentalSideParams{
		ID:     id,
		UserID: userId,
	})
}

func (r *repo) GetRentalsByIds(ctx context.Context, ids []int64, fields []string) ([]model.RentalModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var nonFKFields []string = []string{"id"}
	var fkFields []string
	for _, f := range fields {
		if slices.Contains([]string{"coaps", "minors", "pets", "services", "policies"}, f) {
			fkFields = append(fkFields, f)
		} else {
			nonFKFields = append(nonFKFields, f)
		}
	}

	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("rentals")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	sql, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.RentalModel
	var i database.Rental
	var scanningFields []any = []any{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "creator_id":
			scanningFields = append(scanningFields, &i.CreatorID)
		case "property_id":
			scanningFields = append(scanningFields, &i.PropertyID)
		case "unit_id":
			scanningFields = append(scanningFields, &i.UnitID)
		case "application_id":
			scanningFields = append(scanningFields, &i.ApplicationID)
		case "tenant_id":
			scanningFields = append(scanningFields, &i.TenantID)
		case "profile_image":
			scanningFields = append(scanningFields, &i.ProfileImage)
		case "tenant_type":
			scanningFields = append(scanningFields, &i.TenantType)
		case "tenant_name":
			scanningFields = append(scanningFields, &i.TenantName)
		case "tenant_phone":
			scanningFields = append(scanningFields, &i.TenantPhone)
		case "tenant_email":
			scanningFields = append(scanningFields, &i.TenantEmail)
		case "organization_name":
			scanningFields = append(scanningFields, &i.OrganizationName)
		case "organization_hq_address":
			scanningFields = append(scanningFields, &i.OrganizationHqAddress)
		case "start_date":
			scanningFields = append(scanningFields, &i.StartDate)
		case "movein_date":
			scanningFields = append(scanningFields, &i.MoveinDate)
		case "rental_period":
			scanningFields = append(scanningFields, &i.RentalPeriod)
		case "payment_type":
			scanningFields = append(scanningFields, &i.PaymentType)
		case "rental_price":
			scanningFields = append(scanningFields, &i.RentalPrice)
		case "rental_payment_basis":
			scanningFields = append(scanningFields, &i.RentalPaymentBasis)
		case "rental_intention":
			scanningFields = append(scanningFields, &i.RentalIntention)
		case "grace_period":
			scanningFields = append(scanningFields, &i.GracePeriod)
		case "late_payment_penalty_scheme":
			scanningFields = append(scanningFields, &i.LatePaymentPenaltyScheme)
		case "late_payment_penalty_amount":
			scanningFields = append(scanningFields, &i.LatePaymentPenaltyAmount)
		case "electricity_setup_by":
			scanningFields = append(scanningFields, &i.ElectricitySetupBy)
		case "electricity_payment_type":
			scanningFields = append(scanningFields, &i.ElectricityPaymentType)
		case "electricity_customer_code":
			scanningFields = append(scanningFields, &i.ElectricityCustomerCode)
		case "electricity_provider":
			scanningFields = append(scanningFields, &i.ElectricityProvider)
		case "electricity_price":
			scanningFields = append(scanningFields, &i.ElectricityPrice)
		case "water_setup_by":
			scanningFields = append(scanningFields, &i.WaterSetupBy)
		case "water_payment_type":
			scanningFields = append(scanningFields, &i.WaterPaymentType)
		case "water_customer_code":
			scanningFields = append(scanningFields, &i.WaterCustomerCode)
		case "water_provider":
			scanningFields = append(scanningFields, &i.WaterProvider)
		case "water_price":
			scanningFields = append(scanningFields, &i.WaterPrice)
		case "note":
			scanningFields = append(scanningFields, &i.Note)
		case "status":
			scanningFields = append(scanningFields, &i.Status)
		case "created_at":
			scanningFields = append(scanningFields, &i.CreatedAt)
		case "updated_at":
			scanningFields = append(scanningFields, &i.UpdatedAt)
		}
	}
	for rows.Next() {
		if err := rows.Scan(scanningFields...); err != nil {
			return nil, err
		}
		items = append(items, model.ToRentalModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get fk fields
	for i := 0; i < len(items); i++ {
		item := &items[i]
		if slices.Contains(fkFields, "coaps") {
			u, err := r.dao.GetRentalCoapsByRentalID(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				item.Coaps = append(item.Coaps, model.ToRentalCoapModel(&i))
			}
		}
		if slices.Contains(fkFields, "minors") {
			u, err := r.dao.GetRentalMinorsByRentalID(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				item.Minors = append(item.Minors, model.ToRentalMinor(&i))
			}
		}
		if slices.Contains(fkFields, "pets") {
			u, err := r.dao.GetRentalPetsByRentalID(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				item.Pets = append(item.Pets, model.ToRentalPet(&i))
			}
		}
		if slices.Contains(fkFields, "services") {
			u, err := r.dao.GetRentalServicesByRentalID(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				item.Services = append(item.Services, model.ToRentalService(&i))
			}
		}
		if slices.Contains(fkFields, "policies") {
			u, err := r.dao.GetRentalPoliciesByRentalID(ctx, item.ID)
			if err != nil {
				return nil, err
			}
			for _, i := range u {
				item.Policies = append(item.Policies, model.RentalPolicy(i))
			}
		}
	}
	return items, nil
}

func (r *repo) UpdateRental(ctx context.Context, data *dto.UpdateRental, id int64) error {
	return r.dao.UpdateRental(ctx, data.ToUpdateRentalDB(id))
}

func (r *repo) CheckRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error) {
	return r.dao.CheckRentalVisibility(ctx, database.CheckRentalVisibilityParams{
		ID: id,
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
	})
}

func (r *repo) CheckPreRentalVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error) {
	return r.dao.CheckPreRentalVisibility(ctx, database.CheckPreRentalVisibilityParams{
		ID: id,
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
	})
}

func (r *repo) GetManagedRentals(ctx context.Context, userId uuid.UUID, query *dto.GetRentalsQuery) ([]int64, error) {
	return r.dao.GetManagedRentals(ctx, database.GetManagedRentalsParams{
		UserID:  userId,
		Expired: query.Expired,
		Limit:   *query.Limit,
		Offset:  *query.Offset,
	})
}

func (r *repo) GetMyRentals(ctx context.Context, userId uuid.UUID, query *dto.GetRentalsQuery) ([]int64, error) {
	return r.dao.GetMyRentals(ctx, database.GetMyRentalsParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Expired: query.Expired,
		Limit:   *query.Limit,
		Offset:  *query.Offset,
	})
}
