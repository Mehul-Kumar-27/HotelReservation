package main

import (
	"database/sql"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
)

type DatabaseHandlerInterface interface {
	CreateUserTable() error
	CreateFakeUsers(count int) error
	///

	CreateHotelTable() error
	CreateHotelFakeData(count int) error

	///
	CreateBookingsTable() error
	CreateBookingFakeData(count int) error

	//

	CreateReviewsTable() error
	StoreReviewFakeData(reviews []types.Review) error

	//
	CreateUserBookingTable() error
	StoreUserBookings(userBookings []types.UserBookings) error
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
