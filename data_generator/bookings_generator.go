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
	bookingChn := make(chan *types.Booking)
	reviewChn := make(chan *types.Review)
	userBookingChn := make(chan *types.UserBookings)

	var allBookings []types.Booking
	var reviews []types.Review
	var userBookings []types.UserBookings
	var hotelAndRatings = make(map[string]float32)

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
				bookingChn <- &booking
				review := generateReview(hid, userID)
				reviewChn <- &review

				userBoooking := generateUserBooking(userID, review.HotelID, booking.BookingID)
				userBookingChn <- &userBoooking
			}
		}(hotelId)
	}

	wg.Wait()
	defer close(bookingChn)
	defer close(reviewChn)
	defer close(userBookingChn)

	for hotelID, rating := range hotelAndRatings {
		hotelAndRatings[hotelID] = rating / float32(50)
	}

	log.Printf("Total booking generated are %v", len(allBookings))        ///250,000
	log.Printf("Total reviews generated are %v", len(reviews))            ///250,000
	log.Printf("Total user bookings generated are %v", len(userBookings)) /// 250,000

	for _, review := range reviews {
		hotelAndRatings[review.HotelID] += review.Rating
	}

	for hotelID, rating := range hotelAndRatings {
		hotelAndRatings[hotelID] = rating / float32(50)
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		err := h.StoreReviewFakeData(reviews)
		if err != nil {
			log.Fatalf("Error storing the fake review data , %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := h.AddReviewOfHotelsToDataBase(&hotelAndRatings)
		if err != nil {
			log.Fatalf("Error adding the hotels ratings %v", err)
		}
	}()
	err = h.StoreTheBookingsData(allBookings)
	if err != nil {
		log.Fatalf("Error storing bookings data , %v", err)
	}
	go func() {
		defer wg.Done()
		err := h.StoreUserBookings(userBookings)
		if err != nil {
			log.Fatalf("Error storing user bookings data, %v", err)
		}
	}()
	wg.Wait()
	log.Println("Would you believe this we reached here")
	return nil
}

func collectBookingData(bookingChan <-chan *types.Booking, allBookings *[]types.Booking) {
	for booking := range bookingChan {
		*allBookings = append(*allBookings, *booking)
	}
}

func collectReviews(reviewChn <-chan *types.Review, allReviews *[]types.Review) {
	for review := range reviewChn {
		*allReviews = append(*allReviews, *review)
	}
}

func collectUserBookings(userBookingsChn <-chan *types.UserBookings, userBookings *[]types.UserBookings) {
	for usrBooking := range userBookingsChn {
		*userBookings = append(*userBookings, *usrBooking)
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

func (h *SqlDataHandeler) StoreTheBookingsData(bookings []types.Booking) error {
	batch := 1000
	query := `INSERT INTO BOOKINGS (BOOKINGID, HOTELID, USERID, BOOKINGDATE, NUMBEROFDAYS, NUMBEROFROOMS) VALUES (?,?,?,?,?,?)`
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := h.db.Begin()

	if err != nil {
		return err
	}

	for i := 0; i < len(bookings); i += batch {
		end := i + batch
		if end > len(bookings) {
			end = len(bookings)
		}

		batch := bookings[i:end]
		for _, booking := range batch {
			_, err := stmt.Exec(booking.BookingID, booking.HotelID, booking.UserID, booking.BookingDate, booking.NumberOfDays, booking.NumberOfRooms)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Bookings inserted successfully")

	return nil
}
