package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

type DatabaseHandlerInterface interface {
	CreateUserTable() error
	CreateFakeUsers(count int) error
}

type DatabaseHandler struct {
	databaseHandlerInterface DatabaseHandlerInterface
}

func NewDatabaseHandler(databaseHandlerInterface DatabaseHandlerInterface) *DatabaseHandler {
	return &DatabaseHandler{databaseHandlerInterface: databaseHandlerInterface}
}

type SqlDataHandeler struct {
	db *sql.DB
}

func NewSQLDataHandler(db *sql.DB) DatabaseHandlerInterface {
	return &SqlDataHandeler{db: db}
}

func NewDatabaseHandlerForSQL(db *sql.DB) *DatabaseHandler {
	sqlDataHandler := NewSQLDataHandler(db)
	return NewDatabaseHandler(sqlDataHandler)
}

func (h *SqlDataHandeler) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS USERS (
        USERID VARCHAR(255) PRIMARY KEY NOT NULL,
        FIRSTNAME VARCHAR(255) NOT NULL,
        LASTNAME VARCHAR(255) NOT NULL,
        EMAIL VARCHAR(255) NOT NULL,
        PHONE VARCHAR(255) NOT NULL,
        BookingsID JSON
    )`

	_, err := h.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Table created successfully")
	return nil
}

func (h *SqlDataHandeler) CreateFakeUsers(count int) error {
	users := make([]types.User, count)
	for i := 0; i < count; i++ {
		user := types.User{}
		if err := gofakeit.Struct(&user); err != nil {
			return fmt.Errorf("error generating user: %w", err)
		}
		users[i] = user
	}

	query := "INSERT INTO USERS (USERID, FIRSTNAME, LASTNAME, EMAIL, PHONE, BookingsID) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if Commit is not called

	for _, user := range users {
		_, err := stmt.Exec(user.UserID, user.FirstName, user.LastName, user.Email, user.Phone, user.BookingsID)
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Println("Users inserted successfully")
	return nil
}
