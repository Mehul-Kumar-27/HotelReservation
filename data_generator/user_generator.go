package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func (h *SqlDataHandeler) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS USERS (
        USERID VARCHAR(255) PRIMARY KEY NOT NULL,
        FIRSTNAME VARCHAR(255) NOT NULL,
        LASTNAME VARCHAR(255) NOT NULL,
        EMAIL VARCHAR(255) NOT NULL,
        PHONE VARCHAR(255) NOT NULL,
		PASSWORD VARCHAR(255) NOT NULL,
        BOOKINGSID JSON
    )`

	_, err := h.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Table created successfully")
	return nil
}
func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing the password %v", err)
	}

	return string(hashedPassword), nil
}

func (h *SqlDataHandeler) CreateFakeUsers(count int) error {
	start := time.Now()
	log.Printf("Time started for the users insertions: %v \n", start)
	users := make([]types.User, count)
	for i := 0; i < count; i++ {
		user := types.User{}
		if err := gofakeit.Struct(&user); err != nil {
			return fmt.Errorf("error generating user: %w", err)
		}
		log.Println("Userid: %v, Password: %v", user.UserID, user.Password)
		user.Password, _ = HashedPassword(user.Password)
		users[i] = user
	}

	query := "INSERT INTO USERS (USERID, FIRSTNAME, LASTNAME, EMAIL, PHONE, PASSWORD) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing user insertion: %w", err)
	}
	defer stmt.Close()

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	batchSize := 1000 // Adjust batch size based on performance testing
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		for _, user := range batch {
			_, err := stmt.Exec(user.UserID, user.FirstName, user.LastName, user.Email, user.Phone, user.Password)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting user: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Println("Users inserted successfully")
	log.Printf("Time Taken to insert is %v", time.Since(start))
	return nil

}

func queryToGetRandomUser(count int, db *sql.DB) []string {

	var userList []string
	rows, err := db.Query("SELECT USERID FROM USERS ORDER BY RAND() LIMIT ?", count)
	if err != nil {
		status := db.Stats()
		log.Printf("At Error to get random users")
		log.Printf("%v", status.MaxOpenConnections)
		log.Printf("Number of open connections: %v", status.OpenConnections)
		log.Printf("Number of in use connections: %v", status.InUse)
		log.Printf("Number of idle connections: %v", status.Idle)
		log.Printf("Number of waiting connections: %v", status.WaitCount)
		log.Fatalf("Error fetching users from the user table, %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Error decoding the row obtained from the database, %v", err)
		}
		userList = append(userList, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error closing the rows, %v", err)
	}
	return userList
}
