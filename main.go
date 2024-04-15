package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://localhost:27017"
	dataBase       = "reservation"
	userCollection = "users"
)

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Fatal("Error while connecting to mongo: ", err)
	}

	err = pingTheMongo(mongoClient)
	if err != nil {
		log.Fatal("Error while pinging the mongo: ", err)
	}
	listenAddress := flag.String("listen", ":8080", "server listen address")
	flag.Parse()

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				}
				ctx.Status(code)
				return ctx.JSON(fiber.Map{
					"error": err.Error(),
				})
			},
		},
	)

	app.Listen(*listenAddress)
}

func connectToMongo() (*mongo.Client, error) {
	log.Println("Connecting to mongo")
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	connect, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error while connecting to mongo: ", err)
		return nil, err
	}

	return connect, nil
}

func pingTheMongo(client *mongo.Client) error {
	log.Println("Pinging the mongo")
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Error while pinging the mongo: ", err)
		return err
	}

	return nil
}
