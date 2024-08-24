package api

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/modules/jwt"
	"github.com/gflydev/modules/jwt/dto"
	"github.com/gflydev/modules/jwt/request"
	"github.com/gflydev/modules/jwt/transformer"
	"github.com/gflydev/validation"
)

type SignInApi struct {
	core.Api
}

// NewSignInApi is a constructor
func NewSignInApi() *SignInApi {
	return &SignInApi{}
}

// Validate data from request
func (h *SignInApi) Validate(c *core.Ctx) error {
	// Parse login form
	var signIn request.SignIn
	err := c.ParseBody(&signIn)
	if err != nil {
		c.Status(core.StatusBadRequest)
		return err
	}

	signInDto := signIn.ToDto()
	// Validate login form.
	errorData, err := validation.Check(signInDto)
	if err != nil {
		_ = c.Error(errorData)
		return err
	}

	c.SetData(data, signInDto)
	return nil
}

// Handle func handle sign in user then returns access token and refresh token
// @Description Authenticating user's credentials then return access and refresh token if valid. Otherwise, return an error message.
// @Summary authenticating user's credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.SignIn true "Signin payload"
// @Success 200 {object} response.SignIn
// @Failure 400 {object} response.Error
// @Router /auth/signin [post]
func (h *SignInApi) Handle(c *core.Ctx) error {
	// Get valid data from context
	signInDto := c.GetData(data).(dto.SignIn)

	tokens, err := jwt.SignIn(&signInDto)
	if err != nil {
		return c.Error(errors.New("Error %v", err))
	}

	// Transform to response object.
	signInResponse := transformer.ToSignInResponse(tokens)

	return c.JSON(signInResponse)
}
