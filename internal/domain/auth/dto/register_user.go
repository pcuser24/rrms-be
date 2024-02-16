package dto

type RegisterUser struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=32"`
	LastName  string `json:"lastName" validate:"required,min=1,max=32"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
}
