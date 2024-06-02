package repo

import (
	"context"

	"github.com/google/uuid"
)

func (r *repo) GetRecentListings(ctx context.Context, limit int32) ([]uuid.UUID, error) {
	return r.dao.GetRecentListings(ctx, limit)
}
