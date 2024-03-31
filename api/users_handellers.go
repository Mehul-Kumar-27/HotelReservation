package api

import (
	"context"

	"github.com/Mehul-Kumar-27/HotelReservation/db"
	//"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandeller struct {
	useStore db.UserStore
}

func NewUserHandeller(us db.UserStore) *UserHandeller {
	return &UserHandeller{useStore: us}

}

func (h *UserHandeller) HandleGetUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	context := context.Background()
	u, err := h.useStore.GetUserById(context, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(u)
}

func (h *UserHandeller)HandleGetAllUsers(ctx *fiber.Ctx) error {
	context := context.Background()
	users, err := h.useStore.GetAllUsers(context)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(users)
}
