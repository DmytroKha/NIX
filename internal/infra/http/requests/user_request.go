package requests

import "NIX/internal/domain"

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,gte=6"`
	Name     string `json:"name"`
}

func (r UserRequest) ToDomainModel() (domain.User, error) {
	return domain.User{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}, nil
}
