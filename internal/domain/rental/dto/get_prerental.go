package dto

type GetPreRentalQuery struct {
	ID            int64  `json:"id" validate:"omitempty"`
	ApplicationID *int64 `json:"applicationId" validate:"omitempty"`
}
