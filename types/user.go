package types

type User struct {
	UserID     string   `bson:"_id,omitempty" json:"userid,omitempty" fake:"{uuid}"`
	FirstName  string   `bson:"firstname" json:"firstname" fake:"{firstname}"`
	LastName   string   `bson:"lastname" json:"lastname" fake:"{lastname}"`
	Email      string   `bson:"email" json:"email" fake:"{email}"`
	Phone      string   `bson:"phone" json:"phone" fake:"{phone}"`
	BookingsID []string `bson:"bookingsId" json:"bookingsId" fake:"skip"`
}


