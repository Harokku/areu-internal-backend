package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"internal-backend/utils"
	"strings"
	"time"
)

type Document struct {
	Id           string    `json:"id"`            //Document UUID
	Hash         string    `json:"hash"`          //Document hash
	FileName     string    `json:"file_name"`     //Document filename (full path)
	DisplayName  string    `json:"display_name"`  //Document displayed name
	Category     string    `json:"category"`      //Document category (based on path)
	IsDir        bool      `json:"is_dir"`        //True if file is a directory
	CreationTime time.Time `json:"creation_date"` //Document creation timestamp
}

// GetAll Get all documents
func (d Document) GetAll(dest *[]Document) error {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category,"isDir",creationtime from docs`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error retrieving documents: %v\n", err))
	}

	defer rows.Close()

	for rows.Next() {
		var d Document
		err = rows.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category, &d.IsDir, &d.CreationTime)
		if err != nil {
			return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
		}
		*dest = append(*dest, d)
	}

	return nil
}

// GetByHash Get document by hash
func (d *Document) GetByHash(hash string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category,"isDir",creationtime from docs where hash=$1`

	row = DbConnection.QueryRow(sqlStatement, hash)
	switch err = row.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category, &d.IsDir, &d.CreationTime); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving doc from db: %v\n", err))
	}
}

// GetById Get document by id
func (d *Document) GetById(id string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select id,hash,filename,displayname,category,"isDir",creationtime from docs where id=$1`

	row = DbConnection.QueryRow(sqlStatement, id)
	switch err = row.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category, &d.IsDir, &d.CreationTime); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving doc from db: %v\n", err))
	}
}

// GetRecent Get most recent {num} documents with {mode} aggregation (all, by main category...)
func (d Document) GetRecent(num int, mode string, dest *[]Document) error {
	var (
		err                   error
		rows                  *sql.Rows
		categories            []string
		sqlStatement          string
		sqlFilteredCategories string
		sqlDistinctCategories string
	)

	sqlStatement = `select id,hash,filename,displayname,category,"isDir",creationtime
					from docs
					where "isDir" = false
					order by creationtime desc
					limit $1`

	sqlFilteredCategories = `select id,hash,filename,displayname,category,"isDir",creationtime
							from docs
							where "isDir" = false and category like $2
							order by creationtime desc
							limit $1`

	sqlDistinctCategories = `select distinct category
							 from docs
							 where "isDir" = true`

	switch mode {
	case "split":
		//TODO: implement

		//Get distinct categories from db
		rows, err = DbConnection.Query(sqlDistinctCategories)
		if err != nil {
			return errors.New(fmt.Sprintf("Error retrieving documents: %v\n", err))
		}

		defer rows.Close()

		// Cycle rows and populate root categories array without duplicates
		for rows.Next() {
			var c string
			err = rows.Scan(&c)
			if err != nil {
				return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
			}
			// Check if {c} is  "." and already exist in {categories} if not add to it
			// Before check split category at / and consider only root
			if c != "." && !utils.Contains(categories, strings.Split(c, "/")[0]) {
				categories = append(categories, strings.Split(c, "/")[0])
			}
		}
		// For each category retrieve {num} elements from db and add to dest slice to return
		for _, category := range categories {
			// Query run vs %s/%% to avoid aggregate similar category (eg DOCSrl and DOCSrlombardia)
			rows, err = DbConnection.Query(sqlFilteredCategories, num, fmt.Sprintf("%s/%%", category))
			if err != nil {
				return errors.New(fmt.Sprintf("Error retrieving documents: %v\n", err))
			}

			for rows.Next() {
				var d Document
				err = rows.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category, &d.IsDir, &d.CreationTime)
				if err != nil {
					return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
				}
				*dest = append(*dest, d)
			}
		}
	default:
		rows, err = DbConnection.Query(sqlStatement, num)
		if err != nil {
			return errors.New(fmt.Sprintf("Error retrieving documents: %v\n", err))
		}

		defer rows.Close()

		for rows.Next() {
			var d Document
			err = rows.Scan(&d.Id, &d.Hash, &d.FileName, &d.DisplayName, &d.Category, &d.IsDir, &d.CreationTime)
			if err != nil {
				return errors.New(fmt.Sprintf("Error scanning row: %v\n", err))
			}
			*dest = append(*dest, d)
		}
	}

	return nil
}

// GetHashTable Build hash table in form of hash:path
func (d Document) GetHashTable() (map[string]string, error) {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
		hashTable    map[string]string
	)

	sqlStatement = `select hash,filename from docs`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving hash/filename: %v\n", err))
	}

	defer rows.Close()

	hashTable = make(map[string]string)

	for rows.Next() {
		var d Document
		err = rows.Scan(&d.Hash, &d.FileName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning row %v\n", err))
		}
		hashTable[d.Hash] = d.FileName
	}

	return hashTable, nil
}

// TruncateTable Truncate (clean) actual table
func (d Document) TruncateTable() error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE docs`

	_, err = DbConnection.Exec(sqlStatement)
	if err != nil {
		return errors.New(fmt.Sprintf("Error truncating Document table"))
	}

	return nil
}

// BulkCreate Bulk create passed in document array
func (d Document) BulkCreate(docToAdd []Document) error {
	var (
		err          error
		sqlStatement *sql.Stmt //Prepared sql statement
		txn          *sql.Tx   //DB transaction
	)

	//Begin new transaction
	txn, err = DbConnection.Begin()
	if err != nil {
		return err
	}

	//Prepare insert statement
	sqlStatement, err = txn.Prepare(pq.CopyIn("docs", "hash", "filename", "displayname", "category", "isDir", "creationtime"))

	//Exec insert for every passed document
	for _, doc := range docToAdd {
		_, err = sqlStatement.Exec(doc.Hash, doc.FileName, doc.DisplayName, doc.Category, doc.IsDir, doc.CreationTime)
		if err != nil {
			return err
		}
	}

	//Flush actual data
	_, err = sqlStatement.Exec()
	if err != nil {
		return err
	}

	err = sqlStatement.Close()
	if err != nil {
		return err
	}

	//Execute transaction
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
