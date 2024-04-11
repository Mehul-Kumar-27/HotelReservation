package main

import (
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	var review types.Review
	gofakeit.Struct(&review)
	log.Printf("%v", review.Rating)
	log.Printf("%v", review.Review)
	log.Printf("%v", review.ReviewID)
}
