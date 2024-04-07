package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

func (h *SqlDataHandeler) CreateHotelTable() error {

	query := `CREATE TABLE IF NOT EXISTS HOTELS(
		HOTELID VARCHAR(255) PRIMARY KEY NOT NULL,
		HOTELNAME VARCHAR(255) NOT NULL,
		CITY VARCHAR(255) NOT NULL,
		COUNTRY VARCHAR(255) NOT NULL,
		STREET VARCHAR(255) NOT NULL,
		ROOMS VARCHAR(255) NOT NULL,
		PRICEPERDAY FLOAT NOT NULL,
		EMAIL VARCHAR(255) NOT NULL,
		PHONE VARCHAR(255) NOT NULL	
	)`

	_, err := h.db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating the user table %v", err)
	}

	log.Println("Hotel Table Created successfully")

	return nil
}

func (h *SqlDataHandeler) CreateHotelFakeData(cont int) error {
	start := time.Now()
	var hotels []types.Hotel
	for i := 0; i < cont; i++ {
		var hotel types.Hotel
		gofakeit.Struct(&hotel)
		hotels = append(hotels, hotel)
	}
	query := `INSERT INTO HOTELS 
	(HOTELNAME, CITY, COUNTRY, STREET, ROOMS, PRICEPERDAY, EMAIL, HOTELID, PHONE) 
	VALUES (?,?,?,?,?,?,?,?,?)`
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing hotel insertion: %w", err)
	}
	defer stmt.Close()

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	batchSize := 1000 // Adjust batch size based on performance testing
	for i := 0; i < len(hotels); i += batchSize {
		end := i + batchSize
		if end > len(hotels) {
			end = len(hotels)
		}

		batch := hotels[i:end]
		for _, hotel := range batch {
			_, err := stmt.Exec(hotel.HotelName, hotel.City, hotel.Country, hotel.Street, hotel.Rooms, hotel.PricePerDay, hotel.Email, hotel.HotelID, hotel.Phone)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting hotel: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Println("Hotels inserted successfully")
	log.Printf("Time taken to insert the hotels is %v", time.Since(start))
	return nil
}
