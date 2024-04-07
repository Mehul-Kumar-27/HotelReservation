package main

import (
	"log"
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	var booking types.Booking
	gofakeit.Struct(&booking)
	booking.BookingDate = time.Now()
	log.Printf("Booking %v", booking.BookingDate)
	log.Printf("Days %v", booking.NumberOfDays)
	log.Printf("Rooms %v", booking.NumberOfRooms)

	newBookingDate := booking.BookingDate.AddDate(0, 0, booking.NumberOfDays)
	log.Printf("New Booking Date: %v", newBookingDate)
}
