package db

import (
	"context"
	"log"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	GetUserById(ctx context.Context,id string) (*types.User, error)
	GetAllUsers(ctx context.Context) ([]*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{client: client}
}

func (mus *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User

	usrCol := mus.client.Database(DBNAME).Collection(userCollection)
	oid, er := ToObjectId(id)
	if er != nil {
		log.Println("Error in GetUserById: ", er)
		return nil, er
	
	}
	err := usrCol.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		log.Println("Error in GetUserById: ", err)
		return nil, err
	}

	return &user, nil

}

func (mus *MongoUserStore) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User

	usrCol := mus.client.Database(DBNAME).Collection(userCollection)

	cursor, err := usrCol.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error in GetAllUsers: ", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user types.User
		err := cursor.Decode(&user)
		if err != nil {
			log.Println("Error in GetAllUsers: ", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
