package model

import (
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type ListingPolicyModel struct {
	ListingID uuid.UUID `json:"listingID"`
	PolicyID  int64     `json:"policyID"`
	Note      *string   `json:"note"`
}

func ToListingPolicyModel(lp *database.ListingPolicy) *ListingPolicyModel {
	lm := &ListingPolicyModel{
		ListingID: lp.ListingID,
		PolicyID:  lp.PolicyID,
	}

	if lp.Note.Valid {
		val := lp.Note.String
		lm.Note = &val
	}

	return lm
}
