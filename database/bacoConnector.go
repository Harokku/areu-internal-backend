package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"internal-backend/utils"
	"log"
)

var BacoDbConnection *sql.DB

func BacoConnect() {
	var (
		err   error
		dbUrl string //database url
	)

	// -------------------------
	// .env loading
	// -------------------------

	// Read dbUrls from env
	dbUrl, err = utils.ReadEnv("DATABASE_URL_BACO")
	if err != nil {
		log.Fatalf("Fatal error setting Baco database url: %v", err)
	}
	log.Printf("Baco DbConnection url set")

	// -------------------------
	// DB pool connection
	// -------------------------
	BacoDbConnection, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Can't connect to Baco db: %v", err)
	}
	log.Printf("Baco Connection string set")

	// Try ping db to check for availability
	err = BacoDbConnection.Ping()
	if err != nil {
		log.Fatalf("Can't ping Baco database %v", err)
	}
	log.Printf("Baco DbConnection correctly pinged")

}
