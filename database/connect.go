package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"internal-backend/utils"
	"log"
)

var DbConnection *sql.DB

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
	log.Printf("DbConnection url set")

	// -------------------------
	// DB pool connection
	// -------------------------
	DbConnection, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}
	log.Printf("Connection string set")

	//defer DbConnection.Close()

	// Try ping db to check for availability
	err = DbConnection.Ping()
	if err != nil {
		log.Fatalf("Can't ping database %v", err)
	}
	log.Printf("DbConnection correctly pinged")

	// Init db if not already done
	sqlstatement := `
		create table if not exists docs
		(
			id           uuid    default gen_random_uuid() not null
				constraint docs_pk
					primary key,
			hash         varchar,
			filename     varchar,
			displayname  varchar,
			category     varchar,
			"isDir"      boolean default false,
			creationtime timestamp
		);
		
		comment on table docs is 'Document indexing table';
		
		comment on column docs.hash is 'SHA-1 file hash';
		
		comment on column docs.filename is 'File name as on disk';
		
		comment on column docs.displayname is 'File name as displayed on screen';
		
		comment on column docs.category is 'Builded vs path Use relative path from doc root to define document category';
		
		comment on column docs."isDir" is 'If the scanned file is a directory';
		
		comment on column docs.creationtime is 'File create timestamp';
		
		create unique index if not exists docs_id_uindex
			on docs (id);

`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Document table: %v", err))
	}
}
