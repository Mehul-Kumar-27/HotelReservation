package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	User = NewRouter("user/", userRoutes)
)

var (
	userRoutes = []Route{
		NewRoute("auth", "POST", userLogin),
	}
)

var userHandler *UserHandler

func init() {
	log.Println("Intiallizing the user handler")
	userHandler = NewUserHandler(NewAuth())
}

type UserHandler struct {
	AuthRequest AuthRequest
}

func NewUserHandler(authRequest AuthRequest) *UserHandler {
	return &UserHandler{
		AuthRequest: authRequest,
	}
}

func userLogin(ctx *gin.Context) {
	var loginPayload LoginPayload
	if err := ctx.Bind(&loginPayload); err != nil {
		ctx.IndentedJSON(400, "bad formatted request")
		return
	}
	log.Println(loginPayload.Email)
	log.Println(loginPayload.Userid)
	response := userHandler.AuthRequest.LoginService(ctx, loginPayload)

	ctx.JSON(response.staus, response)
	
}
