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
	reviewChn := make(chan types.Review)
	userBookingChn := make(chan types.UserBookings)

	var allBookings []types.Booking
	var reviews []types.Review
	var userBookings []types.UserBookings

	hotelAndRatings := make(map[string]float32)

	go collectBookingData(bookingChn, &allBookings)
	go collectReviews(reviewChn, &reviews)
	go collectUserBookings(userBookingChn, &userBookings)

	for _, hotelId := range hotelIDs {
		wg.Add(1)
		go func(hid string) {
			defer wg.Done()
			defer func() { <-sem }()

			sem <- struct{}{}

			userIds := queryToGetRandomUser(50, h.db)

			for _, userID := range userIds {
				booking := generateBooking(hid, userID)
				bookingChn <- booking
				review := generateReview(hid, userID)
				hotelAndRatings[hid] += review.Rating

				reviewChn <- review

				userBookingChn <- generateUserBooking(userID, review.HotelID, booking.BookingID)
			}
		}(hotelId)
	}

	wg.Wait()
	close(bookingChn)
	close(reviewChn)
	close(userBookingChn)

	for hotelID, rating := range hotelAndRatings {
		hotelAndRatings[hotelID] = rating / float32(50)
	}

	log.Printf("Total booking generated are %v", len(allBookings)) ///250,000
	log.Printf("Total reviews generated are %v", len(reviews))      ///250,000
	log.Printf("Total user bookings generated are %v", len(userBookings)) /// 250,000

	/// Now I have the bookings data, reviews, user bookings and the hotel reviews
	/// Will make 4 go routines to make them happen simantenously lets see if it works
	
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

func collectUserBookings(userBookingsChn <-chan types.UserBookings, userBookings *[]types.UserBookings) {
	for usrBooking := range userBookingsChn {
		*userBookings = append(*userBookings, usrBooking)
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

func generateUserBooking(userID string, hotelID string, bookingID string) types.UserBookings {
	var userBooking types.UserBookings

	userBooking.BookingID = bookingID
	userBooking.UserID = userID
	userBooking.HotelID = hotelID

	return userBooking
}
