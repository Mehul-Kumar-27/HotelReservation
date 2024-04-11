package main

import (
	"database/sql"
)

type DatabaseHandlerInterface interface {
	CreateUserTable() error
	CreateFakeUsers(count int) error
	///

	CreateHotelTable() error
	CreateHotelFakeData(count int) error

	///
	CreateBookingsTable()error
	CreateBookingFakeData(count int)error
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

///////////////
