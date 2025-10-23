package api

import (
	"github.com/gflydev/core"
	"github.com/gflydev/http"
	"github.com/gflydev/modules/auth"
	"github.com/gflydev/modules/auth/request"
	_ "github.com/gflydev/modules/auth/response" // Used for Swagger documentation
	"github.com/gflydev/modules/auth/services"
	"github.com/gflydev/modules/auth/transformers"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type SignInApi struct {
	Type auth.Type
	core.Api
}

// NewSignInApi is a constructor
func NewSignInApi(authType auth.Type) *SignInApi {
	return &SignInApi{
		Type: authType,
	}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate data from request
func (h *SignInApi) Validate(c *core.Ctx) error {
	return http.ProcessData[request.SignIn](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle func handle sign in user then returns access token and refresh token
// @Description Authenticating user's credentials then return access and refresh token if valid. Otherwise, return an error message.
// @Summary authenticating user's credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.SignIn true "Signin payload"
// @Success 200 {object} response.SignIn
// @Failure 400 {object} http.Error
// @Router /auth/signin [post]
func (h *SignInApi) Handle(c *core.Ctx) error {
	// Get valid data from context
	requestData := c.GetData(http.RequestKey).(request.SignIn)

	tokens, err := services.SignIn(requestData.ToDto())
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		})
	}

	if h.Type == auth.TypeWeb {
		c.SetSession(auth.SessionUsername, requestData.ToDto().Username)

		return c.NoContent()
	}

	return c.JSON(transformers.ToSignInResponse(tokens))
}
