package types

type User struct {
	UserID       string         `bson:"_id,omitempty" json:"userid,omitempty" fake:"{uuid}"`
	FirstName    string         `bson:"firstname" json:"firstname" fake:"{firstname}"`
	LastName     string         `bson:"lastname" json:"lastname" fake:"{lastname}"`
	Email        string         `bson:"email" json:"email" fake:"{email}"`
	Phone        string         `bson:"phone" json:"phone" fake:"{phone}"`
	Password     string 		`bson:"password" json:"password" fake:"{password:true,true,true,false,false,10}"`
}

