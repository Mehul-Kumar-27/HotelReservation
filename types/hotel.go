package types

type Hotel struct {
	HotelName   string  `json:"hotelname" fake:"{company}"`
	City        string  `json:"city" fake:"{city}"`
	Country     string  `json:"country" fake:"{country}"`
	Street      string  `json:"street" fake:"{street}"`
	Rooms       int     `json:"rooms" fake:"{number:1,300}"`
	PricePerDay float64 `json:"pricePerDay" fake:"{number:100,500}"`
	Email       string  `json:"email" fake:"{email}"`
	HotelID     string  `json:"hotelid" fake:"{uuid}"`
	Phone       string  `json:"phone" fake:"{phone}"`
}
