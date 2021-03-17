package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`     //User name
	Surname  string `json:"surname"`  //User surname
	Username string `json:"username"` //User username
	Password string `json:"password"` //User hashed password
}

// Retrieve user information based on username
func (u *User) Get(username string) error {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
	)

	sqlStatement = `select id,name,surname,username,password from users where username=$1`

	row = DbConnection.QueryRow(sqlStatement, username)
	switch err = row.Scan(&u.Id, &u.Name, &u.Surname, &u.Username, &u.Password); err {
	case sql.ErrNoRows:
		return errors.New("no row where retrieved")
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("error retrieving user from db: %v\n", err))
	}
}
