package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const schema = `CREATE TABLE users (
    username text PRIMARY KEY NOT NULL,
    password text NOT NULL)`
	
type User struct {
	Username string
	Password string
}

func TestDb(t *testing.T) {
	fmt.Println("My test!")
	connStr := "host=localhost port=5432 dbname=neurotech user=neurotech password=neurotech"
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// _, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
	
	res := UserExists(db, "kelan")
	fmt.Println(res)
	userInsert := `INSERT INTO users (username, password) VALUES ($1, $2);`
	db.MustExec(userInsert, "kelan", "007")

	
}

func UserExists(db * sqlx.DB, username string) bool {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=$1 LIMIT 1", username)
	
	switch err {
	case nil:
		return true
	case sql.ErrNoRows:
		return false
	default:
		return false
	}
}	


func TestHello(t *testing.T) {
	fmt.Println("Hello test!")
}