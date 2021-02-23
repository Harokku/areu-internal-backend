package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Variable def
	var (
		err error
	)

	// Load .env file
	log.Println("Loading .env file")
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Successfully loaded .env file")

	// Fiber definition
	app := fiberApp()

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))

}

func fiberApp() *fiber.App {
	var (
		// err error
		app *fiber.App
	)
	app = fiber.New()

	// Static route
	app.Static("/", "./static")

	// Debug routes
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	return app
}
