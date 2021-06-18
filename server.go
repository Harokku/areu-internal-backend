package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"internal-backend/crawler"
	"internal-backend/database"
	"internal-backend/router"
	"internal-backend/utils"
	"internal-backend/websocket"
	"log"
	"time"
)

func main() {
	// -------------------------
	// Variable definition
	// -------------------------

	var (
		err  error
		port string //server port from env
	)

	log.Printf("Starting environment init...")
	initStartTime := time.Now() //Startup timer start

	// -------------------------
	// .env loading
	// -------------------------

	// Load YAML with godotenv pkg
	err = godotenv.Load("env.yaml")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read server port from env
	port, err = utils.ReadEnv("PORT")
	if err != nil {
		log.Fatalf("Fatal error setting server port: %v", err)
	}
	log.Printf("Server port set to: %v", port)

	// Read secret from env
	_, err = utils.ReadEnv("SECRET")
	if err != nil {
		log.Fatalf("Fatal error setting secret: %v", err)
	}
	log.Printf("JWT Secret set")

	// -------------------------
	// Database connection
	// -------------------------

	database.Connect()
	defer database.DbConnection.Close()

	crawler.EnumerateDocuments()

	// -------------------------
	// Fiber definition and server start
	// -------------------------

	app := fiberApp()

	// -------------------------
	// Websocket hub start in separate goroutine
	// -------------------------
	go websocket.RunHub()

	initDuration := time.Since(initStartTime) //calculate total startup time
	log.Printf("Enviromnent initialized in %s", initDuration)

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}

// Create a new fiber app and define routes
func fiberApp() *fiber.App {
	var (
		app *fiber.App
	)
	app = fiber.New()
	app.Use(logger.New())  //logger init
	app.Use(cors.New())    //CORS init
	app.Use(recover.New()) //recover init

	// -------------------------
	// Static routes
	// -------------------------

	app.Static("/", "./static")
	app.Static("/intranet", "./static/build")

	// -------------------------
	// Debug routes
	// -------------------------

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	// -------------------------
	// FrontEnd
	// -------------------------

	app.Get("/intranet", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./static/build/index.html")
	})

	// -------------------------
	// Router init (config in router pkg)
	// -------------------------
	router.SetupRoutes(app)

	return app
}
