package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (r *repo) CreatePropertyManagerRequest(ctx context.Context, data *property_dto.CreatePropertyManagerRequest) (property_model.NewPropertyManagerRequest, error) {
	res, err := r.dao.CreateNewPropertyManagerRequest(ctx, database.CreateNewPropertyManagerRequestParams{
		CreatorID:  data.CreatorID,
		PropertyID: data.PropertyID,
		UserID: pgtype.UUID{
			Valid: data.UserID != uuid.Nil,
			Bytes: data.UserID,
		},
		Email: data.Email,
	})
	if err != nil {
		return property_model.NewPropertyManagerRequest{}, err
	}

	return property_model.NewPropertyManagerRequest{
		ID:         res.ID,
		CreatorID:  res.CreatorID,
		PropertyID: res.PropertyID,
		UserID:     res.UserID.Bytes,
		Email:      res.Email,
		Approved:   res.Approved,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, err
}

func (r *repo) UpdatePropertyManagerRequest(ctx context.Context, id int64, uid uuid.UUID, approved bool) error {
	if !approved {
		sb := sqlbuilder.PostgreSQL.NewDeleteBuilder()
		sb.DeleteFrom("new_property_manager_requests")
		sb.Where(sb.Equal("id", id))

		sql, args := sb.Build()
		_, err := r.dao.Exec(ctx, sql, args...)
		return err
	} else {
		err := r.dao.ExecTx(ctx, nil, func(dao database.DAO) error {
			err := dao.UpdateNewPropertyManagerRequest(ctx, database.UpdateNewPropertyManagerRequestParams{
				ID:       id,
				Approved: true,
			})
			if err != nil {
				return err
			}
			return dao.AddPropertyManager(ctx, database.AddPropertyManagerParams{
				RequestID: id,
				UserID:    uid,
			})
		})
		if err == nil {
			return error(nil)
		}
		return err
	}
}

func (r *repo) GetNewPropertyManagerRequest(ctx context.Context, id int64) (property_model.NewPropertyManagerRequest, error) {
	res, err := r.dao.GetNewPropertyManagerRequest(ctx, id)
	if err != nil {
		return property_model.NewPropertyManagerRequest{}, err
	}
	return property_model.NewPropertyManagerRequest{
		ID:         res.ID,
		CreatorID:  res.CreatorID,
		PropertyID: res.PropertyID,
		UserID:     res.UserID.Bytes,
		Email:      res.Email,
		Approved:   res.Approved,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

func (r *repo) GetNewPropertyManagerRequestsToUser(ctx context.Context, uid uuid.UUID, limit, offset int64) ([]property_model.NewPropertyManagerRequest, error) {
	res, err := r.dao.GetNewPropertyManagerRequestsToUser(ctx, database.GetNewPropertyManagerRequestsToUserParams{
		UserID: pgtype.UUID{
			Valid: true,
			Bytes: uid,
		},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	var items []property_model.NewPropertyManagerRequest
	for _, i := range res {
		items = append(items, property_model.NewPropertyManagerRequest{
			ID:         i.ID,
			CreatorID:  i.CreatorID,
			PropertyID: i.PropertyID,
			UserID:     i.UserID.Bytes,
			Email:      i.Email,
			Approved:   i.Approved,
			CreatedAt:  i.CreatedAt,
			UpdatedAt:  i.UpdatedAt,
		})
	}
	return items, nil
}
