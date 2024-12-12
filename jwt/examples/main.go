package main

import (
	"fmt"
	_ "github.com/gflydev/cache/redis"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	_ "github.com/gflydev/db/psql"
	"github.com/gflydev/middleware/cors"
	"github.com/gflydev/modules/jwt"
	"github.com/gflydev/modules/jwt/api"
	"github.com/gflydev/view/pongo"
)

// =========================================================================================
//                                     Default API
// =========================================================================================

// NewDefaultApi As a constructor to create new API.
func NewDefaultApi() *DefaultApi {
	return &DefaultApi{}
}

// DefaultApi API struct.
type DefaultApi struct {
	core.Api
}

func (h *DefaultApi) Handle(c *core.Ctx) error {
	return c.JSON(core.Data{
		"name":   core.AppName,
		"server": core.AppURL,
	})
}

// =========================================================================================
//                                     Home page
// =========================================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	core.Page
}

func (m *HomePage) Handle(c *core.Ctx) error {
	return c.View("home", core.Data{
		"title": "gFly | Laravel inspired web framework written in Go",
	})
}

// =========================================================================================
//                                     Routers
// =========================================================================================

func router(g core.IFlyRouter) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	/* ============================ Authentication Middleware ============================*/

	// API Routers
	g.Group(prefixAPI, func(apiRouter *core.Group) {
		apiRouter.GET("/info", NewDefaultApi())

		apiRouter.Use(jwt.New(
			prefixAPI+"/auth/signin",
			prefixAPI+"/auth/signup",
			prefixAPI+"/auth/refresh",
		))

		/* ============================ Authentication Group ===================================*/
		apiRouter.Group("/auth", func(authGroup *core.Group) {
			authGroup.POST("/signin", api.NewSignInApi())
			authGroup.DELETE("/signout", api.NewSignOutApi())
			authGroup.POST("/signup", api.NewSignUpApi())
			authGroup.PUT("/refresh", api.NewRefreshTokenApi())
		})
	})

	// Web Routers
	g.GET("/home", NewHomePage())
}

// =========================================================================================
//                                     Application
// =========================================================================================

func main() {
	app := core.New()

	// Middleware
	app.Use(cors.New(cors.Data{
		core.HeaderAccessControlAllowOrigin: cors.AllowedOrigin,
	}))

	// Register router
	app.RegisterRouter(router)

	// Load Model builder
	mb.Load()

	// Register view
	core.RegisterView(pongo.New())

	app.Run()
}
