package types

type Review struct {
	ReviewID      string `bson:"_id,omitempty" json:"reviewid,omitempty" fake:"{uuid}"`
	HotelID string `bson:"hotelID" json:"hotelid" fake:"skip"`
	UserID  string `bson:"userID" json:"userid" fake:"skip"`
	Review  string `bson:"review" json:"review" fake:"{sentence:10}"`
}
