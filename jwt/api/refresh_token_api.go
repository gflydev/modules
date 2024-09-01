package api

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/modules/jwt"
	"github.com/gflydev/modules/jwt/dto"
	"github.com/gflydev/modules/jwt/request"
	"github.com/gflydev/modules/jwt/response"
	"github.com/gflydev/validation"
)

// NewRefreshTokenApi As a constructor to create new API.
func NewRefreshTokenApi() *RefreshTokenApi {
	return &RefreshTokenApi{}
}

type RefreshTokenApi struct {
	core.Api
}

// Validate validates request refresh token
func (h *RefreshTokenApi) Validate(c *core.Ctx) error {
	var refreshToken request.RefreshToken
	err := c.ParseBody(&refreshToken)
	if err != nil {
		c.Status(core.StatusBadRequest)
		return err
	}

	refreshTokenDto := refreshToken.ToDto()
	errorData, err := validation.Check(refreshTokenDto)
	if err != nil {
		_ = c.Error(errorData)
		return err
	}

	c.SetData(data, refreshTokenDto)
	return nil
}

// Handle method to refresh user token.
// @Description Refresh user token
// @Summary refresh user token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.RefreshToken true "RefreshToken payload"
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 200 {object} response.SignIn
// @Security ApiKeyAuth
// @Router /auth/refresh [put]
func (h *RefreshTokenApi) Handle(c *core.Ctx) error {
	refreshToken := c.GetData(data).(dto.RefreshToken)
	// Check valid refresh token
	if !jwt.IsValidRefreshToken(refreshToken.Token) {
		return c.Error(errors.New("Invalid JWT token"))
	}

	jwtToken := jwt.ExtractToken(c)
	// Refresh new pairs of access token & refresh token
	tokens, err := jwt.RefreshToken(jwtToken, refreshToken.Token)
	if err != nil {
		return c.Error(errors.New("Error %v", err))
	}

	// Return response.SignIn struct
	return c.JSON(response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	})
}
