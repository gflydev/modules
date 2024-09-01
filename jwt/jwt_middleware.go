package jwt

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/db"
	"github.com/gflydev/modules/jwt/model"
	"slices"
	"time"
)

// New an HTTP middleware that process login via JWT token.
//
// Use:
//
//	app.Use(jwt.New(
//		prefixAPI+"/info",
//		prefixAPI+"/auth/signin",
//		prefixAPI+"/auth/refresh",
//	))
func New(excludes ...string) core.MiddlewareHandler {
	return func(c *core.Ctx) error {
		path := c.Path()
		if slices.Contains(excludes, path) {
			log.Debugf("Skip JWTAuth checking for '%v'", path)

			return nil
		}

		// Forge status code 401 (Unauthorized) instead 500 (internal error)
		c.Root().Response.SetStatusCode(core.StatusUnauthorized)

		jwtToken := ExtractToken(c)
		isBlocked, err := IsBlockedToken(jwtToken)
		if err != nil {
			log.Errorf("Check JWT error '%v'", err)

			return errors.New("JWT token was blocked")
		}

		if isBlocked {
			return errors.New("Invalid JWT token")
		}

		// Get claims from JWT.
		claims, err := ExtractTokenMetadata(jwtToken)
		if err != nil {
			log.Errorf("Parse JWT error '%v'", err)

			return errors.New("Parse JWT error")
		}

		if claims.Expires < time.Now().Unix() {
			log.Errorf("JWT token expired '%v'", jwtToken)

			return errors.New("JWT token expired")
		}

		var user = model.User{}

		// Keep user ID.
		err = db.GetModel(&user, "id", claims.UserID)
		if err != nil {
			log.Errorf("User not found '%v'", err)

			return errors.New("User not found")
		}

		c.Root().Response.SetStatusCode(core.StatusOK)
		c.SetData(User, user)
		return nil
	}
}
