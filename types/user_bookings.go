package types

type UserBookings struct {
	BookingID string `json:"bookingid"`
	UserID    string `json:"userid"`
	HotelID   string `json:"hotelid"`
}
