package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"internal-backend/router"
	"internal-backend/utils"
	"log"
	"time"
)

func main() {
	// -------------------------
	// Variable definition
	// -------------------------

	var (
		err   error
		port  string //server port from env
		dbUrl string //database url
		conn  *sql.DB
		_     string //JWT secret
	)

	log.Printf("Starting environment init...")
	initStartTime := time.Now() //Startup timer start

	// -------------------------
	// .env file loading
	// -------------------------

	// Read server port from env
	port, err = utils.ReadEnv("PORT")
	if err != nil {
		log.Fatalf("Fatal error setting server port: %v", err)
	}
	log.Printf("Server port set to: %v", port)

	// Read db url
	dbUrl, err = utils.ReadEnv("DATABASE_URL")
	if err != nil {
		log.Fatalf("Fatal error setting database url: %v", err)
	}
	log.Printf("DB url set")

	// Read secret from env
	_, err = utils.ReadEnv("SECRET")
	if err != nil {
		log.Fatalf("Fatal error setting secret: %v", err)
	}
	log.Printf("JWT Secret set")

	// -------------------------
	// DB pool connection
	// -------------------------
	conn, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}
	log.Printf("Connection string set")

	defer conn.Close()

	// Try ping db to check for availability
	err = conn.Ping()
	if err != nil {
		log.Fatalf("Can't ping database %v", err)
	}
	log.Printf("DB correctly pinged")

	// -------------------------
	// Fiber definition and server start
	// -------------------------

	app := fiberApp()

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

	// -------------------------
	// Debug routes
	// -------------------------

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	// -------------------------
	// Router init (config in router pkg)
	// -------------------------
	router.SetupRoutes(app)

	return app
}
