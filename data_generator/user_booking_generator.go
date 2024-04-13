package main

import (
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
)

func (h *SqlDataHandeler) CreateUserBookingTable() error {

	query := `CREATE TABLE IF NOT EXISTS USERBOOKINGS (
		USERID VARCHAR(255),
		HOTELID VARCHAR(255),
		BOOKINGID VARCHAR(255),
		PRIMARY KEY (USERID, HOTELID, BOOKINGID),
		FOREIGN KEY (USERID) REFERENCES USERS(USERID),
		FOREIGN KEY (HOTELID) REFERENCES HOTELS(HOTELID),
		FOREIGN KEY (BOOKINGID) REFERENCES BOOKINGS(BOOKINGID)
	)`

	_, err := h.db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Create the user bookings table")
	return nil
}

func (h *SqlDataHandeler) StoreUserBookings(userBookings []types.UserBookings) error {

	batch := 1000

	query := `INSERT INTO USERBOOKINGS
		(USERID, HOTELID, BOOKINGID) VALUES (?, ? , ?)`

	stmt, err := h.db.Prepare(query)

	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing the statement: %v", err)
		}
	}()

	tx, err := h.db.Begin()
	if err != nil {
		log.Printf("Error starting the transactions %v", err)
		return err
	}

	for i := 0; i < len(userBookings); i += batch {
		end := i + batch
		if end > len(userBookings) {
			end = len(userBookings)
		}

		bookings := userBookings[i:end]
		for _, booking := range bookings {
			_, err := stmt.Exec(booking.UserID, booking.HotelID, booking.BookingID)
			if err != nil {
				log.Printf("Error inserting the user bookings in the database %v", err)
				tx.Rollback()
				return err
			}
		}

	}
	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting the user transactions %v", err)
		return err
	}

	log.Println("User Booking Data Added Successfully")

	return nil

}
