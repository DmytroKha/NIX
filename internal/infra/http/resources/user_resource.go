package resources

import (
	"NIX/internal/domain"
	jwt "github.com/dgrijalva/jwt-go"
)

type UserDto struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UsersDto struct {
	Items []UserDto `json:"items"`
	Total uint64    `json:"total"`
	Pages uint64    `json:"pages"`
}

type AuthDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type GoogleUrlDto struct {
	Url string `json:"url"`
}

func (d UserDto) DomainToDto(user domain.User) UserDto {
	return UserDto{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}

func (d AuthDto) DomainToDto(token string, user domain.User) AuthDto {
	var userDto UserDto

	a := AuthDto{
		Token: token,
		User:  userDto.DomainToDto(user),
	}
	return a
}

func (d GoogleUrlDto) DomainToDto(url string) GoogleUrlDto {
	return GoogleUrlDto{
		Url: url,
	}
}

func (d UserDto) DomainToDtoCollection(users domain.Users) UsersDto {
	result := make([]UserDto, len(users.Items))

	for i := range users.Items {
		result[i] = d.DomainToDto(users.Items[i])
	}

	return UsersDto{Items: result, Pages: users.Pages, Total: users.Total}
}
