package main

import (
	"flag"

	"github.com/Mehul-Kumar-27/HotelReservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddress := flag.String("listen", ":8080", "server listen address")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", api.HandleGetUsers)

	app.Listen(*listenAddress)
}