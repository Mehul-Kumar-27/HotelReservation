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
		NewRoute("auth", "POST", jwt, userLogin),
		NewRoute("auth", "GET", jwt),
	}
)

var userHandler *UserHandler

func init() {
	log.Println("Initializing the user handler")
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

	ctx.Header("AcessToken", "Bearer "+response.accesstoken)
	ctx.JSON(response.status, gin.H{"message": response.response})
	ctx.Abort()
}

func jwt(ctx *gin.Context) {
	var jwtPayload JWTPayload

	acesstoken := ctx.GetHeader("Authorization")
	jwtPayload.Token = acesstoken

	response := userHandler.AuthRequest.JwtAuthService(ctx, jwtPayload)
	if response.status == 200 {
		ctx.JSON(response.status, gin.H{"stauts": response.status, "userid": response.accesstoken, "message": response.response})
		ctx.Abort()
	} else {

		ctx.Next()
		ctx.JSON(response.status, gin.H{"stauts": response.status, "userid": response.accesstoken, "message": response.response})

	}
}
