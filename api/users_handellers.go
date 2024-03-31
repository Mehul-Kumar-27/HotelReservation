package api

import (
	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(ctx *fiber.Ctx) error {
	u := types.User{
		ID:        "1",
		FirstName: "Mehul",
		LastName:  "Kumar",
	}

	return ctx.JSON(u)
}
