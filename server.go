package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	_ "github.com/lib/pq"
	"internal-backend/crawler"
	"internal-backend/database"
	"internal-backend/router"
	"internal-backend/utils"
	"internal-backend/websocket"
	"log"
	"runtime"
	"time"
)

// Call overseer to monitor for file change and program self restart
func main() {
	overseer.Run(overseer.Config{
		Program: prog,
		Fetcher: &fetcher.File{
			Path:     pathByOs(),
			Interval: 5 * time.Second,
		},
	},
	)
}

// Detect current running OS and set overseer path accordingly
func pathByOs() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "./server.exe"
	default:
		return "./server"
	}
}

// Prog is the actual program to be started
func prog(state overseer.State) {
	log.Printf("Overseer state: %t", state.Enabled)
	log.Printf("App %s is running...", state.ID)
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
	// Database connection and initi
	// -------------------------

	database.Connect()
	defer database.DbConnection.Close()

	err = crawler.EnumerateDocuments()
	if err != nil {
		log.Fatalf("Fatal error during file crawler operation: %s", err)
	}

	// -------------------------
	// Fiber definition and server start
	// -------------------------

	app := fiberApp()

	// -------------------------
	// Websocket hub start in separate goroutine
	// -------------------------
	log.Println("Starting websocket hub")
	go websocket.RunHub()

	// -------------------------
	// File change monitor
	// -------------------------
	log.Println("Starting data table watcher...")
	err = crawler.WatchDataTableFromEnv()
	if err != nil {
		log.Fatalf("Error starting data table watcher: %s", err)
	}
	log.Println("Data table watcher initialized")

	log.Println("Starting filewatcher...")
	err = crawler.WatchRootFromEnv()
	if err != nil {
		log.Fatalf("Error starting filewatcher: %s", err)
	}
	log.Println("Filewatcher initialized")

	log.Println("Starting fleet watcher")
	err = crawler.WatchFleetFromEnv()
	if err != nil {
		log.Fatalf("Error starting fleet watcher: %s", err)
	}
	log.Println("Fleet watcher initialized")

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
	app.Get("/restart", func(ctx *fiber.Ctx) error {
		overseer.Restart()
		return ctx.SendString("Service restarting...")
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
