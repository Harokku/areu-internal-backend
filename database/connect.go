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
			id            uuid    default gen_random_uuid()        not null
				constraint check_convenzioni_pk
					primary key,
			convenzione   varchar                                  not null,
			ente          varchar                                  not null,
			active_from   timestamp                                not null,
			stazionamento varchar                                  not null,
			minimum       varchar default 'N/A'::character varying not null,
			active_days   varchar default 'L-Ma-Me-G-V-S-D'::character varying not null
		);
		
		comment on table check_convenzioni is 'Tabella appoggio per sistema controllo convenzioni';
		
		comment on column check_convenzioni.convenzione is 'Tipo di convenzione assegnata';
		
		comment on column check_convenzioni.ente is 'Acronimo ente';
		
		comment on column check_convenzioni.active_from is 'Fascia oraria di controllo';
		
		comment on column check_convenzioni.stazionamento is 'Luogo di stazionamento';
		
		comment on column check_convenzioni.minimum is 'Minimum number of personnel on board';

		comment on column public.check_convenzioni.active_days is 'Giorni di disponibilita';
		
		create unique index if not exists check_convenzioni_id_uindex
			on check_convenzioni (id);
`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Check convenzioni table: %v", err))
	}

	// Baco DB Table
	sqlstatement = `
		create table if not exists "db_Baco"
		(
			id            uuid default gen_random_uuid() not null
				constraint db_baco_pk
					primary key,
			ente          varchar,
			mezzo         varchar,
			stazionamento varchar,
			radio         varchar,
			convenzione   varchar
		);
		
		comment on table "db_Baco" is 'Baco DB export';
		
		comment on column "db_Baco".ente is 'Ente di appartenenza';
		
		comment on column "db_Baco".mezzo is 'sigra radio completa di eventuale lotto';
		
		comment on column "db_Baco".stazionamento is 'Luogo di stazionamento (join da linked table)';
		
		comment on column "db_Baco".radio is 'codica radio';
		
		comment on column "db_Baco".convenzione is 'Tipo di convenzione';
		
		create unique index if not exists db_baco_id_uindex
			on "db_Baco" (id);
`

	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Baco DB Table: %v", err))
	}

	// Issue management tables
	sqlstatement = `
		create table if not exists public.issue
		(
			id        uuid      default gen_random_uuid()      not null
				primary key,
			timestamp timestamp default now()                  not null,
			operator  varchar                                  not null,
			priority  integer   default 2                      not null,
			note      varchar                                  not null,
			title     varchar   default '-'::character varying not null,
			open      boolean   default true                   not null
		);
		
		comment on table public.issue is 'Issue Tracker';
		
		comment on column public.issue.operator is 'operator name';
		
		comment on column public.issue.note is 'Issue note';
		
		comment on column public.issue.title is 'Issue Title';
		
		comment on column public.issue.open is 'If issue is still open False to achieve';
`
	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Baco DB Table: %v", err))
	}

	// Issue detail
	sqlstatement = `
		create table if not exists issue_detail
		(
			id        uuid      default gen_random_uuid() not null
				primary key,
			issue_id  uuid
				constraint issue_detail_issue_id_fk
					references issue,
			timestamp timestamp default now(),
			operator  varchar                             not null,
			note      varchar                             not null
		);
`

	_, err = DbConnection.Exec(sqlstatement)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating Baco DB Table: %v", err))
	}
}
