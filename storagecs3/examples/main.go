package main

import (
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/middleware/cors"
	"github.com/gflydev/modules/storagecs3/api"
	"github.com/gflydev/view/pongo"
	_ "github.com/joho/godotenv/autoload"
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

func router(g core.IFly) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	/* ============================ Authentication Middleware ============================*/

	// API Routers
	g.Group(prefixAPI, func(apiRouter *core.Group) {
		apiRouter.GET("/info", NewDefaultApi())

		/* ============================ Storage Group ==========================================*/
		apiRouter.Group("/storage", func(uploadGroup *core.Group) {
			uploadGroup.GET("/presigned-url", api.NewPresignedURLApi())      // Get presigned URL
			uploadGroup.PUT("/legitimize-files", api.NewLegitimizeFileApi()) // Legitimize uploaded file
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

	// Register view
	core.RegisterView(pongo.New())

	app.Run()
}
