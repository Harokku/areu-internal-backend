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
		err          error
		dbUrl        string //database url
		sqlstatement string //SQL statement to exec
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

	// -------------------------
	// Init DB if not already done
	// -------------------------

	// Document table
	sqlstatement = `
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

	// Content table
	sqlstatement = `
		create table if not exists content_links
		(
			id           uuid default gen_random_uuid() not null
				constraint content_links_pk
					primary key,
			display_name varchar                        not null,
			link         varchar                        not null,
			sheet_number integer                        not null
		);
		
		comment on table content_links is 'Contain content links to serve to frontend.
				Auto created reading from config XLSX';
		
		comment on column content_links.display_name is 'Link display name as read from XLSX sheet name';
		
		comment on column content_links.link is 'Sanitized for URL safety.
				Calculated from display_name';
		
		comment on column content_links.sheet_number is 'XLSX sheet number to read for link content';
		
		create unique index if not exists content_links_id_uindex
			on content_links (id);
`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Content table: %v", err))
	}

	// Fleet table
	sqlstatement = `
		create table if not exists check_convenzioni
		(
			id          uuid default gen_random_uuid() not null
				constraint check_convenzioni_pk
					primary key,
			conv_type   varchar                        not null,
			name        varchar                        not null,
			active_from timestamp                      not null
		);
		
		comment on table check_convenzioni is 'Tabella appoggio per sistema controllo convenzioni';
		
		comment on column check_convenzioni.conv_type is 'Tipo di convenzione assegnata';
		
		comment on column check_convenzioni.name is 'Acronimo ente';
		
		comment on column check_convenzioni.active_from is 'Fascia oraria di controllo';
		
		create unique index if not exists check_convenzioni_id_uindex
			on check_convenzioni (id);
`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Content table: %v", err))
	}
}
