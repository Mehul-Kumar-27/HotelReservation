package types

type Booking struct {
	HotelID string `bson:"hotelID" json:"hoteid" fake:"skip"`
	UserID  string `bson:"userID" json:"userid" fake:"skip"`
	BookingID string `bson:"_id,omitempty" json:"bookingid,omitempty" fake:"{uuid}"`
	BookingDate string `bson:"bookingDate" json:"bookingdate" fake:"{date}"`
	NuberOfDays int `bson:"numberOfDays" json:"numberofdays" fake:"{number:1,10}"`
	NumberOfRomms int `bson:"numberOfRooms" json:"numberofrooms" fake:"skip"`
}
