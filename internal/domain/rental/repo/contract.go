package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreateContract(ctx context.Context, data *dto.CreateContract) (*model.ContractModel, error) {
	prdb, err := r.dao.CreateContract(ctx, data.ToCreateContractDB())
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) GetRentalContractsOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRentalContracts) ([]int64, error) {
	return r.dao.GetRentalContractsOfUser(ctx, database.GetRentalContractsOfUserParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: userId != uuid.Nil,
		},
		Limit:  *query.Limit,
		Offset: *query.Offset,
	})
}

func (r *repo) GetContractsByIds(ctx context.Context, ids []int64, fields []string) ([]model.ContractModel, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var nonFKFields []string = []string{"id"}
	nonFKFields = append(nonFKFields, fields...)

	ib := sqlbuilder.PostgreSQL.NewSelectBuilder()
	ib.Select(nonFKFields...)
	ib.From("contracts")
	ib.Where(ib.In("id", sqlbuilder.List(ids)))
	sql, args := ib.Build()
	// log.Println(query, args)
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.ContractModel
	var i database.Contract
	var scanningFields []any = []any{&i.ID}
	for _, f := range nonFKFields {
		switch f {
		case "rental_id":
			scanningFields = append(scanningFields, &i.RentalID)
		case "a_fullname":
			scanningFields = append(scanningFields, &i.AFullname)
		case "a_dob":
			scanningFields = append(scanningFields, &i.ADob)
		case "a_phone":
			scanningFields = append(scanningFields, &i.APhone)
		case "a_address":
			scanningFields = append(scanningFields, &i.AAddress)
		case "a_household_registration":
			scanningFields = append(scanningFields, &i.AHouseholdRegistration)
		case "a_identity":
			scanningFields = append(scanningFields, &i.AIdentity)
		case "a_identity_issued_by":
			scanningFields = append(scanningFields, &i.AIdentityIssuedBy)
		case "a_identity_issued_at":
			scanningFields = append(scanningFields, &i.AIdentityIssuedAt)
		case "a_documents":
			scanningFields = append(scanningFields, &i.ADocuments)
		case "a_bank_account":
			scanningFields = append(scanningFields, &i.ABankAccount)
		case "a_bank":
			scanningFields = append(scanningFields, &i.ABank)
		case "a_registration_number":
			scanningFields = append(scanningFields, &i.ARegistrationNumber)
		case "b_fullname":
			scanningFields = append(scanningFields, &i.BFullname)
		case "b_organization_name":
			scanningFields = append(scanningFields, &i.BOrganizationName)
		case "b_organization_hq_address":
			scanningFields = append(scanningFields, &i.BOrganizationHqAddress)
		case "b_organization_code":
			scanningFields = append(scanningFields, &i.BOrganizationCode)
		case "b_organization_code_issued_at":
			scanningFields = append(scanningFields, &i.BOrganizationCodeIssuedAt)
		case "b_organization_code_issued_by":
			scanningFields = append(scanningFields, &i.BOrganizationCodeIssuedBy)
		case "b_dob":
			scanningFields = append(scanningFields, &i.BDob)
		case "b_phone":
			scanningFields = append(scanningFields, &i.BPhone)
		case "b_address":
			scanningFields = append(scanningFields, &i.BAddress)
		case "b_household_registration":
			scanningFields = append(scanningFields, &i.BHouseholdRegistration)
		case "b_identity":
			scanningFields = append(scanningFields, &i.BIdentity)
		case "b_identity_issued_by":
			scanningFields = append(scanningFields, &i.BIdentityIssuedBy)
		case "b_identity_issued_at":
			scanningFields = append(scanningFields, &i.BIdentityIssuedAt)
		case "b_bank_account":
			scanningFields = append(scanningFields, &i.BBankAccount)
		case "b_bank":
			scanningFields = append(scanningFields, &i.BBank)
		case "b_tax_code":
			scanningFields = append(scanningFields, &i.BTaxCode)
		case "payment_method":
			scanningFields = append(scanningFields, &i.PaymentMethod)
		case "n_copies":
			scanningFields = append(scanningFields, &i.NCopies)
		case "created_at_place":
			scanningFields = append(scanningFields, &i.CreatedAtPlace)
		case "content":
			scanningFields = append(scanningFields, &i.Content)
		case "status":
			scanningFields = append(scanningFields, &i.Status)
		case "created_at":
			scanningFields = append(scanningFields, &i.CreatedAt)
		case "updated_at":
			scanningFields = append(scanningFields, &i.UpdatedAt)
		case "created_by":
			scanningFields = append(scanningFields, &i.CreatedBy)
		case "updated_by":
			scanningFields = append(scanningFields, &i.UpdatedBy)
		}
	}
	for rows.Next() {
		if err := rows.Scan(scanningFields...); err != nil {
			return nil, err
		}
		items = append(items, *model.ToContractModel(&i))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repo) GetContractByRentalID(ctx context.Context, id int64) (*model.ContractModel, error) {
	prdb, err := r.dao.GetContractByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) PingRentalContract(ctx context.Context, id int64) (any, error) {
	res, err := r.dao.PingContractByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	return struct {
		ID        int64                   `json:"id"`
		RentalID  int64                   `json:"rentalId"`
		Status    database.CONTRACTSTATUS `json:"status"`
		UpdatedBy uuid.UUID               `json:"updatedBy"`
		UpdatedAt time.Time               `json:"updatedAt"`
	}{
		ID:        res.ID,
		RentalID:  res.RentalID,
		Status:    res.Status,
		UpdatedBy: res.UpdatedBy,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *repo) GetContractByID(ctx context.Context, id int64) (*model.ContractModel, error) {
	prdb, err := r.dao.GetContractByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ToContractModel(&prdb), nil
}

func (r *repo) UpdateContract(ctx context.Context, data *dto.UpdateContract) error {
	return r.dao.UpdateContract(ctx, data.ToUpdateContractDB())
}

func (r *repo) UpdateContractContent(ctx context.Context, data *dto.UpdateContractContent) error {
	return r.dao.UpdateContractContent(ctx, data.ToUpdateContractContentDB())
}
