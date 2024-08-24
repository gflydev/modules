package request

import (
	"github.com/gflydev/modules/jwt/dto"
	"strings"
)

type SignIn struct {
	dto.SignIn
}

type SignUp struct {
	dto.SignUp
}

// RefreshToken struct to refresh JWT token.
type RefreshToken struct {
	dto.RefreshToken
}

// ToDto convert to SignIn DTO object.
func (r SignIn) ToDto() dto.SignIn {
	return dto.SignIn{
		Username: r.Username,
		Password: r.Password,
	}
}

func (r RefreshToken) ToDto() dto.RefreshToken {
	return dto.RefreshToken{
		Token: r.Token,
	}
}

// ToDto Convert to SignUp DTO object.
func (r SignUp) ToDto() dto.SignUp {
	return dto.SignUp{
		Email:    strings.ToLower(r.Email),
		Password: r.Password,
		Fullname: r.Fullname,
		Phone:    r.Phone,
		Status:   r.Status,
	}
}
