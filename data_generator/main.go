package main

import (
	"fmt"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {

	var hotel types.Hotel
	gofakeit.Struct(&hotel)
	fmt.Println(hotel.HotelName)
	fmt.Println(hotel.City)
	fmt.Println(hotel.Country)
	fmt.Println(hotel.Street)
	fmt.Println(hotel.Rooms)
	fmt.Println(hotel.PricePerDay)
	fmt.Println(hotel.Email)
	fmt.Println(hotel.HotelUUID)
	fmt.Println(hotel.Phone)


	/// 
	var review types.Review
	gofakeit.Struct(&review)
	fmt.Println(review.ReviewID)
	fmt.Println(review.Review)
}
