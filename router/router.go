package router

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/handler"
)

func SetupRoutes(app *fiber.App) {
	// -------------------------
	// Grouping and versioning
	// -------------------------

	api := app.Group("/api")

	v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("Version", "v1")
		return ctx.Next()
	})

	// -------------------------
	// Versions landing
	// -------------------------

	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("API version 1 root")
	})

	// -------------------------
	// Auth
	// -------------------------

	auth := v1.Group("/auth")
	auth.Post("/login", handler.Login)

	// -------------------------
	// Documents
	// -------------------------
	docs := v1.Group("/docs")
	docs.Get("/", handler.Docs{}.GetAll)
	docs.Get("/:id", handler.Docs{}.GetById)
	docs.Get("/serveById/:id", handler.Docs{}.ServeById)
}
