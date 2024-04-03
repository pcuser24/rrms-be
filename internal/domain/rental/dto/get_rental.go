package dto

type GetRentalQuery struct {
	ID            int64  `json:"id" validate:"omitempty"`
	ApplicationID *int64 `json:"applicationId" validate:"omitempty"`
}
