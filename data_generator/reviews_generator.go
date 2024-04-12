package main

import (
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
)

func (h *SqlDataHandeler) CreateReviewsTable() error {
	query := `CREATE TABLE IF NOT EXISTS REVIEWS(
		REVIEWID VARCHAR(255) PRIMARY KEY NOT NULL,
		HOTELID VARCHAR(255) NOT NULL,
		USERID VARCHAR(255) NOT NULL,
		RATING FLOAT NOT NULL,
		REVIEW TEXT NOT NULL,
		FOREIGN KEY (HOTELID) REFERENCES HOTELS(HOTELID),
		FOREIGN KEY (USERID) REFERENCES USERS(USERID)
	)`
	_, err := h.db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating the reviews table %v", err)
	}

	log.Println("Successfully created the reviews table")
	return nil
}

func (h *SqlDataHandeler) StoreReviewFakeData(reviews []types.Review) error {
	batch := 1000

	query := `INSERT INTO REVIEWS (REVIEWID, HOTELID, USERID, RATING, REVIEW) VALUES (?, ?, ?, ?, ?)`
	stmt, err := h.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing the statement: %v", err)
		}
	}()

	tx, err := h.db.Begin()

	if err != nil {
		log.Printf("Error starting the transaction: %v", err)
		return err
	}

	for i := 0; i < len(reviews); i += batch {
		end := i + batch
		if end > len(reviews) {
			end = len(reviews)
		}

		for _, review := range reviews[i:end] {
			if _, err := stmt.Exec(review.ReviewID, review.HotelID, review.UserID, review.Rating, review.Review); err != nil {
				log.Printf("Error inserting the reviews: %v", err)
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing the transaction: %v", err)
		return err
	}

	log.Println("Successfully inserted the reviews")

	return nil
}
