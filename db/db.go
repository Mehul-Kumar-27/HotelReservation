package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBNAME = "reservation"

func ToObjectId(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}
