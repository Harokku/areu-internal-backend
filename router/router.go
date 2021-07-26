package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"internal-backend/handler"
	websocket2 "internal-backend/websocket"
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

	// -------------------------
	// Content
	// -------------------------
	content := v1.Group("/content")
	content.Get("/", handler.Content{}.GetAll)
	content.Get("/:link", handler.Content{}.GetContent)

	// -------------------------
	// Shifts
	// -------------------------
	shift := v1.Group("/shift")
	shift.Get("/serveByPath/:name/:type", handler.Shift{}.ServeByPath)

	// -------------------------
	// Websocket endpoints
	// -------------------------
	ws := v1.Group("/ws")
	ws.Use(func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("remoteIp", ctx.IP())
			return ctx.Next()
		}
		return ctx.SendStatus(fiber.StatusUpgradeRequired)
	})
	ws.Get("/", websocket2.DocsUpdate())
}
