package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to MySQL database
	db := connectToMySql()
	if db == nil {
		log.Fatal("Error connecting to database")
	}
}

func connectToMySql() *sql.DB {
	var count int = 0
	for {
		db, err := sql.Open("mysql", "mehul:mehulpassword@tcp(localhost:3306)/reservation")
		if err != nil {
			if count < 5 {
				log.Println("Error connecting to database. Retrying...")
				count++
			} else {
				log.Fatalf("Error connecting to database: %v", err)
				break
			}
		} else {
			// Try to ping the sql
			err := db.Ping()
			if err != nil {
				if count < 5 {
					log.Println("Error pinging the database. Retrying...")
					count++
					continue
				} else {
					log.Fatalf("Error pinging the database: %v", err)
					break
				}
			}
			return db
		}
	}

	return nil
}
