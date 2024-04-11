package main

import (
	"log"
	"sync"
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

func (h *SqlDataHandeler) CreateBookingsTable() error {
	query := `CREATE TABLE IF NOT EXISTS BOOKINGS(
		BOOKINGID VARCHAR(255) PRIMARY KEY NOT NULL,
		HOTELID VARCHAR(255),
		USERID VARCHAR(255),
		BOOKINGDATE TIMESTAMP NOT NULL,
        NUMBEROFDAYS INT NOT NULL,
        NUMBEROFROOMS INT NOT NULL,
        FOREIGN KEY (HOTELID) REFERENCES HOTELS(HOTELID),
        FOREIGN KEY (USERID) REFERENCES USERS(USERID)
	)`

	_, err := h.db.Exec(query)
	if err != nil {
		log.Fatalf("Error Creating the Bookings Table, %v", err)
	}
	log.Println("Successfully Created the Bookings Table")
	return nil
}

func (h *SqlDataHandeler) CreateBookingFakeData(count int) error {
	/// Lets get all the hotel ids from the database
	hotelIDs, err := queryToGetALLhotelID(h.db)
	if err != nil {
		log.Fatalf("Error fetching the hotels from the database, %v", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 100)
	bookingChn := make(chan types.Booking)
	var allBookings []types.Booking
	var reviews []types.Review

	reviewChn := make(chan types.Review)

	go collectBookingData(bookingChn, &allBookings)
	go collectReviews(reviewChn, &reviews)

	for _, hotelId := range hotelIDs {
		wg.Add(1)
		go func(hid string) {
			defer wg.Done()
			defer func() { <-sem }()

			sem <- struct{}{}

			userIds := queryToGetRandomUser(50, h.db)

			for _, userID := range userIds {

				bookingChn <- generateBooking(hid, userID)
				reviewChn <- generateReview(hid, userID)
			}
		}(hotelId)
	}

	wg.Wait()
	close(bookingChn)
	close(reviewChn)

	log.Printf("Total booking generated are %v", len(allBookings))
	log.Printf("Total reviews generated are %v", len(reviews))

	return nil
}

func collectBookingData(bookingChan <-chan types.Booking, allBookings *[]types.Booking) {
	for booking := range bookingChan {
		*allBookings = append(*allBookings, booking)
	}
}

func collectReviews(reviewChn <-chan types.Review, allReviews *[]types.Review) {
	for review := range reviewChn {
		*allReviews = append(*allReviews, review)
	}
}

func generateBooking(hotelID string, userID string) types.Booking {
	var booking types.Booking
	err := gofakeit.Struct(&booking)
	if err != nil {
		log.Fatalf("Error generating fake booking")

	}
	booking.UserID = userID
	booking.HotelID = hotelID
	booking.BookingDate = time.Now()

	return booking
}

func generateReview(hotelID string, userID string) types.Review {
	var review types.Review

	err := gofakeit.Struct(&review)
	if err != nil {
		log.Fatalf("Error generating fake review %v", err)

	}
	review.HotelID = hotelID
	review.UserID = userID

	return review
}
