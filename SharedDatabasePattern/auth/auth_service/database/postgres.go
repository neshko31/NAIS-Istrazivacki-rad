package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"time"
)

type Database struct{
	DB *sql.DB
}

func NewDatabase(connStr string) (*Database, error){
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected")
	return &Database{DB: db}, nil
}

func addUser(db *sql.DB, new_username, new_password string){
	insertSQL := `insert into users (username, password) values ($1, $2)`

	_, err := db.Exec(insertSQL, new_username, new_password)
	if err != nil {
		log.Fatalf("Failed to insert into users: %w", err)
	}
	fmt.Println("Registered successuffuly")
}

func getUser(db *sql.DB){
	rows, err := db.Query("select * from users;")
	if err != nil {
		log.Fatalf("Failed to read users: %w", err)
	}
	defer rows.Close()

	for rows.Next(){
		var id int
		var username, password string
		if err := rows.Scan(&username, &password, &id); err != nil{
			log.Fatal(err)
		}
		fmt.Println("ID: %d, username: %s, password: %s\n", id, username, password)
	}
}

func updateUser(db *sql.DB, user_id int, new_username, new_password string){
	updateSQL := `update users set username = $2, password = $3 where id = $1;`

	_, err := db.Exec(updateSQL, user_id, new_username, new_password)
	if err != nil{
		log.Fatalf("Failed to update user: %w", err)
	}

	fmt.Println("User updated")
}

func deleteUser(db *sql.DB, user_id int){
	deleteSQL := `delete from users where id = $1`

	_,err := db.Exec(deleteSQL, user_id)
	if err != nil{
		log.Fatalf("Failed to delete user: %w", err)
	}

	fmt.Println("User deleted")
}

func test(test string){
	fmt.Println("Called from another: ", test)
}
