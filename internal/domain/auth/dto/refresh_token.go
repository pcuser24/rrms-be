package dto

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
	AccessToken  string `json:"accessToken" validate:"required"`
}
