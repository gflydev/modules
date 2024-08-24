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
		Token:        null.ScanString(user.Token),
		Status:       user.Status,
		Avatar:       null.ScanString(user.Avatar),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		VerifiedAt:   null.ScanTime(user.VerifiedAt),
		BlockedAt:    null.ScanTime(user.BlockedAt),
		DeletedAt:    null.ScanTime(user.DeletedAt),
		LastAccessAt: null.ScanTime(user.LastAccessAt),
	}
}
