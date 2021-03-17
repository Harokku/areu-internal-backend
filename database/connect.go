package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"internal-backend/utils"
	"log"
)

var DB *sql.DB

func Connect() {
	var (
		err   error
		dbUrl string //database url
	)

	// -------------------------
	// .env loading
	// -------------------------

	// Read dbUrls from env
	dbUrl, err = utils.ReadEnv("DATABASE_URL")
	if err != nil {
		log.Fatalf("Fatal error setting database url: %v", err)
	}
	log.Printf("DB url set")

	// -------------------------
	// DB pool connection
	// -------------------------
	DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}
	log.Printf("Connection string set")

	defer DB.Close()

	// Try ping db to check for availability
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Can't ping database %v", err)
	}
	log.Printf("DB correctly pinged")
}
