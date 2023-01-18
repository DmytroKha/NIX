package requests

import (
	"nix_education/internal/infra/database"
)

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,gte=6"`
	Name     string `json:"name"`
}

func (r UserRequest) ToDatabaseModel() (database.User, error) {
	return database.User{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}, nil
}
