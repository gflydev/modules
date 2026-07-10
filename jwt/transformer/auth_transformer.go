package transformer

import (
	"github.com/gflydev/db/null"
	"github.com/gflydev/modules/jwt"
	"github.com/gflydev/modules/jwt/model"
	"github.com/gflydev/modules/jwt/response"
)

// ToSignInResponse function JWTTokens struct to SignIn response object.
func ToSignInResponse(tokens *jwt.Tokens) response.SignIn {
	return response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}
}

// ToSignUpResponse function convert User model to User response object.
func ToSignUpResponse(user *model.User) response.User {
	return response.User{
		ID:           user.ID,
		Email:        user.Email,
		Fullname:     user.Fullname,
		Phone:        user.Phone,
		Token:        null.StringNil(user.Token),
		Status:       user.Status,
		Avatar:       null.StringNil(user.Avatar),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		VerifiedAt:   null.TimeNil(user.VerifiedAt),
		BlockedAt:    null.TimeNil(user.BlockedAt),
		DeletedAt:    null.TimeNil(user.DeletedAt),
		LastAccessAt: null.TimeNil(user.LastAccessAt),
	}
}
