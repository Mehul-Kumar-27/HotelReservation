package main

import (
	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)

type UserFunctions interface {
	CreateUsers(n int) []types.User
}

type UserFunctionsImpl struct {
}

func NewUserFunctionsImpl() UserFunctions{
	return &UserFunctionsImpl{}
}

func (uf *UserFunctionsImpl) CreateUsers(n int) []types.User {
	var userList []types.User
	for i := 0; i < n; i++ {
		var user types.User
		gofakeit.Struct(&user)
		userList = append(userList, user)
	}

	return userList
}
