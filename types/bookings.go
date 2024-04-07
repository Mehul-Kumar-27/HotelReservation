package types

import "time"

type Booking struct {
	HotelID       string    `bson:"hotelID" json:"hoteid" fake:"skip"`
	UserID        string    `bson:"userID" json:"userid" fake:"skip"`
	BookingID     string    `bson:"_id,omitempty" json:"bookingid,omitempty" fake:"{uuid}"`
	BookingDate   time.Time `bson:"bookingDate" json:"bookingdate" fake:"skip"`
	NumberOfDays  int       `bson:"numberOfDays" json:"numberofdays" fake:"{number:1,10}"`
	NumberOfRooms int       `bson:"numberOfRooms" json:"numberofrooms" fake:"{number:1,10}"`
}
